package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	LogDirs struct {
		DHCP     string
		DNS      string
		Firewall string
		Network  string
	}
	API struct {
		Port string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.Database.Host = getEnv("DB_HOST", "db")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.DBName = getEnv("DB_NAME", "log_analysis")

	basePath := getEnv("LOG_BASE_DIR", "logs")
	cfg.LogDirs.DHCP = filepath.Join(basePath, "dhcp")
	cfg.LogDirs.DNS = filepath.Join(basePath, "dns")
	cfg.LogDirs.Firewall = filepath.Join(basePath, "firewall")
	cfg.LogDirs.Network = filepath.Join(basePath, "network")

	cfg.API.Port = getEnv("API_PORT", "8080")

	// Ensure log directories exist
	os.MkdirAll(cfg.LogDirs.DHCP, 0755)
	os.MkdirAll(cfg.LogDirs.DNS, 0755)
	os.MkdirAll(cfg.LogDirs.Firewall, 0755)
	os.MkdirAll(cfg.LogDirs.Network, 0755)

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
