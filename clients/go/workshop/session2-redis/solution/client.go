package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

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
	BaseURL      string
	Timeout      time.Duration
	RedisAddr    string
	RedisDB      int
	RedisPass    string
	CacheTTL     time.Duration
}

// NewClient creates a new HTTP client with Redis caching
func NewClient(opts ClientOptions) (*Client, error) {
	if opts.BaseURL == "" {
		opts.BaseURL = DefaultAISHEURL
	}

	if opts.Timeout == 0 {
		opts.Timeout = 120 * time.Second
	}

	if opts.RedisAddr == "" {
		opts.RedisAddr = DefaultRedisAddr
	}

	if opts.CacheTTL == 0 {
		opts.CacheTTL = 1 * time.Hour // Default cache TTL
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.RedisAddr,
		Password: opts.RedisPass,
		DB:       opts.RedisDB,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		baseURL: strings.TrimRight(opts.BaseURL, "/"),
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		redisClient: rdb,
		cacheTTL:    opts.CacheTTL,
		ctx:         ctx,
	}, nil
}

// generateCacheKey generates a cache key from the question
func (c *Client) generateCacheKey(question string) string {
	hash := sha256.Sum256([]byte(strings.TrimSpace(strings.ToLower(question))))
	return fmt.Sprintf("aishe:answer:%x", hash)
}

// CheckHealth checks the health of the API server
func (c *Client) CheckHealth() (*HealthResponse, error) {
	url := fmt.Sprintf("%s/health", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &health, nil
}

// AskQuestion sends a question to the RAG system and returns the answer
// It first checks the Redis cache, and if not found, queries the API and caches the result
func (c *Client) AskQuestion(question string) (*AnswerResponse, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	// Generate cache key
	cacheKey := c.generateCacheKey(question)

	// Try to get from cache
	cachedData, err := c.redisClient.Get(c.ctx, cacheKey).Bytes()
	if err == nil {
		// Cache hit
		var answer AnswerResponse
		if err := json.Unmarshal(cachedData, &answer); err == nil {
			return &answer, nil
		}
		// If unmarshal fails, continue to fetch from API
	}

	// Cache miss or error - fetch from API
	answer, err := c.fetchFromAPI(question)
	if err != nil {
		return nil, err
	}

	// Cache the result
	answerData, err := json.Marshal(answer)
	if err == nil {
		// Ignore cache write errors - don't fail the request
		_ = c.redisClient.Set(c.ctx, cacheKey, answerData, c.cacheTTL).Err()
	}

	return answer, nil
}

// fetchFromAPI fetches the answer from the API
func (c *Client) fetchFromAPI(question string) (*AnswerResponse, error) {
	url := fmt.Sprintf("%s/api/v1/ask", c.baseURL)

	reqBody := QuestionRequest{Question: question}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("server error (%d): %s", resp.StatusCode, errResp.Detail)
		}
		return nil, fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	var answer AnswerResponse
	if err := json.Unmarshal(body, &answer); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &answer, nil
}

// ClearCache clears all cached answers
func (c *Client) ClearCache() error {
	iter := c.redisClient.Scan(c.ctx, 0, "aishe:answer:*", 0).Iterator()
	for iter.Next(c.ctx) {
		if err := c.redisClient.Del(c.ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete cache key: %w", err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan cache keys: %w", err)
	}
	return nil
}

// Close closes the HTTP client and Redis connection
func (c *Client) Close() error {
	return c.redisClient.Close()
}

