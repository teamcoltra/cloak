package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// homeHandler serves static files for the root URL and handles the home page requests.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.FileServer(http.Dir(webDirFlag)).ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// addDomainHandler handles requests to add a new domain mapping.
func addDomainHandler(w http.ResponseWriter, r *http.Request) {
	babbler := NewBabbler()
	query := r.URL.Query()
	url := query.Get("url")
	key := babbler.Babble() // Generate a unique key
	email := query.Get("email")

	if url == "" || key == "" || email == "" {
		http.Error(w, "Missing url, key, or email parameters", http.StatusBadRequest)
		return
	}

	if err := addMappingToFile(mapFileFlag, key, url, email); err != nil {
		http.Error(w, "Failed to add domain mapping", http.StatusInternalServerError)
		return
	}

	if err := loadMappings(mapFileFlag); err != nil {
		http.Error(w, "Failed to reload domain mappings", http.StatusInternalServerError)
		return
	}

	md5sum, err := calculateMD5(mapFileFlag)
	if err != nil {
		http.Error(w, "Failed to calculate MD5 sum", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, md5sum)
}

// getMapHandler returns the MD5 checksum of the mappings file.
func getMapHandler(w http.ResponseWriter, r *http.Request) {
	md5sum, err := calculateMD5(mapFileFlag)
	if err != nil {
		http.Error(w, "Failed to calculate MD5 sum", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, md5sum)
}

// proxyHandler proxies requests to their corresponding mapped domains.
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	subdomain := strings.Split(r.Host, ".")[0]
	mutex.RLock()
	targetDomain, exists := DomainMap[subdomain]
	mutex.RUnlock()

	if !exists {
		http.Error(w, "Domain mapping not found", http.StatusNotFound)
		return
	}

	target, err := url.Parse("https://" + targetDomain)
	if err != nil {
		log.Printf("Error parsing target URL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	reverseProxy.ServeHTTP(w, r)
}

// setupRoutes defines the HTTP routes and associates them with their handlers.
func setupRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/addDomain", addDomainHandler)
	http.HandleFunc("/getMap", getMapHandler)

	// Serve static files from the specified directory
	fs := http.FileServer(http.Dir(webDirFlag))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

// addMappingToFile adds a new domain mapping to the mappings file.
func addMappingToFile(filePath, key, domain, email string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s|%s|%s\n", key, domain, email))
	return err
}
