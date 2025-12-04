# AISHE Go Client Workshop

Welcome to the AISHE Go Client Workshop! This hands-on workshop will teach you how to build HTTP clients in Go with different caching strategies.

## 🎯 Workshop Overview

This workshop is divided into three progressive sessions:

1. **Session 1: Basic HTTP Client** - Learn the fundamentals of building HTTP clients in Go
2. **Session 2: Redis Cache Client** - Add exact-match caching with Redis
3. **Session 3: LangCache Client** - Implement semantic caching with Redis LangCache

Each session builds on the previous one, introducing new concepts and techniques.

## 📋 Prerequisites

### Required
- Go 1.21 or higher
- Docker and Docker Compose (for running AISHE server)
- Basic understanding of Go programming
- Basic understanding of HTTP and REST APIs

### Recommended
- Familiarity with JSON in Go
- Basic understanding of caching concepts
- Text editor or IDE with Go support (VS Code, GoLand, etc.)

## 🚀 Setup

### 1. Start the AISHE Server

First, start the AISHE server with Docker Compose:

```bash
# From the repository root
./docker-setup.sh
```

This will start:
- AISHE API server on http://localhost:8000
- Ollama LLM service on http://localhost:11434
- Redis on localhost:6379

Verify everything is running:

```bash
./test-docker.sh
```

### 2. Configure Constants

The workshop uses constants in the code instead of environment variables for easier configuration.

**For Sessions 1 and 2:** No configuration needed! The defaults work with the Docker setup.

**For Session 3 (LangCache):** You'll need to update the constants in `client.go`:

1. Sign up for Redis LangCache at https://redis.io/langcache/
2. Get your credentials (URL, API Key, Cache ID)
3. Open `session3-langcache/starter/client.go` (or `solution/client.go`)
4. Update these constants at the top of the file:
   ```go
   const (
       LangCacheURL    = "https://your-actual-url.com"
       LangCacheAPIKey = "your-actual-api-key"
       LangCacheCacheID = "your-actual-cache-id"
   )
   ```

**Note:** If your AISHE server or Redis are running on different ports, you can also update `DefaultAISHEURL` and `DefaultRedisAddr` constants in the respective `client.go` files.

## 📚 Workshop Sessions

### Session 1: Basic HTTP Client

**Duration:** 30-45 minutes

**Learning Objectives:**
- Understand HTTP client basics in Go
- Work with JSON encoding/decoding
- Handle errors properly
- Make GET and POST requests

**Files:**
- `session1-basic/starter/` - Your starting point
- `session1-basic/solution/` - Reference solution

**Tasks:**
1. Implement `NewClient()` - Create a client with configuration
2. Implement `CheckHealth()` - Make a GET request to check server health
3. Implement `AskQuestion()` - Make a POST request to ask questions

**Run the starter:**

```bash
cd session1-basic/starter
go run .
```

**Key Concepts:**
- `http.Client` and `http.Request`
- JSON marshaling/unmarshaling
- Error handling patterns
- Configuration with constants

---

### Session 2: Redis Cache Client

**Duration:** 45-60 minutes

**Learning Objectives:**
- Integrate Redis for caching
- Implement cache-aside pattern
- Generate cache keys with hashing
- Measure performance improvements

**Files:**
- `session2-redis/starter/` - Your starting point
- `session2-redis/solution/` - Reference solution

**Tasks:**
1. Implement `NewClient()` - Create client with Redis connection
2. Implement `generateCacheKey()` - Create SHA-256 hash-based cache keys
3. Implement `fetchFromAPI()` - Separate API fetching logic
4. Implement `AskQuestion()` - Add caching logic (check cache → fetch → store)
5. Implement `ClearCache()` - Clear all cached entries

**Run the starter:**

```bash
cd session2-redis/starter
# go.mod already exists with dependencies
go run .
```

**Key Concepts:**
- Redis client library (go-redis)
- Cache-aside pattern
- SHA-256 hashing for cache keys
- TTL (Time To Live) for cache entries
- Performance measurement

**Expected Results:**
- First request: ~5 seconds (API call)
- Second request: <0.1 seconds (cache hit)
- Speedup: ~50-100x faster!

---

### Session 3: LangCache Client

**Duration:** 60-75 minutes

**Learning Objectives:**
- Understand semantic caching vs exact-match caching
- Integrate with external REST APIs
- Implement Bearer token authentication
- Compare different caching strategies

**Files:**
- `session3-langcache/starter/` - Your starting point
- `session3-langcache/solution/` - Reference solution

