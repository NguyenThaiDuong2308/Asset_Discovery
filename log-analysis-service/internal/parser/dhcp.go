package parser

import (
	"regexp"
	"time"
)

type DHCPParser struct {
	leasePattern *regexp.Regexp
}

func NewDHCPParser() *DHCPParser {
	return &DHCPParser{
		leasePattern: regexp.MustCompile(`DHCPACK.*?(\d+\.\d+\.\d+\.\d+).*?([0-9A-Fa-f:]+)`),
	}
}

func (p *DHCPParser) Parse(rawContent string, logTime time.Time) (map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	matches := p.leasePattern.FindStringSubmatch(rawContent)
	if matches == nil {
		return nil, nil, nil
	}

	ipAddress := matches[1]
	macAddress := matches[2]

	parsedLog := map[string]interface{}{
		"log_type":       "dhcp_lease",
		"log_time":       logTime,
		"destination_ip": ipAddress,
		"action":         "assign",
		"details":        map[string]string{"mac_address": macAddress},
	}

	asset := map[string]interface{}{
		"ip_address":  ipAddress,
		"mac_address": macAddress,
		"first_seen":  logTime,
		"last_seen":   logTime,
		"asset_type":  "client",
	}
	services := []map[string]interface{}{nil}
	return parsedLog, []map[string]interface{}{asset}, services
}
