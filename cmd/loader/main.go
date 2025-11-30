package main

import (
	"flag"
	"log"

	"settlements/internal/config"
	"settlements/internal/db"
	"settlements/internal/service"
)

func main() {
	filePath := flag.String("file", "datasets/employee_travel_data.csv", "Path to the CSV file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	database, err := db.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Create data loader service
	loaderService := service.NewDataLoaderService(database)

	// Load data
	log.Printf("Loading data from %s...", *filePath)
	err = loaderService.LoadEmployeeTravelData(*filePath)
	if err != nil {
		log.Fatalf("failed to load data: %v", err)
	}

	log.Println("Data loaded successfully!")
}
