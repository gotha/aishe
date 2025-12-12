package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

// CachedResponse wraps a Response with optional similarity score
type CachedResponse struct {
	Response   *Response
	Similarity *float64
}

// LangCacheClient represents a client for LangCache semantic caching
type LangCacheClient struct {
	ServerURL  string
	CacheID    string
	APIKey     string
	HTTPClient *http.Client
}

// LangCacheSearchRequest represents a search request to LangCache
type LangCacheSearchRequest struct {
	Prompt              string  `json:"prompt"`
	SimilarityThreshold float64 `json:"similarity_threshold"`
}

// LangCacheSearchEntry represents a single search result entry
type LangCacheSearchEntry struct {
	Prompt     string   `json:"prompt"`
	Response   string   `json:"response"`
	Similarity *float64 `json:"similarity,omitempty"` // Optional similarity score
}

// LangCacheSearchResponse represents the search response from LangCache
type LangCacheSearchResponse struct {
	Data []LangCacheSearchEntry `json:"data"`
}

// LangCacheSetRequest represents a set request to LangCache
type LangCacheSetRequest struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

// NewLangCacheClient creates a new LangCache client
func NewLangCacheClient(serverURL, cacheID, apiKey string) *LangCacheClient {
	return &LangCacheClient{
		ServerURL:  serverURL,
		CacheID:    cacheID,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// getFromCache searches for a cached response using semantic search
func getFromCache(client *LangCacheClient, question string, threshold float64) (*CachedResponse, error) {
	// Prepare search request
	searchReq := LangCacheSearchRequest{
		Prompt:              question,
		SimilarityThreshold: threshold,
	}

	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, err
	}

	// Build the search URL
	url := fmt.Sprintf("%s/v1/caches/%s/entries/search", client.ServerURL, client.CacheID)

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.APIKey)

	// Send request
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status %d", resp.StatusCode)
	}

	// Parse response
	var searchResp LangCacheSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	// Check if we got any results
	if len(searchResp.Data) == 0 {
		return nil, nil // No cache hit
	}

	// Get the first (most similar) entry
	entry := searchResp.Data[0]

	// Parse the cached response from JSON string
	var cachedData Response
	if err := json.Unmarshal([]byte(entry.Response), &cachedData); err != nil {
		return nil, err
	}

	// Return cached response with similarity score (if available)
	return &CachedResponse{
		Response:   &cachedData,
		Similarity: entry.Similarity,
	}, nil
}

// saveToCache saves response to semantic cache
func saveToCache(client *LangCacheClient, question string, response *Response) error {
	// Convert response to JSON string
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// Prepare set request
	setReq := LangCacheSetRequest{
		Prompt:   question,
		Response: string(responseJSON),
	}

	jsonData, err := json.Marshal(setReq)
	if err != nil {
		return err
	}

	// Build the set URL
	url := fmt.Sprintf("%s/v1/caches/%s/entries", client.ServerURL, client.CacheID)

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.APIKey)

	// Send request
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("set failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func main() {
	// Start timing
	startTime := time.Now()

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

	// Get similarity threshold from environment variable (default: 0.8)
	threshold := 0.8
	if thresholdStr := os.Getenv("SIMILARITY_THRESHOLD"); thresholdStr != "" {
		if parsedThreshold, err := strconv.ParseFloat(thresholdStr, 64); err == nil {
			threshold = parsedThreshold
		}
	}

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
	var similarity *float64

	cachedResponse, err := getFromCache(langCache, question, threshold)
	if err != nil {
		fmt.Printf("⚠ Cache lookup error: %v\n", err)
	}
	if cachedResponse != nil {
		fmt.Println("✓ Found in semantic cache! (no API call needed)")
		if cachedResponse.Similarity != nil {
			fmt.Printf("  Similarity score: %.4f\n", *cachedResponse.Similarity)
		}
		fmt.Println()
		data = cachedResponse.Response
		similarity = cachedResponse.Similarity
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

	// Print processing time or cache info
	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	if fromCache {
		fmt.Println("Source: Semantic Cache (LangCache)")
		if similarity != nil {
			fmt.Printf("Similarity score: %.4f\n", *similarity)
		}
		fmt.Printf("Original processing time: %.2f seconds\n", data.ProcessingTime)
	} else {
		fmt.Printf("Processing time: %.2f seconds\n", data.ProcessingTime)
	}
	fmt.Println(strings.Repeat("=", 70))

	// Print total execution time
	executionTime := time.Since(startTime).Seconds()
	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Execution time: %.2f seconds\n", executionTime)
	fmt.Println(strings.Repeat("=", 70))
}
