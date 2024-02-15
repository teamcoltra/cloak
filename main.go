package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialization functions from config.go
	// init() is automatically called due to its special function name in Go,
	// which initializes flags and loads configuration.

	// Load domain mappings and perform any necessary initial setup
	if err := loadMappings(mapFileFlag); err != nil {
		log.Fatalf("Failed to load domain mappings: %v", err)
	}

	// Ensure the necessary directories exist or are created
	ensureDirectoriesExist()

	// Set up HTTP routes from http.go
	setupRoutes()

	// Start the HTTP server
	log.Printf("Starting server on port %s\n", portFlag)
	if err := http.ListenAndServe(":"+portFlag, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ensureDirectoriesExist checks and creates necessary directories, such as log directories.
func ensureDirectoriesExist() {
	// Ensure logDir directory exists
	if _, err := os.Stat(logDirFlag); os.IsNotExist(err) {
		err := os.MkdirAll(logDirFlag, 0755)
		if err != nil {
			log.Fatalf("Failed to create log directory: %s", err)
		}
	}

	// Add more directories to check/create as needed
}
