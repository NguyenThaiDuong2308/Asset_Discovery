package api

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) getAssetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	rows, err := s.db.Query(`
		SELECT id, ip_address, mac_address, hostname, first_seen, last_seen, asset_type 
		FROM assets 
		ORDER BY last_seen DESC
		LIMIT $1
	`, limit)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assets []map[string]interface{}

	for rows.Next() {
		var id int
		var ipAddress, macAddress, hostname, assetType sql.NullString
		var firstSeen, lastSeen time.Time

		if err := rows.Scan(&id, &ipAddress, &macAddress, &hostname, &firstSeen, &lastSeen, &assetType); err != nil {
			continue
		}

		asset := map[string]interface{}{
			"id":          id,
			"ip_address":  nullStringValue(ipAddress),
			"mac_address": nullStringValue(macAddress),
			"hostname":    nullStringValue(hostname),
			"first_seen":  firstSeen.Format(time.RFC3339),
			"last_seen":   lastSeen.Format(time.RFC3339),
			"asset_type":  nullStringValue(assetType),
		}

		assets = append(assets, asset)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"assets": assets,
		"count":  len(assets),
	})
}

func (s *Server) getAssetByIPHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip"]

	w.Header().Set("Content-Type", "application/json")

	var id int
	var ipAddress, macAddress, hostname, assetType sql.NullString
	var firstSeen, lastSeen time.Time
	var metadata []byte

	err := s.db.QueryRow(`
		SELECT id, ip_address, mac_address, hostname, first_seen, last_seen, asset_type, metadata
		FROM assets WHERE ip_address = $1
	`, ip).Scan(&id, &ipAddress, &macAddress, &hostname, &firstSeen, &lastSeen, &assetType, &metadata)

	if err != nil {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	asset := map[string]interface{}{
		"id":          id,
		"ip_address":  nullStringValue(ipAddress),
		"mac_address": nullStringValue(macAddress),
		"hostname":    nullStringValue(hostname),
		"first_seen":  firstSeen.Format(time.RFC3339),
		"last_seen":   lastSeen.Format(time.RFC3339),
		"asset_type":  nullStringValue(assetType),
		"metadata":    string(metadata),
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"asset": asset,
	})
}

func (s *Server) getLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logType := r.URL.Query().Get("type")

	var rows *sql.Rows
	var err error

	if logType != "" {
		rows, err = s.db.Query(`
			SELECT id, log_type, log_time, source_ip, destination_ip, action
			FROM parsed_logs
			WHERE log_type = $1
			ORDER BY log_time DESC
			LIMIT $2
		`, logType, limit)
	} else {
		rows, err = s.db.Query(`
			SELECT id, log_type, log_time, source_ip, destination_ip, action
			FROM parsed_logs
			ORDER BY log_time DESC
			LIMIT $1
		`, limit)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}

	for rows.Next() {
		var id int
		var logType, sourceIP, destIP, action sql.NullString
		var logTime time.Time

		if err := rows.Scan(&id, &logType, &logTime, &sourceIP, &destIP, &action); err != nil {
			continue
		}

		log := map[string]interface{}{
			"id":        id,
			"type":      nullStringValue(logType),
			"time":      logTime.Format(time.RFC3339),
			"source_ip": nullStringValue(sourceIP),
			"dest_ip":   nullStringValue(destIP),
			"action":    nullStringValue(action),
		}

		logs = append(logs, log)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}

func (s *Server) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var totalAssets int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM assets").Scan(&totalAssets); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rows, err := s.db.Query("SELECT asset_type, COUNT(*) FROM assets GROUP BY asset_type")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	assetTypes := make(map[string]int)

	for rows.Next() {
		var assetType sql.NullString
		var count int

		if err := rows.Scan(&assetType, &count); err != nil {
			continue
		}

		typeStr := "unknown"
		if assetType.Valid {
			typeStr = assetType.String
		}

		assetTypes[typeStr] = count
	}

	rows, err = s.db.Query("SELECT log_type, COUNT(*) FROM parsed_logs GROUP BY log_type")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	logTypes := make(map[string]int)

	for rows.Next() {
		var logType sql.NullString
		var count int

		if err := rows.Scan(&logType, &count); err != nil {
			continue
		}

		typeStr := "unknown"
		if logType.Valid {
			typeStr = logType.String
		}

		logTypes[typeStr] = count
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_assets":   totalAssets,
		"assets_by_type": assetTypes,
		"logs_by_type":   logTypes,
	})
}

func (s *Server) getServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	rows, err := s.db.Query(`SELECT id, asset_ip, name, port, protocol
	FROM services
	ORDER BY asset_ip
	LIMIT $1
	`, limit)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var services []map[string]interface{}
	for rows.Next() {
		var id, port int
		var assetIP, name, protocol sql.NullString
		if err := rows.Scan(&id, &assetIP, &name, &port, &protocol); err != nil {
			continue
		}

		service := map[string]interface{}{
			"id":       id,
			"asset_ip": nullStringValue(assetIP),
			"name":     nullStringValue(name),
			"port":     port,
			"protocol": nullStringValue(protocol),
		}
		services = append(services, service)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"services": services,
		"count":    len(services),
	})
}
func (s *Server) getServicesByAssetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip"]

	w.Header().Set("Content-Type", "application/json")

	rows, err := s.db.Query(`
		SELECT name, port, protocol
		FROM services
		WHERE asset_ip = $1
	`, ip)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var services []map[string]interface{}

	for rows.Next() {
		var name sql.NullString
		var port int
		var protocol sql.NullString

		if err := rows.Scan(&name, &port, &protocol); err != nil {
			continue
		}

		service := map[string]interface{}{
			"name":     nullStringValue(name),
			"port":     port,
			"protocol": nullStringValue(protocol),
		}

		services = append(services, service)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"services": services,
		"count":    len(services),
	})
}

func (s *Server) uploadLogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logType := vars["type"]

	validTypes := map[string]bool{
		"dhcp":     true,
		"dns":      true,
		"firewall": true,
		"network":  true,
	}

	if !validTypes[logType] {
		http.Error(w, "Invalid log type", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, _, err := r.FormFile("logfile")
	if err != nil {
		http.Error(w, "Failed to get log file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Prepare statement for inserting logs
	stmt, err := tx.Prepare("INSERT INTO raw_logs (source_type, log_time, raw_content) VALUES ($1, $2, $3)")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Process uploaded file
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Insert log
		_, err := stmt.Exec(logType, time.Now(), line)
		if err != nil {
			http.Error(w, "Failed to insert log", http.StatusInternalServerError)
			return
		}

		lineCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  fmt.Sprintf("Uploaded %d log lines", lineCount),
		"log_type": logType,
	})
}

// nullStringValue returns the string value of a sql.NullString or an empty string if it's null
func nullStringValue(s sql.NullString) interface{} {
	if !s.Valid {
		return nil
	}
	return s.String
}
