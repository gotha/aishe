package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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

// getCacheKey generates a cache key from the question
func getCacheKey(question string) string {
	// Normalize the question (lowercase, strip whitespace)
	normalized := strings.ToLower(strings.TrimSpace(question))
	// Create a hash for the cache key
	hash := sha256.Sum256([]byte(normalized))
	questionHash := hex.EncodeToString(hash[:])
	return fmt.Sprintf("aishe:question:%s", questionHash)
}

// getFromCache retrieves cached response for a question
func getFromCache(client *redis.Client, question string) (*Response, error) {
	ctx := context.Background()
	cacheKey := getCacheKey(question)

	cachedData, err := client.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// saveToCache saves response to cache
func saveToCache(client *redis.Client, question string, response *Response) error {
	ctx := context.Background()
	cacheKey := getCacheKey(question)

	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// Store with 24 hour expiration (86400 seconds)
	return client.Set(ctx, cacheKey, jsonData, 24*time.Hour).Err()
}

func main() {
	// Start timing
	startTime := time.Now()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional, continue with system environment variables
	}

	// Check if question was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <your question>")
		fmt.Println("Example: go run main.go 'What is the capital of France?'")
		os.Exit(1)
	}

	// Get question from command line arguments
	question := strings.Join(os.Args[1:], " ")

	// Get Redis address from environment variable (default: localhost:6379)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("Error: Could not connect to Redis at localhost:6379")
		fmt.Println("Make sure Redis is running in Docker.")
		os.Exit(1)
	}

	fmt.Printf("Asking: %s\n", question)

	// Check cache first
	var data *Response
	var fromCache bool

	cachedResponse, err := getFromCache(rdb, question)
	if err == nil && cachedResponse != nil {
		fmt.Println("✓ Found in cache! (no API call needed)\n")
		data = cachedResponse
		fromCache = true
	} else {
		fmt.Println("✗ Not in cache, calling AISHE API...")
		fmt.Println("Waiting for response...\n")

		// Get AISHE server URL from environment variable (default: http://localhost:8000)
		aisheURL := os.Getenv("AISHE_URL")
		if aisheURL == "" {
			aisheURL = "http://localhost:8000"
		}
		url := aisheURL + "/api/v1/ask"

		// Prepare request payload
		payload := Request{Question: question}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("Error marshaling request: %v\n", err)
			os.Exit(1)
		}

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
		data = &Response{}
		if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
			fmt.Printf("Error parsing response: %v\n", err)
			os.Exit(1)
		}

		// Save to cache for future use
		if err := saveToCache(rdb, question, data); err != nil {
			fmt.Printf("Warning: Error saving to cache: %v\n", err)
		} else {
			fmt.Println("✓ Response saved to cache\n")
		}
		fromCache = false
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
	if fromCache {
		fmt.Println("Source: Redis Cache")
	} else {
		fmt.Printf("Processing time: %.2f seconds\n", data.ProcessingTime)
	}
	fmt.Println(strings.Repeat("=", 70))

	// Print total execution time
	executionTime := time.Since(startTime).Seconds()
	fmt.Println()
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Execution time: %.2f seconds\n", executionTime)
	fmt.Println(strings.Repeat("-", 70))
}

