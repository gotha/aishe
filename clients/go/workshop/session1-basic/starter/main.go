package main

import (
	"fmt"
	"log"
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
	
	// TODO: Print the health status
	// Expected output:
	// Server Status: healthy
	// Ollama Accessible: true
	
	fmt.Println()

	// TODO: Ask a question
	question := "What is the capital of France?"
	fmt.Printf("Asking: %s\n", question)
	
	// Hint: Call client.AskQuestion(question) and handle the error
	
	// TODO: Print the answer and sources
	// Expected output:
	// Answer: <the answer>
	// Processing Time: <time> seconds
	// Sources:
	//   [1] <title>
	//       <url>
	
	log.Println("Workshop Session 1 Complete!")
}

