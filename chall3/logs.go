package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var logDir string

var filenameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+\.(txt|log)`)

func init() {
	logDir = os.Getenv("logDir")

	if logDir == "" {
		log.Printf("logDir environment variable is required")
		logDir = "logs"

	}
}
func validateFilename(path string) error {

	if strings.Contains(path, "flag") || strings.Contains(path, "app") || strings.Contains(path, "etc") {
		return fmt.Errorf("attack detected")
	}
	if !filenameRegex.MatchString(path) {
		return fmt.Errorf("invalid filename")
	}
	return nil
}
func readLogFile(filename string) ([]byte, error) {
	fmt.Printf("Reading log filename: %s\n", filename)
	if err := validateFilename(filename); err != nil {
		return nil, err
	}
	logpath := filepath.Join(logDir, filename)
	fmt.Printf("Reading logpath: %s\n", logpath)
	file, err := os.Open(logpath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 1024))
	if err != nil {
		return nil, err
	}
	return data, nil
}
func listLogFiles() (json.RawMessage, error) {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	data, err := json.Marshal(names)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func handleLogs(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		names, err := listLogFiles()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(names)
		return
	}
	data, err := readLogFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
