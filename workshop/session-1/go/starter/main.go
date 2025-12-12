package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
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

	// Start timing
	startTime := time.Now()

	fmt.Printf("Asking: %s\n", question)
	fmt.Println("Waiting for response...\n")

	// TODO: Implement the AISHE API call
	//
	// Hints:
	// 1. Define structs to represent the API request and response
	// 2. Build a JSON payload with the question
	// 3. Make an HTTP POST request to the AISHE API endpoint
	// 4. Parse the JSON response
	// 5. Display the answer, sources, and processing time
	// 6. Display the total execution time using: time.Since(startTime).Seconds()
	// 7. Handle errors appropriately (connection, timeout, HTTP errors)

	_ = startTime // TODO: Remove this line when you implement execution time tracking
	panic("not implemented ...")
}

