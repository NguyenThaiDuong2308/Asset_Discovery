package models

import (
	"database/sql"
	"time"
)

type RawLog struct {
	ID         int       `json:"id"`
	SourceType string    `json:"source_type"`
	LogTime    time.Time `json:"log_time"`
	RawContent string    `json:"raw_content"`
	Processed  bool      `json:"processed"`
	IngestedAt time.Time `json:"ingested_at"`
}

type ParsedLog struct {
	ID              int            `json:"id"`
	RawLogID        int            `json:"raw_log_id"`
	LogType         string         `json:"log_type"`
	LogTime         time.Time      `json:"log_time"`
	SourceIP        sql.NullString `json:"source_ip"`
	DestinationIP   sql.NullString `json:"destination_ip"`
	SourcePort      sql.NullInt32  `json:"source_port"`
	DestinationPort sql.NullInt32  `json:"destination_port"`
	Protocol        sql.NullString `json:"protocol"`
	Action          sql.NullString `json:"action"`
	Details         []byte         `json:"details"`
}

type Asset struct {
	ID         int            `json:"id"`
	IPAddress  sql.NullString `json:"ip_address"`
	MACAddress sql.NullString `json:"mac_address"`
	Hostname   sql.NullString `json:"hostname"`
	FirstSeen  time.Time      `json:"first_seen"`
	LastSeen   time.Time      `json:"last_seen"`
	AssetType  sql.NullString `json:"asset_type"`
	Metadata   []byte         `json:"metadata"` // JSONB in database
	Services   []Service
}

type Service struct {
	ID       int    `json:"id"`
	AssetID  int    `json:"asset_id"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}
