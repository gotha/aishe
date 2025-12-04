package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// Check if benchmark mode is requested
	if len(os.Args) > 1 && os.Args[1] == "benchmark" {
		runBenchmark()
		return
	}

	// Create a new LangCache client
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

	// Check server health
	fmt.Println("Checking server health...")
	health, err := client.CheckHealth()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Printf("Server Status: %s\n", health.Status)
	fmt.Printf("Ollama Accessible: %v\n", health.OllamaAccessible)
	fmt.Println()

	fmt.Println("=== Demonstrating Semantic Caching with LangCache ===")
	fmt.Println()

	// First question (cache miss)
	question1 := "What is the capital of France?"
	fmt.Printf("Question 1: %s\n", question1)
	start := time.Now()
	answer1, err := client.AskQuestion(question1)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed1 := time.Since(start).Seconds()
	
	fmt.Printf("Answer: %s\n", answer1.Answer)
	fmt.Printf("Time taken: %.2f seconds\n", elapsed1)
	fmt.Println()
	
	// Similar question (should be cache hit with semantic matching!)
	question2 := "What is the capital city of France?"
	fmt.Printf("Question 2 (semantically similar): %s\n", question2)
	start = time.Now()
	answer2, err := client.AskQuestion(question2)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed2 := time.Since(start).Seconds()
	
	fmt.Printf("Answer: %s\n", answer2.Answer)
	fmt.Printf("Time taken: %.2f seconds (semantic cache hit!)\n", elapsed2)
	if elapsed2 < 0.1 {
		fmt.Println("✅ Cache hit confirmed - response was instant!")
	}
	fmt.Println()
	
	// Another similar question
	question3 := "Tell me the capital of France"
	fmt.Printf("Question 3 (also similar): %s\n", question3)
	start = time.Now()
	answer3, err := client.AskQuestion(question3)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed3 := time.Since(start).Seconds()
	
	fmt.Printf("Answer: %s\n", answer3.Answer)
	fmt.Printf("Time taken: %.2f seconds\n", elapsed3)
	if elapsed3 < 0.1 {
		fmt.Println("✅ Cache hit confirmed - response was instant!")
	}
	fmt.Println()
	
	// Different question (cache miss)
	question4 := "What is the capital of Germany?"
	fmt.Printf("Question 4 (different topic): %s\n", question4)
	start = time.Now()
	answer4, err := client.AskQuestion(question4)
	if err != nil {
		log.Fatalf("Failed to ask question: %v", err)
	}
	elapsed4 := time.Since(start).Seconds()
	
	fmt.Printf("Answer: %s\n", answer4.Answer)
	fmt.Printf("Time taken: %.2f seconds\n", elapsed4)
	fmt.Println()
	
	// Flush the cache
	fmt.Println("Flushing LangCache...")
	if err := client.FlushCache(); err != nil {
		log.Fatalf("Failed to flush cache: %v", err)
	}
	fmt.Println("Cache flushed successfully!")
	fmt.Println()
	
	log.Println("Workshop Session 3 Complete!")
	
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🎯 Key Takeaways:")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()
	fmt.Println("1. SEMANTIC CACHING:")
	fmt.Println("   LangCache uses AI embeddings to understand question similarity.")
	fmt.Println("   These questions all matched the same cache entry:")
	fmt.Println("   - 'What is the capital of France?'")
	fmt.Println("   - 'What is the capital city of France?'")
	fmt.Println("   - 'Tell me the capital of France'")
	fmt.Println()
	fmt.Println("2. PERFORMANCE:")
	fmt.Printf("   - First request (no cache):  %.2fs\n", elapsed1)
	fmt.Printf("   - Semantic match 1:          %.2fs (%.0fx faster)\n", elapsed2, elapsed1/elapsed2)
	fmt.Printf("   - Semantic match 2:          %.2fs (%.0fx faster)\n", elapsed3, elapsed1/elapsed3)
	fmt.Printf("   - Different question:        %.2fs (cache miss)\n", elapsed4)
	fmt.Println()
	fmt.Println("3. COMPARISON WITH SESSION 2 (Redis Cache):")
	fmt.Println("   - Redis Cache: Exact match only (SHA-256 hash)")
	fmt.Println("     'What is the capital of France?' ≠ 'What is the capital city of France?'")
	fmt.Println("   - LangCache: Semantic similarity matching")
	fmt.Println("     'What is the capital of France?' ≈ 'What is the capital city of France?'")
	fmt.Println()
	fmt.Println("4. USE CASES:")
	fmt.Println("   - Basic Client: Simple, no dependencies, good for testing")
	fmt.Println("   - Redis Cache: Fast exact-match caching, self-hosted")
	fmt.Println("   - LangCache: Intelligent semantic caching, managed service")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
}

