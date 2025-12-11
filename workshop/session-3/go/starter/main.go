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

	"github.com/joho/godotenv"
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

// LangCacheClient represents a client for LangCache semantic caching
type LangCacheClient struct {
	ServerURL string
	CacheID   string
	APIKey    string
	HTTPClient *http.Client
}

// NewLangCacheClient creates a new LangCache client
func NewLangCacheClient(serverURL, cacheID, apiKey string) *LangCacheClient {
	return &LangCacheClient{
		ServerURL: serverURL,
		CacheID:   cacheID,
		APIKey:    apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getFromCache searches for a cached response using semantic search
//
// Hints:
// 1. Create a search request:
//    - Define a struct with fields: Prompt (string) and SimilarityThreshold (float64)
//    - Set Prompt to the question
//    - Set SimilarityThreshold to 0.8 (allows semantic matches)
//    - Marshal to JSON using json.Marshal()
//
// 2. Build the API URL:
//    - Format: "{ServerURL}/cache/{CacheID}/search"
//    - Use fmt.Sprintf() to construct the URL
//
// 3. Create and send HTTP POST request:
//    - Use http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
//    - Set headers:
//      * Content-Type: application/json
//      * Authorization: Bearer {APIKey}
//    - Use client.HTTPClient.Do(req) to send
//
// 4. Parse the response:
//    - Define a struct for the response with a Data field (array of entries)
//    - Each entry should have Prompt and Response fields
//    - Use json.NewDecoder(resp.Body).Decode()
//
// 5. Extract the cached data:
//    - If Data array is empty, return nil (cache miss)
//    - Get the first entry (most similar)
//    - The Response field contains JSON string - unmarshal it to Response struct
//
// 6. Handle errors appropriately
func getFromCache(client *LangCacheClient, question string) (*Response, error) {
	panic("not implemented")
}

// saveToCache saves response to semantic cache
//
// Hints:
// 1. Convert response to JSON string:
//    - Use json.Marshal() to convert Response struct to bytes
//    - Convert bytes to string
//
// 2. Create a set request:
//    - Define a struct with fields: Prompt (string) and Response (string)
//    - Set Prompt to the question
//    - Set Response to the JSON string from step 1
//    - Marshal to JSON using json.Marshal()
//
// 3. Build the API URL:
//    - Format: "{ServerURL}/cache/{CacheID}/set"
//    - Use fmt.Sprintf() to construct the URL
//
// 4. Create and send HTTP POST request:
//    - Use http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
//    - Set headers:
//      * Content-Type: application/json
//      * Authorization: Bearer {APIKey}
//    - Use client.HTTPClient.Do(req) to send
//
// 5. Check the response:
//    - Status code should be 200 or 201
//    - Return error if not successful
//
// 6. Return nil on success
func saveToCache(client *LangCacheClient, question string, response *Response) error {
	panic("not implemented")
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// Check if question was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <your question>")
		fmt.Println("Example: go run main.go 'What is the capital of France?'")
		os.Exit(1)
	}

	// Get question from command line arguments
	question := strings.Join(os.Args[1:], " ")

	// Get credentials from environment variables
	apiKey := os.Getenv("API_KEY")
	cacheID := os.Getenv("CACHE_ID")
	serverURL := os.Getenv("SERVER_URL")

	// Validate required credentials
	var missingFields []string
	if apiKey == "" {
		missingFields = append(missingFields, "API_KEY")
	}
	if cacheID == "" {
		missingFields = append(missingFields, "CACHE_ID")
	}
	if serverURL == "" || serverURL == "YOUR_REDIS_CLOUD_LANGCACHE_HOST_HERE" {
		missingFields = append(missingFields, "SERVER_URL")
	}

	if len(missingFields) > 0 {
		fmt.Println("Error: Missing or invalid credentials in .env file")
		fmt.Printf("Missing fields: %s\n", strings.Join(missingFields, ", "))
		fmt.Println("\nPlease update .env file with your Redis Cloud LangCache credentials:")
		fmt.Println("- SERVER_URL: Your Redis Cloud LangCache host (e.g., 'your-instance.redis.cloud')")
		fmt.Println("- CACHE_ID: Your cache ID")
		fmt.Println("- API_KEY: Your LangCache API key")
		os.Exit(1)
	}

	// Ensure server_url has https:// prefix
	if !strings.HasPrefix(serverURL, "http") {
		serverURL = "https://" + serverURL
	}

	// Initialize LangCache client
	langCache := NewLangCacheClient(serverURL, cacheID, apiKey)

	fmt.Printf("Asking: %s\n", question)

	// Check cache first using semantic search
	var data *Response
	var fromCache bool

	cachedResponse, err := getFromCache(langCache, question)
	if err == nil && cachedResponse != nil {
		fmt.Println("✓ Found in semantic cache! (no API call needed)\n")
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

		// Save to semantic cache for future use
		if err := saveToCache(langCache, question, data); err != nil {
			fmt.Printf("Warning: Error saving to cache: %v\n", err)
		} else {
			fmt.Println("✓ Response saved to semantic cache\n")
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
		fmt.Println("Source: Semantic Cache (LangCache)")
	} else {
		fmt.Printf("Processing time: %.2f seconds\n", data.ProcessingTime)
	}
	fmt.Println(strings.Repeat("=", 70))
}

