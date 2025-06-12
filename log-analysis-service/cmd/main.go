package main

import (
	"log"
	"log-analysis-service/internal/api"
	"log-analysis-service/internal/config"
	"log-analysis-service/internal/database"
	"log-analysis-service/internal/ingestor"
	"log-analysis-service/internal/parser"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logIngestor := ingestor.New(db, cfg)
	logParser := parser.New(db)
	apiServer := api.New(db)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		logIngestor.Start()
	}()

	go func() {
		defer wg.Done()
		logParser.StartParsingLoop()
	}()

	go func() {
		if err := apiServer.Start(cfg.API.Port); err != nil {
			log.Fatalf("API server error: %v", err)
		}
	}()

	log.Println("Log Analysis Service is running")
	log.Printf("API server listening on port %s", cfg.API.Port)

	<-stop
	log.Println("Shutting down...")

	wg.Wait()
	log.Println("Service stopped")
}
