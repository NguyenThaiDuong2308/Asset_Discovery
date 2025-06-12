package parser

import (
	"regexp"
	"time"
)

type DNSParser struct {
	queryPattern *regexp.Regexp
	replyPattern *regexp.Regexp
}

func NewDNSParser() *DNSParser {
	return &DNSParser{
		queryPattern: regexp.MustCompile(`client ([\d\.]+)#(\d+) \(([^)]+)\).*?\(([\d\.]+)\)$`),
		replyPattern: regexp.MustCompile(`client ([\d\.]+)#(\d+): response: ([^\s]+) IN A ([\d\.]+)`),
	}
}

func (p *DNSParser) Parse(rawContent string, logTime time.Time) (map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	queryMatches := p.queryPattern.FindStringSubmatch(rawContent)
	if queryMatches != nil {
		clientIP := queryMatches[1]
		clientPort := queryMatches[2]
		hostname := queryMatches[3]
		dnsServerIP := queryMatches[4]
		parsedLog := map[string]interface{}{
			"log_type":         "dns_query",
			"log_time":         logTime,
			"source_ip":        clientIP,
			"source_port":      clientPort,
			"destination_ip":   dnsServerIP,
			"destination_port": 53,
			"action":           "query",
			"details":          map[string]string{},
		}

		asset := map[string]interface{}{
			"ip_address": clientIP,
			"hostname":   hostname,
			"first_seen": logTime,
			"last_seen":  logTime,
			"asset_type": "client",
		}
		services := []map[string]interface{}{nil}
		return parsedLog, []map[string]interface{}{asset}, services
	}

	replyMatches := p.replyPattern.FindStringSubmatch(rawContent)
	if replyMatches != nil {
		serverIP := replyMatches[1]
		serverPort := replyMatches[2]
		serverName := replyMatches[3]
		clientIP := replyMatches[4]
		parsedLog := map[string]interface{}{
			"log_type":       "dns_reply",
			"log_time":       logTime,
			"source_ip":      serverIP,
			"hostname":       serverName,
			"destination_ip": clientIP,
			"source_port":    serverPort,
			"action":         "reply",
		}

		clientAsset := map[string]interface{}{
			"ip_address": clientIP,
			"hostname":   "",
			"first_seen": logTime,
			"last_seen":  logTime,
			"asset_type": "client",
		}

		serverAsset := map[string]interface{}{
			"ip_address": serverIP,
			"hostname":   "",
			"first_seen": logTime,
			"last_seen":  logTime,
			"asset_type": "server",
		}

		serverService := map[string]interface{}{
			"name":       "dns",
			"port":       serverPort,
			"protocol":   "",
			"ip_address": serverIP,
		}

		return parsedLog,
			[]map[string]interface{}{clientAsset, serverAsset},
			[]map[string]interface{}{nil, serverService}
	}

	return nil, nil, nil
}
