package parser

import (
	"regexp"
	"time"
)

type NetworkParser struct {
	macIPPattern *regexp.Regexp
}

func NewNetworkParser() *NetworkParser {
	return &NetworkParser{
		macIPPattern: regexp.MustCompile(`MAC\s+([0-9A-Fa-f:]+).*?IP\s+(\d+\.\d+\.\d+\.\d+)`),
	}
}

func (p *NetworkParser) Parse(rawContent string, logTime time.Time) (map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	matches := p.macIPPattern.FindStringSubmatch(rawContent)
	if matches == nil {
		return nil, nil, nil
	}

	macAddress := matches[1]
	ipAddress := matches[2]

	parsedLog := map[string]interface{}{
		"log_type":  "mac_ip_mapping",
		"log_time":  logTime,
		"source_ip": ipAddress,
		"action":    "mapping",
		"details":   map[string]string{"mac_address": macAddress},
	}

	asset := map[string]interface{}{
		"ip_address":  ipAddress,
		"mac_address": macAddress,
		"first_seen":  logTime,
		"last_seen":   logTime,
		"asset_type":  "network_device",
	}

	return parsedLog, []map[string]interface{}{asset}, nil
}
