package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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
//
// Hints:
// 1. Normalize the question:
//    - Convert to lowercase using strings.ToLower()
//    - Trim whitespace using strings.TrimSpace()
//    This ensures "What is Python?" and "what is python?" produce the same key
//
// 2. Hash the normalized question:
//    - Use crypto/sha256 package
//    - hash := sha256.Sum256([]byte(normalized))
//    - Convert to hex string: hex.EncodeToString(hash[:])
//
// 3. Add a namespace prefix:
//    - Format: "aishe:question:{hash}"
//    - Use fmt.Sprintf() to create the final key
func getCacheKey(question string) string {
	panic("not implemented")
}

// getFromCache retrieves cached response for a question
//
// Hints:
// 1. Generate the cache key using getCacheKey()
//
// 2. Query Redis:
//    - Use client.Get(ctx, cacheKey).Result()
//    - You'll need context.Background() for ctx
//
// 3. Handle cache miss:
//    - If err == redis.Nil, return nil (cache miss)
//    - If other error, return the error
//
// 4. Parse the cached JSON:
//    - The cached value is a JSON string
//    - Use json.Unmarshal() to convert to Response struct
//    - Return pointer to the Response
func getFromCache(client *redis.Client, question string) (*Response, error) {
	panic("not implemented")
}

// saveToCache saves response to cache
//
// Hints:
// 1. Generate the cache key using getCacheKey()
//
// 2. Convert response to JSON:
//    - Use json.Marshal() to convert Response struct to JSON bytes
//
// 3. Save to Redis with expiration:
//    - Use client.Set(ctx, key, value, expiration)
//    - Set expiration to 24 hours: 24 * time.Hour
//    - You'll need context.Background() for ctx
//
// 4. Return any errors from the Set operation
func saveToCache(client *redis.Client, question string, response *Response) error {
	panic("not implemented")
}

func main() {
	// Check if question was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <your question>")
		fmt.Println("Example: go run main.go 'What is the capital of France?'")
		os.Exit(1)
	}

	// Get question from command line arguments
	question := strings.Join(os.Args[1:], " ")

	// Connect to Redis (running in Docker on port 6379)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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

		// AISHE server URL (running in Docker on port 8000)
		url := "http://localhost:8000/api/v1/ask"

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
}

