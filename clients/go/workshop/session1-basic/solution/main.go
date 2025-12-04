package main

import (
	"fmt"
	"log"
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

	// Ask a question
	question := "What is the capital of France?"
	fmt.Printf("Asking: %s\n", question)
	
	answer, err := client.AskQuestion(question)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}

	fmt.Printf("\nAnswer: %s\n", answer.Answer)
	fmt.Printf("Processing Time: %.2f seconds\n", answer.ProcessingTime)
	
	if len(answer.Sources) > 0 {
		fmt.Println("\nSources:")
		for _, source := range answer.Sources {
			fmt.Printf("  [%d] %s\n", source.Number, source.Title)
			fmt.Printf("      %s\n", source.URL)
		}
	}
}

