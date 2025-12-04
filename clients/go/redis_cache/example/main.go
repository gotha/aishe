package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ndyakov/aishe/clients/go/redis_cache"
)

func main() {
	// Create a new client with Redis caching
	client, err := redis_cache.NewClient(redis_cache.ClientOptions{
		BaseURL:   "",                // Uses AISHE_API_URL env or default
		Timeout:   120 * time.Second,
		RedisAddr: "localhost:6379",  // Uses REDIS_ADDR env or default
		RedisDB:   0,
		CacheTTL:  1 * time.Hour,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Check health
	fmt.Println("Checking server health...")
	health, err := client.CheckHealth()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Printf("Server Status: %s\n", health.Status)
	fmt.Printf("Ollama Accessible: %v\n", health.OllamaAccessible)
	fmt.Println()

	// Ask a question (first time - will hit API and cache)
	question := "What is the capital of France?"
	fmt.Printf("Asking (first time): %s\n", question)
	
	start := time.Now()
	answer, err := client.AskQuestion(question)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	firstDuration := time.Since(start)

	fmt.Printf("\nAnswer: %s\n", answer.Answer)
	fmt.Printf("Processing Time: %.2f seconds\n", answer.ProcessingTime)
	fmt.Printf("Total Time (including network): %.2f seconds\n", firstDuration.Seconds())
	
	// Ask the same question again (should be cached)
	fmt.Printf("\n\nAsking (second time - from cache): %s\n", question)
	
	start = time.Now()
	answer2, err := client.AskQuestion(question)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	secondDuration := time.Since(start)

	fmt.Printf("\nAnswer: %s\n", answer2.Answer)
	fmt.Printf("Total Time (from cache): %.2f seconds\n", secondDuration.Seconds())
	fmt.Printf("Speedup: %.2fx faster\n", firstDuration.Seconds()/secondDuration.Seconds())

	if len(answer.Sources) > 0 {
		fmt.Println("\nSources:")
		for _, source := range answer.Sources {
			fmt.Printf("  [%d] %s\n", source.Number, source.Title)
			fmt.Printf("      %s\n", source.URL)
		}
	}

	// Optional: Clear cache
	// fmt.Println("\nClearing cache...")
	// if err := client.ClearCache(); err != nil {
	// 	log.Printf("Failed to clear cache: %v", err)
	// }
}

