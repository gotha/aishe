package main

import (
	"fmt"
	"log"
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

	// Ask the same question twice to demonstrate caching
	question := "What is the capital of France?"
	
	// First request (cache miss - will be slow)
	fmt.Printf("First request: %s\n", question)
	start := time.Now()
	answer1, err := client.AskQuestion(question)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed1 := time.Since(start).Seconds()
	
	fmt.Printf("Answer: %s\n", answer1.Answer)
	fmt.Printf("Time taken: %.2f seconds\n", elapsed1)
	fmt.Println()
	
	// Second request (cache hit - should be fast!)
	fmt.Printf("Second request (cached): %s\n", question)
	start = time.Now()
	answer2, err := client.AskQuestion(question)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed2 := time.Since(start).Seconds()
	
	fmt.Printf("Answer: %s\n", answer2.Answer)
	fmt.Printf("Time taken: %.2f seconds (from cache!)\n", elapsed2)
	fmt.Printf("Speedup: %.0fx faster!\n", elapsed1/elapsed2)
	fmt.Println()
	
	// Try a slightly different question (should be cache miss)
	question2 := "What is the capital city of France?"
	fmt.Printf("Different question: %s\n", question2)
	start = time.Now()
	answer3, err := client.AskQuestion(question2)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed3 := time.Since(start).Seconds()
	
	fmt.Printf("Answer: %s\n", answer3.Answer)
	fmt.Printf("Time taken: %.2f seconds\n", elapsed3)
	fmt.Println()
	
	// Clear the cache
	fmt.Println("Clearing cache...")
	if err := client.ClearCache(); err != nil {
		log.Fatalf("Failed to clear cache: %v", err)
	}
	fmt.Println("Cache cleared successfully!")
	fmt.Println()
	
	log.Println("Workshop Session 2 Complete!")
	fmt.Println("\n📊 Performance Summary:")
	fmt.Printf("  - First request (no cache): %.2fs\n", elapsed1)
	fmt.Printf("  - Second request (cached):  %.2fs (%.0fx faster)\n", elapsed2, elapsed1/elapsed2)
	fmt.Printf("  - Different question:       %.2fs\n", elapsed3)
}

