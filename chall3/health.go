package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

type SystemHealth struct {
	Status     string `json:"status"`
	Uptime     string `json:"uptime"`
	Goroutines int    `json:"goroutines"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	Flag       string `json:"flag"`
	DBStatus   string `json:"db_status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

var startPoint = time.Now()

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{"method not allowed"})
		return
	}

	flag := ""
	dbStatus := "DOWN"

	targetfile, err := os.Open("/flag.txt")
	if err == nil {
		defer targetfile.Close()
	}

	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		pingErr := MongoClient.Ping(ctx, nil)
		cancel()

		if pingErr == nil {
			dbStatus = "UP"
		}
	}

	if err == nil {
		data, readErr := io.ReadAll(io.LimitReader(targetfile, 1024))
		if readErr == nil {
			flag = string(data)
		} else {
			flag = "flag not found"
		}
	} else {
		flag = "flag not found"
	}

	// ==========================================

	resp := SystemHealth{
		Status:     "UP",
		Uptime:     time.Since(startPoint).Truncate(time.Second).String(),
		Goroutines: runtime.NumGoroutine(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Flag:       flag,
		DBStatus:   dbStatus,
	}

	statusCode := http.StatusOK
	if flag == "flag not found" || dbStatus == "DOWN" {
		resp.Status = "DOWN"
	}

	writeJSON(w, statusCode, resp)
}