**Tasks:**
1. Implement `NewClient()` - Configure LangCache credentials
2. Implement `searchLangCache()` - Search for semantically similar questions
3. Implement `storeLangCache()` - Store responses in LangCache
4. Implement `fetchFromAPI()` - Fetch from AISHE API
5. Implement `AskQuestion()` - Orchestrate semantic caching
6. Implement `FlushCache()` - Clear the cache

**Run the starter:**

```bash
cd session3-langcache/starter
# First, update the LangCache constants in client.go
# Then run:
go run .
```

**Key Concepts:**
- Semantic similarity vs exact matching
- REST API integration
- Bearer token authentication
- AI embeddings for caching
- Trade-offs between caching strategies

**Expected Results:**
- "What is the capital of France?" → Cache miss
- "What is the capital city of France?" → **Cache hit!** (semantic match)
- "Tell me the capital of France" → **Cache hit!** (semantic match)

---

## 🎓 Learning Path

### Recommended Order

1. **Start with Session 1** - Build a solid foundation
2. **Progress to Session 2** - Learn caching fundamentals
3. **Complete Session 3** - Explore advanced semantic caching

### Tips for Success

1. **Read the TODOs carefully** - Each TODO has hints and guidance
2. **Test frequently** - Run your code after each function implementation
3. **Compare with solutions** - If stuck, peek at the solution for guidance
4. **Experiment** - Try different questions, cache keys, thresholds
5. **Ask questions** - Don't hesitate to ask the instructor

### Common Pitfalls

1. **Forgetting to close resources** - Always `defer client.Close()`
2. **Not checking errors** - Go requires explicit error handling
3. **Incorrect JSON tags** - Make sure struct tags match API response
4. **Missing configuration** - Update constants in `client.go` for Session 3
5. **Not trimming/normalizing input** - Always clean user input

## 📊 Performance Comparison

After completing all three sessions, you'll understand the trade-offs:

| Approach | Speed | Cache Hit Rate | Complexity | Cost |
|----------|-------|----------------|------------|------|
| **Basic Client** | Slow (5s) | N/A | Low | Free |
| **Redis Cache** | Fast (0.05s) | Low (exact match) | Medium | Low (self-hosted) |
| **LangCache** | Fast (0.05s) | High (semantic) | Medium | Medium (managed) |

### When to Use Each

- **Basic Client**: Testing, development, low-traffic applications
- **Redis Cache**: High-traffic with repeated exact queries, self-hosted infrastructure
- **LangCache**: High-traffic with varied phrasing, managed service preferred

## 🔧 Troubleshooting

### AISHE Server Not Running

```bash
# Check if containers are running
docker-compose ps

# Restart if needed
docker-compose down
./docker-setup.sh
```

### Redis Connection Failed

```bash
# Check Redis is running
docker exec aishe-redis redis-cli ping
# Should return: PONG

# Check the constant in client.go
# DefaultRedisAddr should be: "localhost:6379"
```

### LangCache Authentication Failed

```bash
# Verify constants are set in client.go
# Open session3-langcache/starter/client.go or solution/client.go
# Check that these constants are updated:
# - LangCacheURL
# - LangCacheAPIKey
# - LangCacheCacheID

# Test with curl (replace with your actual values)
curl -H "Authorization: Bearer YOUR_API_KEY" \
  YOUR_LANGCACHE_URL/v1/caches/YOUR_CACHE_ID/entries/search \
  -d '{"prompt":"test"}'
```

### Import Errors

```bash
# Make sure you've initialized the module
go mod init workshop-sessionX

# Download dependencies
go mod tidy
```

## 🎉 Completion

After completing all three sessions, you will have:

✅ Built three different HTTP clients in Go  
✅ Implemented exact-match caching with Redis  
✅ Implemented semantic caching with LangCache  
✅ Understood caching patterns and trade-offs  
✅ Gained hands-on experience with Go HTTP clients  
✅ Learned to integrate with external APIs  

## 📚 Additional Resources

- [Go HTTP Client Documentation](https://pkg.go.dev/net/http)
- [go-redis Documentation](https://redis.uptrace.dev/)
- [Redis LangCache Documentation](https://redis.io/docs/latest/develop/ai/langcache/)
- [AISHE API Documentation](http://localhost:8000/docs) (when server is running)

## 🤝 Contributing

Found an issue or have a suggestion? Please open an issue or submit a pull request!

## 📝 License

This workshop is part of the AISHE project.

---

**Happy Coding! 🚀**

