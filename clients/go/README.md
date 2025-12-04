# AISHE Go Clients

Three Go client implementations for the AISHE (AI Search & Help Engine) HTTP API.

## Clients

### 1. Basic HTTP Client (`basic/`)

A simple HTTP client without caching.

**Features:**
- Direct HTTP requests to the AISHE API
- Health check endpoint
- Question answering endpoint
- No external dependencies (except standard library)

**Usage:**

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/ndyakov/aishe/clients/go/basic"
)

func main() {
    client := basic.NewClient("http://localhost:8000", 120*time.Second)
    defer client.Close()
    
    answer, err := client.AskQuestion("What is the capital of France?")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(answer.Answer)
}
```

### 2. Redis Cache Client (`redis_cache/`)

HTTP client with basic Redis caching using go-redis.

**Features:**
- All features from basic client
- Redis-based response caching
- Configurable TTL
- Cache key generation using SHA-256
- Cache clearing functionality

**Usage:**

```go
package main

import (
    "log"
    "time"
    
    "github.com/ndyakov/aishe/clients/go/redis_cache"
)

func main() {
    client, err := redis_cache.NewClient(redis_cache.ClientOptions{
        BaseURL:   "http://localhost:8000",
        Timeout:   120 * time.Second,
        RedisAddr: "localhost:6379",
        RedisDB:   0,
        CacheTTL:  1 * time.Hour,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // First call hits API and caches result
    answer, _ := client.AskQuestion("What is the capital of France?")
    
    // Second call returns cached result (much faster)
    answer, _ = client.AskQuestion("What is the capital of France?")
}
```

### 3. Redis LangCache Client (`langcache/`)

HTTP client with Redis LangCache integration for semantic caching.

**What is LangCache?**

[Redis LangCache](https://redis.io/langcache/) is a fully-managed semantic caching service available via REST API. Unlike traditional exact-match caching, LangCache uses semantic similarity to match queries, meaning similar questions can retrieve the same cached response.

**Features:**
- Semantic caching via LangCache REST API
- Similarity threshold configuration
- Automatic cache storage after API calls
- Cache flushing capability
- Works with any LangCache-compatible service

**Requirements:**
- LangCache API URL (e.g., `https://api.langcache.redis.io`)
- LangCache API key
- LangCache cache ID

**Usage:**

```go
package main

import (
    "log"
    "os"

    "github.com/ndyakov/aishe/clients/go/langcache"
)

func main() {
    // Set environment variables:
    // export LANGCACHE_URL="https://api.langcache.redis.io"
    // export LANGCACHE_API_KEY="your-api-key"
    // export LANGCACHE_CACHE_ID="your-cache-id"

    client, err := langcache.NewClient(langcache.ClientOptions{
        BaseURL:             "http://localhost:8000",
        SimilarityThreshold: 0.9, // 0.0-1.0, higher = stricter matching
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Ask question (searches LangCache first, then API if miss)
    answer, _ := client.AskQuestion("What is the capital of France?")

    // Semantically similar questions may hit the cache
    answer, _ = client.AskQuestion("What's the capital city of France?")
}
```

## Installation

```bash
cd clients/go
go mod download
```

## Running Examples

### Basic Client

```bash
cd clients/go/basic/example
go run main.go
```

### Redis Cache Client

Make sure Redis is running:

```bash
# Start Redis (if not already running)
redis-server

# Run example
cd clients/go/redis_cache/example
go run main.go
```

### LangCache Client

Make sure you have LangCache credentials:

```bash
# Set LangCache environment variables
export LANGCACHE_URL="https://api.langcache.redis.io"
export LANGCACHE_API_KEY="your-api-key"
export LANGCACHE_CACHE_ID="your-cache-id"

# Run example
cd clients/go/langcache/example
go run main.go
```

## Environment Variables

All clients support the following environment variables:

- `AISHE_API_URL`: Base URL of the AISHE API (default: `http://localhost:8000`)
- `REDIS_ADDR`: Redis server address (default: `localhost:6379`) - for Redis cache client only

LangCache client additionally requires:

- `LANGCACHE_URL`: LangCache API base URL (required)
- `LANGCACHE_API_KEY`: LangCache API key (required)
- `LANGCACHE_CACHE_ID`: LangCache cache ID (required)

## API Reference

### Common Types

```go
type Source struct {
    Number int    `json:"number"`
    Title  string `json:"title"`
    URL    string `json:"url"`
}

type AnswerResponse struct {
    Answer         string   `json:"answer"`
    Sources        []Source `json:"sources"`
    ProcessingTime float64  `json:"processing_time"`
}

type HealthResponse struct {
    Status           string  `json:"status"`
    OllamaAccessible bool    `json:"ollama_accessible"`
    Message          *string `json:"message,omitempty"`
}
```

### Client Methods

All clients implement:

- `CheckHealth() (*HealthResponse, error)` - Check API server health
- `AskQuestion(question string) (*AnswerResponse, error)` - Ask a question
- `Close() error` - Close the client and cleanup resources

Redis cache client additionally implements:

- `ClearCache() error` - Clear all cached responses

LangCache client additionally implements:

- `FlushCache() error` - Flush all entries from LangCache

## Dependencies

- Basic client: No external dependencies
- Redis cache client: `github.com/redis/go-redis/v9`
- LangCache client: No external dependencies (uses standard library HTTP client)

## Testing

Make sure the AISHE server is running before testing:

```bash
# In the project root
nix run .#server
```

Then run the examples or your own code.

## Performance Comparison

Typical performance (approximate):

- **Basic client**: ~2-5 seconds per question (depends on LLM processing)
- **Redis cache client**: ~2-5 seconds (first call), ~10-50ms (cached, exact match only)
- **LangCache client**: ~2-5 seconds (first call), ~50-200ms (cached, semantic matching)

The cache clients provide significant speedup for repeated questions. LangCache additionally provides semantic matching, so similar questions can hit the cache even if not exactly the same.

## License

Same as the main AISHE project.

