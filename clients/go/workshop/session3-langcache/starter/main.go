package main

import (
	"fmt"
	"log"
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
	// Hint: Same as previous sessions
	
	fmt.Println()

	// TODO: Demonstrate semantic caching with similar questions
	
	// First question (cache miss)
	question1 := "What is the capital of France?"
	fmt.Printf("Question 1: %s\n", question1)
	start := time.Now()
	// TODO: Call AskQuestion
	
	fmt.Printf("Answer: %s\n", "TODO")
	fmt.Printf("Time taken: %.2f seconds\n", time.Since(start).Seconds())
	fmt.Println()
	
	// Similar question (should be cache hit with semantic matching!)
	question2 := "What is the capital city of France?"
	fmt.Printf("Question 2 (semantically similar): %s\n", question2)
	start = time.Now()
	// TODO: Call AskQuestion
	
	fmt.Printf("Answer: %s\n", "TODO")
	fmt.Printf("Time taken: %.2f seconds (semantic cache hit!)\n", time.Since(start).Seconds())
	fmt.Println()
	
	// Another similar question
	question3 := "Tell me the capital of France"
	fmt.Printf("Question 3 (also similar): %s\n", question3)
	start = time.Now()
	// TODO: Call AskQuestion
	
	fmt.Printf("Answer: %s\n", "TODO")
	fmt.Printf("Time taken: %.2f seconds\n", time.Since(start).Seconds())
	fmt.Println()
	
	// Different question (cache miss)
	question4 := "What is the capital of Germany?"
	fmt.Printf("Question 4 (different topic): %s\n", question4)
	start = time.Now()
	// TODO: Call AskQuestion
	
	fmt.Printf("Answer: %s\n", "TODO")
	fmt.Printf("Time taken: %.2f seconds\n", time.Since(start).Seconds())
	fmt.Println()
	
	// TODO: Flush the cache
	fmt.Println("Flushing LangCache...")
	// Hint: Call client.FlushCache()
	
	log.Println("Workshop Session 3 Complete!")
	fmt.Println("\n🎯 Key Takeaway:")
	fmt.Println("LangCache uses semantic similarity to match questions,")
	fmt.Println("so 'What is the capital of France?' and 'What is the capital city of France?'")
	fmt.Println("are treated as the same question!")
}

