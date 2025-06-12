package ingestor

import (
	"bufio"
	"database/sql"
	"io/ioutil"
	"log"
	"log-analysis-service/internal/config"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogIngestor struct {
	db     *sql.DB
	config *config.Config
}

func New(db *sql.DB, cfg *config.Config) *LogIngestor {
	return &LogIngestor{
		db:     db,
		config: cfg,
	}
}

func (i *LogIngestor) Start() {
	log.Println("Starting log ingestor test")

	dirToType := map[string]string{
		i.config.LogDirs.DHCP:     "dhcp",
		i.config.LogDirs.DNS:      "dns",
		i.config.LogDirs.Firewall: "firewall",
		i.config.LogDirs.Network:  "network",
	}

	for dir, logType := range dirToType {
		i.processDirectory(dir, logType)
	}

	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		for dir, logType := range dirToType {
			i.processDirectory(dir, logType)
		}
	}
}

func (i *LogIngestor) processDirectory(dir string, logType string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Error reading directory %s: %v", dir, err)
		return
	}

	for _, file := range files {
		if file.IsDir() || strings.HasSuffix(file.Name(), ".processed") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		if err := i.processLogFile(filePath, logType); err != nil {
			log.Printf("Error processing log file %s: %v", filePath, err)
		} else {
			newPath := filePath + ".processed"
			if err := os.Rename(filePath, newPath); err != nil {
				log.Printf("Error renaming processed file %s: %v", filePath, err)
			}
		}
	}
}

func (i *LogIngestor) processLogFile(filePath string, logType string) error {
	log.Printf("Processing log file: %s (type: %s)", filePath, logType)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO raw_logs (source_type, log_time, raw_content) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		logTime := i.extractTimestamp(line)

		_, err := stmt.Exec(logType, logTime, line)
		if err != nil {
			return err
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Processed %d lines from %s", lineCount, filePath)
	return nil
}

func (i *LogIngestor) extractTimestamp(line string) time.Time {
	return time.Now()
}
