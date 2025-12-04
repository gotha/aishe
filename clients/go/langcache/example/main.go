package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ndyakov/aishe/clients/go/langcache"
)

func main() {
	// Create client with LangCache support
	// You need to set these environment variables:
	// - LANGCACHE_URL: LangCache API base URL (e.g., https://api.langcache.redis.io)
	// - LANGCACHE_API_KEY: Your LangCache API key
	// - LANGCACHE_CACHE_ID: Your LangCache cache ID
	client, err := langcache.NewClient(langcache.ClientOptions{
		BaseURL:             "http://localhost:8000",
		SimilarityThreshold: 0.9, // Semantic similarity threshold (0.0-1.0)
	})
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
	fmt.Printf("Server status: %s (Ollama accessible: %v)\n\n", health.Status, health.OllamaAccessible)

	// Ask a question (first time - will be cached in LangCache)
	question := "What is the capital of France?"
	fmt.Printf("Asking question (first time): %s\n", question)
	start := time.Now()
	answer, err := client.AskQuestion(question)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Answer: %s\n", answer.Answer)
	fmt.Printf("Processing time: %.2fs (API: %.2fs)\n", elapsed.Seconds(), answer.ProcessingTime)
	fmt.Printf("Sources: %d\n\n", len(answer.Sources))

	// Ask the same question again (should be faster - from LangCache)
	fmt.Printf("Asking the same question again (from LangCache)...\n")
	start = time.Now()
	answer2, err := client.AskQuestion(question)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed2 := time.Since(start)
	fmt.Printf("Answer: %s\n", answer2.Answer)
	fmt.Printf("Processing time: %.2fs (cached!)\n", elapsed2.Seconds())
	fmt.Printf("Speedup: %.2fx faster\n\n", elapsed.Seconds()/elapsed2.Seconds())

	// Ask a semantically similar question (should also hit cache)
	similarQuestion := "What's the capital city of France?"
	fmt.Printf("Asking semantically similar question: %s\n", similarQuestion)
	start = time.Now()
	answer3, err := client.AskQuestion(similarQuestion)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed3 := time.Since(start)
	fmt.Printf("Answer: %s\n", answer3.Answer)
	fmt.Printf("Processing time: %.2fs\n", elapsed3.Seconds())
	if elapsed3 < elapsed {
		fmt.Printf("Cache hit! Speedup: %.2fx faster\n\n", elapsed.Seconds()/elapsed3.Seconds())
	} else {
		fmt.Println("Cache miss (similarity below threshold)\n")
	}

	// Flush cache (optional)
	// fmt.Println("Flushing LangCache...")
	// if err := client.FlushCache(); err != nil {
	// 	log.Printf("Failed to flush cache: %v", err)
	// } else {
	// 	fmt.Println("Cache flushed successfully")
	// }
}

