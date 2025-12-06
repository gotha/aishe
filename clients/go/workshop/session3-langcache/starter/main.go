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
	// TODO: Create a new LangCache client
	// Note: You'll need to update the LangCache constants in client.go:
	// - LangCacheURL
	// - LangCacheAPIKey
	// - LangCacheCacheID
	opts := ClientOptions{
		BaseURL:             "", // Will use DefaultAISHEURL constant
		Timeout:             120 * time.Second,
		SimilarityThreshold: 0.9, // 90% similarity threshold
	}

	client, err := NewClient(opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// TODO: Check server health
	fmt.Println("Checking server health...")
	health, err := client.CheckHealth()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Printf("Server Status: %s\n", health.Status)
	fmt.Printf("Ollama Accessible: %v\n", health.OllamaAccessible)
	fmt.Println()

	// Interactive question loop
	fmt.Println("=== AISHE Question Answering (Session 3: LangCache Semantic Caching) ===")
	fmt.Println("Type your question and press Enter. Type 'exit' to quit.")
	fmt.Println("Note: Semantically similar questions will be cached!")
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

		// TODO: Check if question is in LangCache
		// Hint: Call client.checkCache(question) to check for semantic match
		// isCached, err := client.checkCache(question)
		// if err != nil {
		//     isCached = false // Treat errors as cache miss
		// }
		isCached := false // TODO: Replace with actual cache check

		// TODO: Measure execution time and ask the question
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
		// Source: LangCache ✅ (semantic cache hit)  OR  Source: AISHE API (cache miss, now cached)
		// Execution Time: X.XX seconds
		//
		// Wikipedia Sources:
		//   [1] <title>
		//       <url>
		//

		fmt.Println()
		fmt.Println("TODO: Display the answer here")
		fmt.Println()

		// Show source (LangCache or AISHE API)
		if isCached {
			fmt.Printf("Source: LangCache ✅ (semantic cache hit)\n")
		} else {
			fmt.Printf("Source: AISHE API (cache miss, now cached)\n")
		}
		fmt.Printf("Execution Time: %.2f seconds\n", executionTime)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

