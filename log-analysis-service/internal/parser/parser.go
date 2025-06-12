package parser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type LogParser interface {
	Parse(rawContent string, logTime time.Time) (map[string]interface{}, []map[string]interface{}, []map[string]interface{})
}

type Parser struct {
	db      *sql.DB
	parsers map[string]LogParser
}

func New(db *sql.DB) *Parser {
	p := &Parser{
		db:      db,
		parsers: make(map[string]LogParser),
	}

	p.parsers["dhcp"] = NewDHCPParser()
	p.parsers["dns"] = NewDNSParser()
	p.parsers["firewall"] = NewFirewallParser()
	p.parsers["network"] = NewNetworkParser()

	return p
}

func (p *Parser) ProcessLogs(batchSize int) error {
	rows, err := p.db.Query(`
		SELECT id, source_type, log_time, raw_content 
		FROM raw_logs 
		WHERE processed = false 
		ORDER BY log_time
		LIMIT $1
	`, batchSize)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var sourceType string
		var logTime time.Time
		var content string

		if err := rows.Scan(&id, &sourceType, &logTime, &content); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		parser, exists := p.parsers[sourceType]
		if !exists {
			log.Printf("No parser for source type: %s", sourceType)
			continue
		}

		parsedLog, assets, services := parser.Parse(content, logTime)

		if parsedLog != nil {
			parsedLog["raw_log_id"] = id
			if err := p.saveParsedLog(parsedLog); err != nil {
				log.Printf("Error saving parsed log: %v", err)
			}
		}

		for _, asset := range assets {
			if err := p.saveAsset(asset); err != nil {
				log.Printf("Error saving asset: %v", err)
			}
		}

		for _, service := range services {
			log.Println(service)
			if err := p.saveService(service); err != nil {
				log.Printf("Error saving service: %v", err)
			}
		}

		if _, err := p.db.Exec("UPDATE raw_logs SET processed = true WHERE id = $1", id); err != nil {
			log.Printf("Error marking log as processed: %v", err)
		}
	}

	return rows.Err()
}

func (p *Parser) saveParsedLog(parsedLog map[string]interface{}) error {

	var detailsJSON []byte
	var err error

	if details, ok := parsedLog["details"]; ok && details != nil {
		detailsJSON, err = json.Marshal(details)
		if err != nil {
			return err
		}
	} else {
		detailsJSON = []byte("null")
	}

	_, err = p.db.Exec(`
		INSERT INTO parsed_logs (
			raw_log_id, log_type, log_time, source_ip, destination_ip, 
			source_port, destination_port, protocol, action, details
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		parsedLog["raw_log_id"],
		parsedLog["log_type"],
		parsedLog["log_time"],
		parsedLog["source_ip"],
		parsedLog["destination_ip"],
		parsedLog["source_port"],
		parsedLog["destination_port"],
		parsedLog["protocol"],
		parsedLog["action"],
		detailsJSON,
	)

	return err
}

func (p *Parser) saveAsset(asset map[string]interface{}) error {
	var id int
	var firstSeen time.Time

	err := p.db.QueryRow(`
		SELECT id, first_seen FROM assets 
		WHERE (ip_address = $1 OR mac_address = $2) 
		AND (ip_address IS NOT NULL OR mac_address IS NOT NULL)
		LIMIT 1
	`,
		asset["ip_address"],
		asset["mac_address"],
	).Scan(&id, &firstSeen)

	var metadataJSON []byte
	if metadata, ok := asset["metadata"]; ok && metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return err
		}
	} else {
		metadataJSON = []byte("null")
	}

	if err == nil {
		_, err = p.db.Exec(`
			UPDATE assets SET 
				ip_address = COALESCE($1, ip_address),
				mac_address = COALESCE($2, mac_address),
				hostname = COALESCE($3, hostname),
				last_seen = $4,
				asset_type = COALESCE($5, asset_type),
				metadata = COALESCE($6, metadata)
			WHERE id = $7
		`,
			asset["ip_address"],
			asset["mac_address"],
			asset["hostname"],
			asset["last_seen"],
			asset["asset_type"],
			metadataJSON,
			id,
		)
		return err
	}

	_, err = p.db.Exec(`
		INSERT INTO assets (
			ip_address, mac_address, hostname, first_seen, last_seen, asset_type, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		asset["ip_address"],
		asset["mac_address"],
		asset["hostname"],
		asset["first_seen"],
		asset["last_seen"],
		asset["asset_type"],
		metadataJSON,
	)

	return err
}

func (p *Parser) saveService(service map[string]interface{}) error {
	ip, ok := service["ip_address"].(string)
	if !ok || ip == "" {
		return fmt.Errorf("invalid or missing ip_address")
	}

	var existingID int
	err := p.db.QueryRow(`
		SELECT id FROM services
		WHERE asset_ip = $1 AND port = $2 AND protocol = $3 AND name = $4
	`, ip, service["port"], service["protocol"], service["name"]).Scan(&existingID)

	if err == sql.ErrNoRows {
		_, err := p.db.Exec(`
			INSERT INTO services (asset_ip, name, port, protocol)
			VALUES ($1, $2, $3, $4)
		`, ip, service["name"], service["port"], service["protocol"])
		return err
	} else if err != nil {
		return err
	}

	return nil
}

func (p *Parser) StartParsingLoop() {
	log.Println("Starting log parsing loop")

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		if err := p.ProcessLogs(100); err != nil {
			log.Printf("Error processing logs: %v", err)
		}
	}
}
