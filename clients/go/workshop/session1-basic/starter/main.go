package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// TODO: Create a new basic client
	// Hint: Use NewClient("", 120*time.Second) to use default URL from constant
	client := NewClient("", 120*time.Second)
	defer client.Close()

	// TODO: Check server health
	fmt.Println("Checking server health...")
	// Hint: Call client.CheckHealth() and handle the error
	health, err := client.CheckHealth()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}

	// TODO: Print the health status
	fmt.Printf("Server Status: %s\n", health.Status)
	fmt.Printf("Ollama Accessible: %v\n", health.OllamaAccessible)
	if health.Message != nil {
		fmt.Printf("Message: %s\n", *health.Message)
	}
	fmt.Println()

	// Interactive question loop
	fmt.Println("=== AISHE Question Answering (Session 1: Basic Client) ===")
	fmt.Println("Type your question and press Enter. Type 'exit' to quit.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Question: ")
		if !scanner.Scan() {
			break
		}

		question := strings.TrimSpace(scanner.Text())

		// Check for exit command
		if strings.ToLower(question) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Skip empty questions
		if question == "" {
			continue
		}

		// TODO: Measure execution time and ask the question
		// Hint: Use time.Now() before calling AskQuestion, then time.Since(start).Seconds()
		start := time.Now()
		// TODO: Call client.AskQuestion(question) here
		// answer, err := client.AskQuestion(question)
		executionTime := time.Since(start).Seconds()

		// TODO: Handle errors
		// if err != nil {
		//     fmt.Printf("❌ Error: %v\n\n", err)
		//     continue
		// }

		// TODO: Display results
		// Expected output format:
		//
		// Answer: <the answer>
		//
		// Source: AISHE API
		// Execution Time: X.XX seconds
		//
		// Wikipedia Sources:
		//   [1] <title>
		//       <url>
		//

		fmt.Println()
		fmt.Println("TODO: Display the answer here")
		fmt.Printf("Execution Time: %.2f seconds\n", executionTime)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

