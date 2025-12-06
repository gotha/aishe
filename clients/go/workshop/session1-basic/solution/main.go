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
	// Create a new basic client
	client := NewClient("", 120*time.Second)
	defer client.Close()

	// Check health
	fmt.Println("Checking server health...")
	health, err := client.CheckHealth()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
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

		// Measure execution time
		start := time.Now()
		answer, err := client.AskQuestion(question)
		executionTime := time.Since(start).Seconds()

		if err != nil {
			fmt.Printf("❌ Error: %v\n\n", err)
			continue
		}

		// Display results
		fmt.Println()
		fmt.Println("Answer:", answer.Answer)
		fmt.Println()
		fmt.Printf("Source: AISHE API\n")
		fmt.Printf("Execution Time: %.2f seconds\n", executionTime)
		fmt.Println()

		if len(answer.Sources) > 0 {
			fmt.Println("Wikipedia Sources:")
			for _, source := range answer.Sources {
				fmt.Printf("  [%d] %s\n", source.Number, source.Title)
				fmt.Printf("      %s\n", source.URL)
			}
			fmt.Println()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

