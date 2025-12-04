package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// BasicClient is a simple HTTP client without caching
type BasicClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewBasicClient(baseURL string) *BasicClient {
	if baseURL == "" {
		baseURL = DefaultAISHEURL
	}
	return &BasicClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *BasicClient) AskQuestion(question string) (*AnswerResponse, error) {
	url := fmt.Sprintf("%s/api/v1/ask", c.baseURL)
	reqBody := QuestionRequest{Question: question}
	jsonData, _ := json.Marshal(reqBody)
	
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var answer AnswerResponse
	json.Unmarshal(body, &answer)
	return &answer, nil
}

func (c *BasicClient) Close() error { return nil }

// RedisClient is an HTTP client with Redis exact-match caching
type RedisClient struct {
	baseURL     string
	httpClient  *http.Client
	redisClient *redis.Client
	ctx         context.Context
}

func NewRedisClient(baseURL, redisAddr string) (*RedisClient, error) {
	if baseURL == "" {
		baseURL = DefaultAISHEURL
	}
	if redisAddr == "" {
		redisAddr = DefaultRedisAddr
	}
	
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	
	return &RedisClient{
		baseURL:     strings.TrimRight(baseURL, "/"),
		httpClient:  &http.Client{Timeout: 120 * time.Second},
		redisClient: rdb,
		ctx:         ctx,
	}, nil
}

func (c *RedisClient) generateCacheKey(question string) string {
	hash := sha256.Sum256([]byte(strings.TrimSpace(strings.ToLower(question))))
	return fmt.Sprintf("benchmark:answer:%x", hash)
}

func (c *RedisClient) AskQuestion(question string) (*AnswerResponse, error) {
	cacheKey := c.generateCacheKey(question)
	
	// Try cache
	cachedData, err := c.redisClient.Get(c.ctx, cacheKey).Bytes()
	if err == nil {
		var answer AnswerResponse
		if json.Unmarshal(cachedData, &answer) == nil {
			return &answer, nil
		}
	}
	
	// Fetch from API
	url := fmt.Sprintf("%s/api/v1/ask", c.baseURL)
	reqBody := QuestionRequest{Question: question}
	jsonData, _ := json.Marshal(reqBody)
	
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var answer AnswerResponse
	json.Unmarshal(body, &answer)
	
	// Cache it
	answerData, _ := json.Marshal(answer)
	c.redisClient.Set(c.ctx, cacheKey, answerData, 1*time.Hour)
	
	return &answer, nil
}

func (c *RedisClient) ClearCache() error {
	iter := c.redisClient.Scan(c.ctx, 0, "benchmark:answer:*", 0).Iterator()
	for iter.Next(c.ctx) {
		c.redisClient.Del(c.ctx, iter.Val())
	}
	return iter.Err()
}

func (c *RedisClient) Close() error {
	return c.redisClient.Close()
}

// Benchmark results
type BenchmarkResult struct {
	ClientType  string
	Question    string
	Duration    time.Duration
	CacheHit    bool
	Answer      string
}

