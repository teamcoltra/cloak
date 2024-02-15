package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Global variables for configuration flags
var (
	webDirFlag             string
	mapFileFlag            string
	logDirFlag             string
	portFlag               string
	apiKeyFlag             string
	dictionaryLocationFlag string
	configFilePathFlag     string // Special flag for the config file path
)

func init() {
	// Define all flags here
	flag.StringVar(&webDirFlag, "webDir", "/var/www/html", "Directory where the index.html file lives")
	flag.StringVar(&mapFileFlag, "map", "mappings.txt", "Path to the domain mappings file")
	flag.StringVar(&logDirFlag, "logDir", "/var/log/cloak", "Directory where logs should go")
	flag.StringVar(&portFlag, "port", "8080", "What port to listen to")
	flag.StringVar(&apiKeyFlag, "apiKey", "", "API KEY - THIS IS USED TO MESS STUFF UP")
	flag.StringVar(&dictionaryLocationFlag, "dictionary", "dictionary.txt", "Dictionary of words to be used by babble")
	flag.StringVar(&configFilePathFlag, "config", "/etc/cloak/cloak.cfg", "Path to the configuration file")

	// Parse flags to capture command line inputs
	flag.Parse()

	// After parsing, load the configuration file if specified to potentially override flag values
	if configFilePathFlag != "" {
		loadConfig(configFilePathFlag)
	}

	// Environment variable fallback for API key if not set via flag or config
	if apiKeyFlag == "" {
		apiKeyFlag = os.Getenv("API_KEY")
	}
}

func loadConfig(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file is optional, so exit gracefully if not found
			return
		}
		fmt.Printf("Error opening config file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			// Ignore comments and empty lines
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Skip invalid lines
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		// Override flags based on config file values
		switch key {
		case "webDir":
			webDirFlag = value
		case "map":
			mapFileFlag = value
		case "logDir":
			logDirFlag = value
		case "port":
			portFlag = value
		case "apiKey":
			apiKeyFlag = value
		case "dictionary":
			dictionaryLocationFlag = value
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
	}
}
