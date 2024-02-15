package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// Global variable to hold domain mappings with a mutex for thread-safe access
var (
	DomainMap map[string]string
	mutex     sync.RWMutex // Mutex to safely update the DomainMap
)

// readAvailableDictionary reads the dictionary file specified by the dictionaryLocationFlag
// and returns a slice of words.
func readAvailableDictionary() (words []string) {
	file, err := os.Open(dictionaryLocationFlag)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	words = strings.Split(string(bytes), "\n")
	return
}

// Babbler is a struct that generates random strings from a given dictionary of words.
type Babbler struct {
	Words []string
}

// NewBabbler creates a new Babbler instance with a loaded dictionary of words.
func NewBabbler() Babbler {
	return Babbler{
		Words: readAvailableDictionary(),
	}
}

// Babble generates a random string composed of three randomly chosen words from the dictionary.
func (b Babbler) Babble() string {
	pieces := make([]string, 3)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		pieces[i] = b.Words[rand.Intn(len(b.Words))]
	}
	return strings.Join(pieces, "-")
}

// loadMappings loads domain mappings from a specified file into the global DomainMap variable.
func loadMappings(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tempMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if len(parts) >= 2 {
			key := parts[0]
			domain := parts[1]
			tempMap[key] = domain
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	mutex.Lock()
	DomainMap = tempMap
	mutex.Unlock()

	return nil
}

// calculateMD5 calculates the MD5 hash of the file specified by filePath.
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// appendToDomainLog appends a log entry to the domain-specific log file.
func appendToDomainLog(domain, entry string) error {
	logFilePath := fmt.Sprintf("%s/%s.log", logDirFlag, domain)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(entry + "\n"); err != nil {
		return err
	}

	return nil
}
