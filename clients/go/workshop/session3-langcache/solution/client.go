package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Configuration constants
const (
	// AISHE API base URL - change this if your server is running elsewhere
	DefaultAISHEURL = "http://localhost:8000"

	// Redis address - used by benchmark.go for comparison
	DefaultRedisAddr = "localhost:6379"

	// LangCache configuration - update these with your LangCache credentials
	// Sign up at https://redis.io/langcache/ to get your credentials
	LangCacheURL    = "https://your-langcache-url.com"  // TODO: Update this
	LangCacheAPIKey = "your-api-key-here"                // TODO: Update this
	LangCacheCacheID = "your-cache-id-here"              // TODO: Update this
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

// LangCacheSearchRequest represents a search request to LangCache
type LangCacheSearchRequest struct {
	Prompt              string  `json:"prompt"`
	SimilarityThreshold float64 `json:"similarity_threshold,omitempty"`
}

// LangCacheSearchResponse represents a search response from LangCache
type LangCacheSearchResponse struct {
	EntryID  string  `json:"entry_id,omitempty"`
	Prompt   string  `json:"prompt,omitempty"`
	Response string  `json:"response,omitempty"`
	Score    float64 `json:"score,omitempty"`
}

// LangCacheSetRequest represents a set request to LangCache
type LangCacheSetRequest struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

// Client is an HTTP client with Redis LangCache integration for the AISHE API
type Client struct {
	baseURL          string
	httpClient       *http.Client
	langCacheURL     string
	langCacheAPIKey  string
	langCacheCacheID string
	threshold        float64
}

// ClientOptions contains configuration options for the client
type ClientOptions struct {
	BaseURL              string
	Timeout              time.Duration
	LangCacheURL         string  // LangCache API base URL
	LangCacheAPIKey      string  // LangCache API key
	LangCacheCacheID     string  // LangCache cache ID
	SimilarityThreshold  float64 // Similarity threshold for semantic matching (default: 0.9)
}

// NewClient creates a new HTTP client with Redis LangCache integration
func NewClient(opts ClientOptions) (*Client, error) {
	if opts.BaseURL == "" {
		opts.BaseURL = DefaultAISHEURL
	}

	if opts.Timeout == 0 {
		opts.Timeout = 120 * time.Second
	}

	if opts.LangCacheURL == "" {
		opts.LangCacheURL = LangCacheURL
	}
	if opts.LangCacheURL == "" || opts.LangCacheURL == "https://your-langcache-url.com" {
		return nil, fmt.Errorf("LangCache URL is required (update LangCacheURL constant in client.go)")
	}

	if opts.LangCacheAPIKey == "" {
		opts.LangCacheAPIKey = LangCacheAPIKey
	}
	if opts.LangCacheAPIKey == "" || opts.LangCacheAPIKey == "your-api-key-here" {
		return nil, fmt.Errorf("LangCache API key is required (update LangCacheAPIKey constant in client.go)")
	}

	if opts.LangCacheCacheID == "" {
		opts.LangCacheCacheID = LangCacheCacheID
	}
	if opts.LangCacheCacheID == "" || opts.LangCacheCacheID == "your-cache-id-here" {
		return nil, fmt.Errorf("LangCache cache ID is required (update LangCacheCacheID constant in client.go)")
	}

	if opts.SimilarityThreshold == 0 {
		opts.SimilarityThreshold = 0.9 // Default threshold
	}

	return &Client{
		baseURL: strings.TrimRight(opts.BaseURL, "/"),
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		langCacheURL:     strings.TrimRight(opts.LangCacheURL, "/"),
		langCacheAPIKey:  opts.LangCacheAPIKey,
		langCacheCacheID: opts.LangCacheCacheID,
		threshold:        opts.SimilarityThreshold,
	}, nil
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

// searchLangCache searches LangCache for a semantically similar cached response
func (c *Client) searchLangCache(question string) (*AnswerResponse, error) {
	url := fmt.Sprintf("%s/v1/caches/%s/entries/search", c.langCacheURL, c.langCacheCacheID)

	reqBody := LangCacheSearchRequest{
		Prompt:              question,
		SimilarityThreshold: c.threshold,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal LangCache request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create LangCache request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.langCacheAPIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search LangCache: %w", err)
	}
	defer resp.Body.Close()

	// 404 means no cache hit
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("LangCache returned status %d: %s", resp.StatusCode, string(body))
	}

	var cacheResp LangCacheSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&cacheResp); err != nil {
		return nil, fmt.Errorf("failed to decode LangCache response: %w", err)
	}

	// Parse the cached answer
	var answer AnswerResponse
	if err := json.Unmarshal([]byte(cacheResp.Response), &answer); err != nil {
		return nil, fmt.Errorf("failed to parse cached answer: %w", err)
	}

	return &answer, nil
}

// storeLangCache stores a response in LangCache
func (c *Client) storeLangCache(question string, answer *AnswerResponse) error {
	url := fmt.Sprintf("%s/v1/caches/%s/entries", c.langCacheURL, c.langCacheCacheID)

	// Serialize the answer
	answerJSON, err := json.Marshal(answer)
	if err != nil {
		return fmt.Errorf("failed to marshal answer: %w", err)
	}

	reqBody := LangCacheSetRequest{
		Prompt:   question,
		Response: string(answerJSON),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal LangCache set request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create LangCache set request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.langCacheAPIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to store in LangCache: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LangCache store returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// AskQuestion sends a question to the RAG system and returns the answer
// It uses Redis LangCache for semantic caching
func (c *Client) AskQuestion(question string) (*AnswerResponse, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	// Try to get from LangCache
	cachedAnswer, err := c.searchLangCache(question)
	if err != nil {
		// Log error but continue to API
		// In production, use proper logging
		_ = err
	} else if cachedAnswer != nil {
		// Cache hit!
		return cachedAnswer, nil
	}

	// Cache miss - fetch from API
	answer, err := c.fetchFromAPI(question)
	if err != nil {
		return nil, err
	}

	// Store in LangCache for future requests
	if err := c.storeLangCache(question, answer); err != nil {
		// Log error but don't fail the request
		_ = err
	}

	return answer, nil
}

// fetchFromAPI fetches the answer from the AISHE API
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

// FlushCache flushes all entries from the LangCache
func (c *Client) FlushCache() error {
	url := fmt.Sprintf("%s/v1/caches/%s/flush", c.langCacheURL, c.langCacheCacheID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create flush request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.langCacheAPIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to flush LangCache: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LangCache flush returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Close closes the HTTP client (no-op for standard http.Client)
func (c *Client) Close() error {
	// Standard http.Client doesn't need explicit closing
	return nil
}