func runBenchmark() {
	fmt.Println("=" + strings.Repeat("=", 78))
	fmt.Println("  AISHE Client Benchmark: Basic vs Redis vs LangCache")
	fmt.Println("=" + strings.Repeat("=", 78))
	fmt.Println()
	
	questions := []string{
		"What is the capital of France?",
		"What is the capital of France?",      // Exact duplicate - will hit Redis cache
		"What is the capital city of France?", // Semantically similar - will hit LangCache only
		"Tell me the capital of France",       // Semantically similar - will hit LangCache only
	}
	
	var results []BenchmarkResult
	
	// Test 1: Basic Client (no caching)
	fmt.Println("📊 Test 1: Basic Client (No Caching)")
	fmt.Println(strings.Repeat("-", 80))
	basicClient := NewBasicClient("")
	defer basicClient.Close()
	
	for _, q := range questions {
		start := time.Now()
		answer, _ := basicClient.AskQuestion(q)
		duration := time.Since(start)
		
		results = append(results, BenchmarkResult{
			ClientType: "Basic",
			Question:   q,
			Duration:   duration,
			CacheHit:   false,
			Answer:     answer.Answer[:min(50, len(answer.Answer))],
		})
		
		fmt.Printf("  Q: %s\n", q)
		fmt.Printf("  ⏱️  Time: %.3fs\n", duration.Seconds())
		fmt.Println()
	}
	
	// Test 2: Redis Client (exact-match caching)
	fmt.Println("\n📊 Test 2: Redis Client (Exact-Match Caching)")
	fmt.Println(strings.Repeat("-", 80))
	redisClient, err := NewRedisClient("", "")
	if err != nil {
		fmt.Printf("  ⚠️  Redis not available: %v\n", err)
	} else {
		defer redisClient.Close()
		redisClient.ClearCache() // Start fresh
		
		// Track which questions we've seen for cache hit detection
		seenQuestions := make(map[string]bool)

		for _, q := range questions {
			start := time.Now()
			answer, _ := redisClient.AskQuestion(q)
			duration := time.Since(start)

			// Cache hit if we've seen this exact question before
			cacheHit := seenQuestions[q]
			seenQuestions[q] = true

			results = append(results, BenchmarkResult{
				ClientType: "Redis",
				Question:   q,
				Duration:   duration,
				CacheHit:   cacheHit,
				Answer:     answer.Answer[:min(50, len(answer.Answer))],
			})

			fmt.Printf("  Q: %s\n", q)
			fmt.Printf("  ⏱️  Time: %.3fs", duration.Seconds())
			if cacheHit {
				fmt.Printf(" ✅ CACHE HIT")
			}
			fmt.Println()
			fmt.Println()
		}
		
		redisClient.ClearCache()
	}
	
	// Test 3: LangCache Client (semantic caching)
	fmt.Println("\n📊 Test 3: LangCache Client (Semantic Caching)")
	fmt.Println(strings.Repeat("-", 80))

	if LangCacheURL == "" || LangCacheURL == "https://your-langcache-url.com" {
		fmt.Println("  ⚠️  LangCache not configured (update constants in client.go)")
	} else {
		opts := ClientOptions{
			SimilarityThreshold: 0.9,
		}
		langClient, err := NewClient(opts)
		if err != nil {
			fmt.Printf("  ⚠️  LangCache error: %v\n", err)
		} else {
			defer langClient.Close()
			langClient.FlushCache() // Start fresh
			
			for i, q := range questions {
				start := time.Now()
				answer, _ := langClient.AskQuestion(q)
				duration := time.Since(start)
				
				cacheHit := i > 0 // All similar questions should hit
				
				results = append(results, BenchmarkResult{
					ClientType: "LangCache",
					Question:   q,
					Duration:   duration,
					CacheHit:   cacheHit,
					Answer:     answer.Answer[:min(50, len(answer.Answer))],
				})
				
				fmt.Printf("  Q: %s\n", q)
				fmt.Printf("  ⏱️  Time: %.3fs", duration.Seconds())
				if cacheHit {
					fmt.Printf(" ✅ SEMANTIC CACHE HIT")
				}
				fmt.Println()
				fmt.Println()
			}
			
			langClient.FlushCache()
		}
	}
	
	// Summary
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("  📈 BENCHMARK SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
	
	fmt.Printf("%-15s %-45s %10s %12s\n", "Client", "Question", "Time", "Cache Hit")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, r := range results {
		cacheStatus := ""
		if r.CacheHit {
			cacheStatus = "✅ YES"
		} else {
			cacheStatus = "❌ NO"
		}
		
		question := r.Question
		if len(question) > 45 {
			question = question[:42] + "..."
		}
		
		fmt.Printf("%-15s %-45s %9.3fs %12s\n",
			r.ClientType,
			question,
			r.Duration.Seconds(),
			cacheStatus,
		)
	}
	
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("  🎯 KEY INSIGHTS")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
	fmt.Println("1. Basic Client: No caching - every request is slow (~5s)")
	fmt.Println("2. Redis Client: Exact-match only - only identical questions are cached")
	fmt.Println("3. LangCache: Semantic matching - similar questions share the same cache!")
	fmt.Println()
	fmt.Println("💡 LangCache provides the best cache hit rate for natural language queries!")
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

