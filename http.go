package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"
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

// proxyHandler proxies requests to their corresponding mapped domains.// proxyHandler decides the proxy method based on the target domain.
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	subdomain := strings.Split(r.Host, ".")[0]

	mutex.RLock()
	targetDomain, exists := DomainMap[subdomain]
	mutex.RUnlock()

	if !exists {
		http.Error(w, "Domain mapping not found", http.StatusNotFound)
		return
	}

	log.Printf("Proxying request for subdomain: %s with URL: %s", subdomain, r.URL.String())

	// Special handling for /CLOAK-CGI path
	if r.URL.Path == "/CLOAK-CGI" {
		http.FileServer(http.Dir(webDirFlag)).ServeHTTP(w, r)
		return
	}

	if strings.HasSuffix(targetDomain, ".onion") {
		proxyThroughSOCKS5(w, r, targetDomain)
	} else {
		proxyDirectly(w, r, targetDomain)
	}
}

// proxyDirectly handles direct proxying for non-.onion domains.
func proxyDirectly(w http.ResponseWriter, r *http.Request, targetDomain string) {

	// Parse the target URL from the domain map
	target, err := url.Parse("https://" + targetDomain)
	if err != nil {
		log.Printf("Error parsing target URL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	reverseProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Update the request with the target's host
	r.Host = target.Host

	// Proxy the request
	reverseProxy.ServeHTTP(w, r)
}

// proxyThroughSOCKS5 routes requests to .onion domains through a SOCKS5 proxy.
func proxyThroughSOCKS5(w http.ResponseWriter, r *http.Request, targetDomain string) {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		log.Printf("Error creating SOCKS5 dialer: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	transport := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{Transport: transport}

	// Modify the request to the target .onion address
	r.URL.Scheme = "http"
	r.URL.Host = targetDomain
	r.RequestURI = "" // The RequestURI field must be empty for client requests

	resp, err := client.Do(r)
	if err != nil {
		log.Printf("Error forwarding request through SOCKS5 proxy: %v", err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// copyHeader copies headers from the response to the original writer.
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// setupRoutes defines the HTTP routes and associates them with their handlers.
func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		subdomain := "" // Declare subdomain variable outside of the if-else block
		parts := strings.Split(r.Host, ".")
		if len(parts) > 2 {
			subdomain = parts[0] // Assign value to subdomain
		}
		if subdomain != "" {
			proxyHandler(w, r) // Route requests with subdomains to proxyHandler
		} else {
			homeHandler(w, r) // Route requests without subdomains to homeHandler
		}
	})

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
