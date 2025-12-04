package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// TODO: Create a new Redis cache client
	// Hint: Use ClientOptions struct with default values
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

	// TODO: Check server health
	fmt.Println("Checking server health...")
	// Hint: Same as Session 1
	
	fmt.Println()

	// TODO: Ask the same question twice to demonstrate caching
	question := "What is the capital of France?"
	
	// First request (cache miss - will be slow)
	fmt.Printf("First request: %s\n", question)
	start := time.Now()
	// TODO: Call AskQuestion and measure time
	
	fmt.Printf("Answer: %s\n", "TODO")
	fmt.Printf("Time taken: %.2f seconds\n", time.Since(start).Seconds())
	fmt.Println()
	
	// Second request (cache hit - should be fast!)
	fmt.Printf("Second request (cached): %s\n", question)
	start = time.Now()
	// TODO: Call AskQuestion again and measure time
	
	fmt.Printf("Answer: %s\n", "TODO")
	fmt.Printf("Time taken: %.2f seconds (from cache!)\n", time.Since(start).Seconds())
	fmt.Println()
	
	// TODO: Try a slightly different question (should be cache miss)
	question2 := "What is the capital city of France?"
	fmt.Printf("Different question: %s\n", question2)
	start = time.Now()
	// TODO: Call AskQuestion with question2
	
	fmt.Printf("Answer: %s\n", "TODO")
	fmt.Printf("Time taken: %.2f seconds\n", time.Since(start).Seconds())
	fmt.Println()
	
	// TODO: Clear the cache
	fmt.Println("Clearing cache...")
	// Hint: Call client.ClearCache()
	
	log.Println("Workshop Session 2 Complete!")
}

