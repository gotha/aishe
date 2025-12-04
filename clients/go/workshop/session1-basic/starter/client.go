package main

import (
	"fmt"
	"net/http"
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
	// TODO: Implement this function
	// 1. If baseURL is empty, use DefaultAISHEURL constant
	// 2. If timeout is 0, set it to 120 seconds
	// 3. Create and return a Client with the baseURL and an http.Client
	// Hint: Use strings.TrimRight to remove trailing slashes from baseURL

	return nil
}

// CheckHealth checks the health of the API server
func (c *Client) CheckHealth() (*HealthResponse, error) {
	// TODO: Implement this function
	// 1. Build the URL: baseURL + "/health"
	// 2. Make a GET request using c.httpClient.Get()
	// 3. Check the status code (should be 200)
	// 4. Decode the JSON response into a HealthResponse struct
	// 5. Return the health response

	return nil, fmt.Errorf("not implemented")
}

// AskQuestion sends a question to the RAG system and returns the answer
func (c *Client) AskQuestion(question string) (*AnswerResponse, error) {
	// TODO: Implement this function
	// 1. Validate that question is not empty (after trimming whitespace)
	// 2. Build the URL: baseURL + "/api/v1/ask"
	// 3. Create a QuestionRequest with the question
	// 4. Marshal it to JSON
	// 5. Make a POST request with Content-Type: application/json
	// 6. Check the status code (should be 200)
	// 7. Decode the JSON response into an AnswerResponse struct
	// 8. Return the answer response

	return nil, fmt.Errorf("not implemented")
}

// Close closes the HTTP client (no-op for basic client, but included for consistency)
func (c *Client) Close() error {
	// HTTP client doesn't need explicit closing
	return nil
}
