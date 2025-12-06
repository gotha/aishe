package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// Create a new Redis cache client
	opts := ClientOptions{
		BaseURL:   "", // Will use DefaultAISHEURL constant
		Timeout:   120 * time.Second,
		RedisAddr: "", // Will use DefaultRedisAddr constant
		CacheTTL:  1 * time.Hour,
	}

	client, err := NewClient(opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Check server health
	fmt.Println("Checking server health...")
	health, err := client.CheckHealth()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Printf("Server Status: %s\n", health.Status)
	fmt.Printf("Ollama Accessible: %v\n", health.OllamaAccessible)
	fmt.Println()

	// Interactive question loop
	fmt.Println("=== AISHE Question Answering (Session 2: Redis Cache) ===")
	fmt.Println("Type your question and press Enter. Type 'exit' to quit.")
	fmt.Println("Note: Exact same questions will be cached!")
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

		// Check if question is in cache
		cacheKey := generateCacheKey(question)
		ctx := context.Background()
		cachedAnswer, err := client.redisClient.Get(ctx, cacheKey).Result()
		isCached := err == nil && cachedAnswer != ""

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

		// Show source (Redis cache or AISHE API)
		if isCached {
			fmt.Printf("Source: Redis Cache ✅ (cache hit)\n")
		} else {
			fmt.Printf("Source: AISHE API (cache miss, now cached)\n")
		}
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

// generateCacheKey creates a SHA-256 hash of the question for use as a cache key
func generateCacheKey(question string) string {
	hash := sha256.Sum256([]byte(question))
	return "aishe:question:" + hex.EncodeToString(hash[:])
}

