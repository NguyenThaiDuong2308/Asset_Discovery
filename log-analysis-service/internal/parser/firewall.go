package parser

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

type FirewallParser struct {
	connectionPattern *regexp.Regexp
}

func NewFirewallParser() *FirewallParser {
	return &FirewallParser{
		connectionPattern: regexp.MustCompile(`SRC=(\d+\.\d+\.\d+\.\d+)\s+DST=(\d+\.\d+\.\d+\.\d+)\s+.*?PROTO=(\w+)\s+SPT=(\d+)\s+DPT=(\d+)`),
	}
}

func (p *FirewallParser) Parse(rawContent string, logTime time.Time) (map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	matches := p.connectionPattern.FindStringSubmatch(rawContent)
	if matches == nil {
		return nil, nil, nil
	}
	log.Println(matches)
	sourceIP := matches[1]
	destIP := matches[2]
	protocol := matches[3]
	sourcePort, _ := strconv.Atoi(matches[4])
	destPort, _ := strconv.Atoi(matches[5])

	parsedLog := map[string]interface{}{
		"log_type":         "firewall_connection",
		"log_time":         logTime,
		"source_ip":        sourceIP,
		"destination_ip":   destIP,
		"source_port":      sourcePort,
		"destination_port": destPort,
		"protocol":         protocol,
		"action":           "connection",
		"details":          map[string]interface{}{},
	}

	sourceAsset := map[string]interface{}{
		"ip_address": sourceIP,
		"hostname":   "",
		"first_seen": logTime,
		"last_seen":  logTime,
		"asset_type": "client",
	}

	destAsset := map[string]interface{}{
		"ip_address": destIP,
		"hostname":   "",
		"first_seen": logTime,
		"last_seen":  logTime,
		"asset_type": "server",
	}
	var service map[string]interface{}
	switch destPort {
	case 80, 443, 8080:
		destAsset["asset_type"] = "web_server"
		service = map[string]interface{}{
			"name":       "http",
			"port":       destPort,
			"protocol":   protocol,
			"ip_address": destAsset["ip_address"],
		}
	case 22:
		destAsset["asset_type"] = "ssh_server"
		service = map[string]interface{}{
			"name":       "ssh",
			"port":       destPort,
			"protocol":   protocol,
			"ip_address": destAsset["ip_address"],
		}
	case 445:
		destAsset["asset_type"] = "smb_server"
		service = map[string]interface{}{
			"name":       "smb",
			"port":       destPort,
			"protocol":   protocol,
			"ip_address": destAsset["ip_address"],
		}
	case 3306:
		destAsset["asset_type"] = "database_server"
		service = map[string]interface{}{
			"name":       "mysql",
			"port":       destPort,
			"protocol":   protocol,
			"ip_address": destAsset["ip_address"],
		}
	case 5432:
		destAsset["asset_type"] = "database_server"
		service = map[string]interface{}{
			"name":       "postgresql",
			"port":       destPort,
			"protocol":   protocol,
			"ip_address": destAsset["ip_address"],
		}
	case 3389:
		destAsset["asset_type"] = "windows_server"
		service = map[string]interface{}{
			"name":       "rdp",
			"port":       destPort,
			"protocol":   protocol,
			"ip_address": destAsset["ip_address"],
		}
	}

	assets := []map[string]interface{}{sourceAsset, destAsset}
	var services []map[string]interface{}
	if service != nil {
		services = append(services, service)
	}

	return parsedLog, assets, services
}
