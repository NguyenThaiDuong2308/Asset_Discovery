package models

import (
	"time"
)

type Asset struct {
	IPAddress       string    `json:"ip_address" db:"ip_address"`
	MACAddress      string    `json:"mac_address" db:"mac_address"`
	Hostname        string    `json:"hostname" db:"hostname"`
	AssetType       string    `json:"asset_type" db:"asset_type"`
	Location        string    `json:"location" db:"location"`
	OperatingSystem string    `json:"operating_system" db:"operating_system"`
	FirstSeen       time.Time `json:"first_seen" db:"first_seen"`
	LastSeen        time.Time `json:"last_seen" db:"last_seen"`
	IsManaged       bool      `json:"is_managed" db:"is_managed"`
	Services        []Service `json:"services,omitempty"`
}

type Service struct {
	ID          int    `json:"id" db:"id"`
	AssetIP     string `json:"asset_ip" db:"asset_ip"`
	Name        string `json:"name" db:"name"`
	Port        int    `json:"port" db:"port"`
	Protocol    string `json:"protocol" db:"protocol"`
	Description string `json:"description" db:"description"`
	IsManaged   bool   `json:"is_managed" db:"is_managed"`
}
