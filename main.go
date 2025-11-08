package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type TestRequest struct {
	Command string `json:"command"` // curl, ping, tracert, custom
	URL     string `json:"url"`
	Custom  string `json:"custom,omitempty"` // for custom commands
}

type TestResponse struct {
	Success    bool              `json:"success"`
	Command    string            `json:"command"`
	Output     string            `json:"output"`
	Error      string            `json:"error,omitempty"`
	Duration   string            `json:"duration"`
	Connection ConnectionInfo    `json:"connection"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type ConnectionInfo struct {
	Target    string `json:"target"`
	Timestamp string `json:"timestamp"`
	OS        string `json:"os"`
}

func main() {
	r := mux.NewRouter()

	// Enable CORS
	r.Use(corsMiddleware)

	// API endpoints
	r.HandleFunc("/api/test", handleTest).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/health", handleHealth).Methods("GET")

	// Serve static files (React frontend)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/build")))

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("OS: %s, ARCH: %s", runtime.GOOS, runtime.GOARCH)
	log.Fatal(http.ListenAndServe(port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	response := TestResponse{
		Command: req.Command,
		Connection: ConnectionInfo{
			Target:    req.URL,
			Timestamp: startTime.Format(time.RFC3339),
			OS:        runtime.GOOS,
		},
		Metadata: make(map[string]string),
	}

	var output string
	var err error

	switch req.Command {
	case "curl":
		output, err = executeCurl(req.URL)
	case "ping":
		output, err = executePing(req.URL)
	case "tracert":
		output, err = executeTraceroute(req.URL)
	case "custom":
		output, err = executeCustom(req.Custom, req.URL)
	default:
		err = fmt.Errorf("unknown command: %s", req.Command)
	}

	duration := time.Since(startTime)
	response.Duration = duration.String()

	if err != nil {
		response.Success = false
		response.Error = err.Error()
		response.Output = output
	} else {
		response.Success = true
		response.Output = output
	}

	// Add metadata
	response.Metadata["command_type"] = req.Command
	response.Metadata["execution_time"] = duration.String()

	json.NewEncoder(w).Encode(response)
}

func executeCurl(url string) (string, error) {
	// Add timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "curl", "-i", "-L", "--max-time", "30", url)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	if stderr.String() != "" {
		output += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		return output, fmt.Errorf("curl error: %v", err)
	}

	return output, nil
}

func executePing(target string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	// Remove protocol if present
	target = strings.TrimPrefix(target, "http://")
	target = strings.TrimPrefix(target, "https://")

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "ping", "-n", "4", target)
	} else {
		cmd = exec.CommandContext(ctx, "ping", "-c", "4", target)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	if stderr.String() != "" {
		output += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		return output, fmt.Errorf("ping error: %v", err)
	}

	return output, nil
}

func executeTraceroute(target string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Remove protocol if present
	target = strings.TrimPrefix(target, "http://")
	target = strings.TrimPrefix(target, "https://")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "tracert", target)
	} else {
		cmd = exec.CommandContext(ctx, "traceroute", target)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	if stderr.String() != "" {
		output += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		// Traceroute often returns non-zero even on success
		return output, nil
	}

	return output, nil
}

func executeCustom(customCmd, url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Replace {url} placeholder with actual URL
	cmdStr := strings.ReplaceAll(customCmd, "{url}", url)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", cmdStr)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", cmdStr)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	if stderr.String() != "" {
		output += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		return output, fmt.Errorf("custom command error: %v", err)
	}

	return output, nil
}
