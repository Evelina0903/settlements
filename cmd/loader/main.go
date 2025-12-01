package main

import (
	"flag"
	"log"
	"settlements/internal/config"
	"settlements/internal/db"
	"settlements/internal/service/data_loader"
)

func main() {
	filePath := flag.String("file", "datasets/dataset.csv", "Path to the dataset CSV file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := db.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create data loader service
	loader := data_loader.New(db)

	// Load data
	log.Printf("Loading data from %s...", *filePath)
	err = loader.LoadCityData(*filePath)
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	log.Println("Data loaded successfully!")
}
