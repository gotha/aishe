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

// Client is a basic HTTP client for the AISHE API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new basic HTTP client
// If baseURL is empty, it uses DefaultAISHEURL
func NewClient(baseURL string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = DefaultAISHEURL
	}

	if timeout == 0 {
		timeout = 120 * time.Second
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
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
func (c *Client) AskQuestion(question string) (*AnswerResponse, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

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

// Close closes the HTTP client (no-op for basic client, but included for consistency)
func (c *Client) Close() error {
	// HTTP client doesn't need explicit closing
	return nil
}

