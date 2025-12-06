package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO: You'll need to import these packages as you implement the functions:
// - "bytes"
// - "crypto/sha256"
// - "encoding/json"
// - "io"
// - "strings"

// Configuration constants
const (
	// AISHE API base URL - change this if your server is running elsewhere
	DefaultAISHEURL = "http://localhost:8000"

	// Redis address - change this if your Redis is running elsewhere
	DefaultRedisAddr = "localhost:6379"
)

// Source represents a source reference in the API response
type Source struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

// AnswerResponse represents the response from the /api/v1/ask endpoint
type AnswerResponse struct {
	Answer         string   `json:"answer"`
	Sources        []Source `json:"sources"`
	ProcessingTime float64  `json:"processing_time"`
}

// HealthResponse represents the response from the /health endpoint
type HealthResponse struct {
	Status           string  `json:"status"`
	OllamaAccessible bool    `json:"ollama_accessible"`
	Message          *string `json:"message,omitempty"`
}

// QuestionRequest represents the request to the /api/v1/ask endpoint
type QuestionRequest struct {
	Question string `json:"question"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// Client is an HTTP client with Redis caching for the AISHE API
type Client struct {
	baseURL     string
	httpClient  *http.Client
	redisClient *redis.Client
	cacheTTL    time.Duration
	ctx         context.Context
}

// ClientOptions contains configuration options for the client
type ClientOptions struct {
	BaseURL   string
	Timeout   time.Duration
	RedisAddr string
	RedisDB   int
	RedisPass string
	CacheTTL  time.Duration
}

// NewClient creates a new HTTP client with Redis caching
func NewClient(opts ClientOptions) (*Client, error) {
	// TODO: Implement this function
	// 1. Set defaults for BaseURL (use DefaultAISHEURL constant if empty)
	// 2. Set default timeout (120 seconds)
	// 3. Set default RedisAddr (use DefaultRedisAddr constant if empty)
	// 4. Set default CacheTTL (1 hour)
	// 5. Create a Redis client using redis.NewClient()
	// 6. Test the Redis connection with Ping()
	// 7. Return the Client struct

	return nil, fmt.Errorf("not implemented")
}

// generateCacheKey generates a cache key from the question
func (c *Client) generateCacheKey(question string) string {
	// TODO: Implement this function
	// 1. Normalize the question (trim whitespace, convert to lowercase)
	// 2. Generate a SHA-256 hash of the normalized question
	// 3. Return a cache key in the format: "aishe:answer:<hash>"
	// Hint: Use crypto/sha256 and fmt.Sprintf
	
	return ""
}

// CheckHealth checks the health of the API server
func (c *Client) CheckHealth() (*HealthResponse, error) {
	// TODO: Copy implementation from Session 1
	// This is the same as the basic client
	
	return nil, fmt.Errorf("not implemented")
}

// fetchFromAPI fetches the answer from the API
func (c *Client) fetchFromAPI(question string) (*AnswerResponse, error) {
	// TODO: Implement this function
	// This is similar to AskQuestion from Session 1, but without caching logic
	// 1. Build the URL: baseURL + "/api/v1/ask"
	// 2. Create a QuestionRequest and marshal to JSON
	// 3. Make a POST request
	// 4. Check status code and decode response
	// 5. Return the AnswerResponse
	
	return nil, fmt.Errorf("not implemented")
}

// AskQuestion sends a question to the RAG system and returns the answer
// It first checks the Redis cache, and if not found, queries the API and caches the result
func (c *Client) AskQuestion(question string) (*AnswerResponse, error) {
	// TODO: Implement this function with caching logic
	// 1. Validate the question is not empty
	// 2. Generate a cache key using generateCacheKey()
	// 3. Try to get the cached answer from Redis using c.redisClient.Get()
	//    - If found, unmarshal and return it (cache hit!)
	//    - If not found (redis.Nil error), continue to step 4
	// 4. Fetch from API using fetchFromAPI()
	// 5. Marshal the answer to JSON
	// 6. Store in Redis using c.redisClient.Set() with the TTL
	// 7. Return the answer
	
	return nil, fmt.Errorf("not implemented")
}

// ClearCache clears all cached answers
func (c *Client) ClearCache() error {
	// TODO: Implement this function
	// 1. Use c.redisClient.Scan() to find all keys matching "aishe:answer:*"
	// 2. Delete each key using c.redisClient.Del()
	// Hint: Use an iterator pattern with Scan
	
	return fmt.Errorf("not implemented")
}

// Close closes the HTTP client and Redis connection
func (c *Client) Close() error {
	return c.redisClient.Close()
}

