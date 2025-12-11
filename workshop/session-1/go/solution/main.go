package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Request represents the API request payload
type Request struct {
	Question string `json:"question"`
}

// Source represents a source citation
type Source struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

// Response represents the API response
type Response struct {
	Answer         string   `json:"answer"`
	Sources        []Source `json:"sources"`
	ProcessingTime float64  `json:"processing_time"`
}

func main() {
	// Start timing
	startTime := time.Now()

	// Check if question was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <your question>")
		fmt.Println("Example: go run main.go 'What is the capital of France?'")
		os.Exit(1)
	}

	// Get question from command line arguments
	question := strings.Join(os.Args[1:], " ")

	// AISHE server URL (running in Docker on port 8000)
	url := "http://localhost:8000/api/v1/ask"

	// Prepare request payload
	payload := Request{Question: question}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Asking: %s\n", question)
	fmt.Println("Waiting for response...\n")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	// Send POST request to AISHE server
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: Could not connect to AISHE server at %s\n", url)
		fmt.Println("Make sure the server is running in Docker.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error: Server returned status %d\n", resp.StatusCode)
		if len(body) > 0 {
			fmt.Printf("Details: %s\n", string(body))
		}
		os.Exit(1)
	}

	// Parse response
	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	// Print answer
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("ANSWER:")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(data.Answer)

	// Print sources if available
	if len(data.Sources) > 0 {
		fmt.Println()
		fmt.Println(strings.Repeat("=", 70))
		fmt.Println("SOURCES:")
		fmt.Println(strings.Repeat("=", 70))
		for _, source := range data.Sources {
			fmt.Printf("[%d] %s\n", source.Number, source.Title)
			fmt.Printf("    %s\n", source.URL)
		}
	}

	// Print processing time
	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Processing time: %.2f seconds\n", data.ProcessingTime)
	fmt.Println(strings.Repeat("=", 70))

	// Print total execution time
	executionTime := time.Since(startTime).Seconds()
	fmt.Println()
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Execution time: %.2f seconds\n", executionTime)
	fmt.Println(strings.Repeat("-", 70))
}

