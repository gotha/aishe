package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Check if question was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <your question>")
		fmt.Println("Example: go run main.go 'What is the capital of France?'")
		os.Exit(1)
	}

	// Get question from command line arguments
	question := strings.Join(os.Args[1:], " ")

	fmt.Printf("Asking: %s\n", question)
	fmt.Println("Waiting for response...\n")

	// TODO: Implement the AISHE API call
	//
	// Hints:
	// 1. Define structs for Request and Response to match the API format
	//    - Request: {"question": "..."}
	//    - Response: {"answer": "...", "sources": [...], "processing_time": 2.45}
	//
	// 2. Create the request payload:
	//    - Use json.Marshal() to convert your Request struct to JSON
	//    - The API endpoint is: http://localhost:8000/api/v1/ask
	//
	// 3. Make the HTTP POST request:
	//    - Use http.Client with a timeout (e.g., 120 seconds)
	//    - Set Content-Type header to "application/json"
	//    - Use client.Post() or http.NewRequest() + client.Do()
	//
	// 4. Handle the response:
	//    - Check for HTTP errors (resp.StatusCode != 200)
	//    - Use json.NewDecoder(resp.Body).Decode() to parse the JSON response
	//    - Don't forget to defer resp.Body.Close()
	//
	// 5. Display the results:
	//    - Print the answer
	//    - Print sources (if available)
	//    - Print processing time
	//
	// 6. Handle errors gracefully:
	//    - Connection errors (server not running)
	//    - Timeout errors
	//    - HTTP errors (non-200 status codes)

	panic("not implemented ...")
}

