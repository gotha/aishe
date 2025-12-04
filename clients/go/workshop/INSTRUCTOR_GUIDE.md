# AISHE Go Client Workshop - Instructor Guide

## 📋 Workshop Overview

This workshop teaches students how to build HTTP clients in Go with progressive caching strategies. It's designed for a 3-session format, with each session building on the previous one.

**Total Duration:** 2.5 - 3 hours  
**Skill Level:** Intermediate Go developers  
**Prerequisites:** Basic Go knowledge, HTTP/REST understanding

## 🎯 Learning Objectives

By the end of this workshop, students will:

1. Build HTTP clients in Go using the standard library
2. Implement JSON encoding/decoding for API communication
3. Integrate Redis for exact-match caching
4. Understand and implement semantic caching with LangCache
5. Compare different caching strategies and their trade-offs
6. Apply best practices for error handling and resource management

## 📚 Session Breakdown

### Session 1: Basic HTTP Client (30-45 min)

**Objectives:**
- Understand HTTP client fundamentals
- Work with JSON marshaling/unmarshaling
- Implement proper error handling
- Configure clients with constants

**Key Teaching Points:**
1. `http.Client` configuration (timeouts, etc.)
2. Making GET and POST requests
3. Reading and parsing JSON responses
4. Error handling patterns in Go
5. Resource cleanup with `defer`

**Common Student Mistakes:**
- Forgetting to close response bodies
- Not checking errors
- Incorrect JSON struct tags
- Not handling empty/nil responses

**Solution Walkthrough:**
1. Start with `NewClient()` - explain configuration patterns
2. Show `CheckHealth()` - simple GET request
3. Build `AskQuestion()` - POST with JSON body
4. Emphasize error handling at each step

---

### Session 2: Redis Cache Client (45-60 min)

**Objectives:**
- Integrate Redis for caching
- Implement cache-aside pattern
- Generate deterministic cache keys
- Measure performance improvements

**Key Teaching Points:**
1. Cache-aside pattern (check cache → fetch → store)
2. SHA-256 hashing for cache keys
3. Redis client library usage
4. TTL (Time To Live) for cache entries
5. Performance measurement techniques

**Common Student Mistakes:**
- Not normalizing questions before hashing (case, whitespace)
- Forgetting to handle cache misses gracefully
- Not setting TTL on cache entries
- Failing to close Redis connection

**Solution Walkthrough:**
1. Show Redis client initialization and connection testing
2. Explain `generateCacheKey()` - why SHA-256? Why normalize?
3. Walk through `AskQuestion()` caching logic step-by-step
4. Demonstrate the performance difference (first vs second request)
5. Show `ClearCache()` for cleanup

**Demo:**
Run the solution twice with the same question to show:
- First request: ~5 seconds
- Second request: <0.1 seconds
- Speedup: ~50-100x!

---

### Session 3: LangCache Client (60-75 min)

**Objectives:**
- Understand semantic vs exact-match caching
- Integrate with external REST APIs
- Implement Bearer token authentication
- Compare all three caching strategies

**Key Teaching Points:**
1. Semantic similarity vs exact matching
2. REST API integration patterns
3. Bearer token authentication
4. AI embeddings for caching (conceptual)
5. Trade-offs between caching strategies

**Common Student Mistakes:**
- Forgetting Authorization header
- Not handling 404 (cache miss) correctly
- Incorrect JSON structure for LangCache API
- Not understanding similarity threshold

**Solution Walkthrough:**
1. Explain LangCache architecture (managed service, not Redis library)
2. Show `searchLangCache()` - REST API call with auth
3. Explain `storeLangCache()` - storing serialized responses
4. Walk through `AskQuestion()` orchestration
5. Demonstrate semantic matching with similar questions

**Demo:**
Run the solution with these questions:
1. "What is the capital of France?" → Cache miss
2. "What is the capital city of France?" → **Cache hit!** (semantic)
3. "Tell me the capital of France" → **Cache hit!** (semantic)

**Benchmark:**
Run `go run . benchmark` to compare all three approaches side-by-side.

---

## 🛠️ Setup Instructions

### Before the Workshop

1. **Test the Docker environment:**
   ```bash
   cd /path/to/aishe
   ./docker-setup.sh
   ./test-docker.sh
   ```

2. **Verify all solutions work:**
   ```bash
   # Session 1
   cd clients/go/workshop/session1-basic/solution
   go run .
   
   # Session 2
   cd ../session2-redis/solution
   go run .
   
   # Session 3 (requires LangCache credentials)
   cd ../session3-langcache/solution
   # First, update the LangCache constants in client.go
   # Then run:
   go run .
   go run . benchmark
   ```

3. **Prepare LangCache accounts (Session 3):**
   - Option A: Provide shared credentials for all students (update constants in solution files)
   - Option B: Have students sign up individually at https://redis.io/langcache/
   - Option C: Skip Session 3 if LangCache is not available

### Student Setup (15 min before Session 1)

Have students:

1. Clone the repository
2. Start Docker environment:
   ```bash
   ./docker-setup.sh
   ```
3. Verify setup:
   ```bash
   ./test-docker.sh
   ```
4. **No environment variables needed!** Configuration is done via constants in the code.

---

## 🎓 Teaching Tips

### General

1. **Live coding vs starter files:** 
   - Option A: Students fill in TODOs in starter files (recommended)
   - Option B: Live code from scratch with students following along
   - Option C: Hybrid - live code key parts, students complete the rest

2. **Pacing:**
   - Don't rush through Session 1 - it's the foundation
   - Session 2 builds directly on Session 1
   - Session 3 is more advanced - adjust based on time

3. **Debugging:**
   - Encourage students to test after each function
   - Use `fmt.Printf()` for debugging
   - Show how to read Go error messages

### Session-Specific Tips

**Session 1:**
- Start with a quick HTTP/REST refresher
- Show the AISHE API docs at http://localhost:8000/docs
- Test with `curl` first to show expected responses
- Emphasize error handling - it's idiomatic Go

**Session 2:**
- Explain why caching matters (cost, latency, rate limits)
- Show Redis CLI to inspect cache entries:
  ```bash
  docker exec -it aishe-redis redis-cli
  KEYS aishe:answer:*
  GET <key>
  ```
- Discuss cache invalidation strategies
- Measure and compare performance

**Session 3:**
- Start with "what is semantic similarity?"
- Show examples of similar vs different questions
- Explain embeddings at a high level (vectors, cosine similarity)
- Run the benchmark to compare all three approaches
- Discuss when to use each strategy

---

## 📊 Assessment

### Completion Criteria

Students should be able to:

- [ ] Build a working HTTP client in Go
- [ ] Implement JSON encoding/decoding
- [ ] Integrate Redis for caching
- [ ] Explain cache-aside pattern
- [ ] Understand semantic vs exact-match caching
- [ ] Compare trade-offs between caching strategies

### Optional Challenges

For advanced students who finish early:

1. **Add retry logic** - Retry failed requests with exponential backoff
2. **Add metrics** - Track cache hit rate, average latency
3. **Add circuit breaker** - Stop calling API if it's down
4. **Add batch requests** - Ask multiple questions in parallel
5. **Add custom cache key strategy** - Different hashing or normalization

---

## 🐛 Troubleshooting

### Common Issues

**"Connection refused" errors:**
- Check Docker containers are running: `docker-compose ps`
- Verify ports: `lsof -i :8000` (AISHE), `lsof -i :6379` (Redis)
- Restart: `docker-compose down && ./docker-setup.sh`

**"Module not found" errors:**
- Run `go mod tidy` in the session directory
- Check `go.mod` exists
- Verify Go version: `go version` (need 1.21+)

**Redis connection errors:**
- Check Redis is running: `docker exec aishe-redis redis-cli ping`
- Verify `DefaultRedisAddr` constant in `client.go` is set to `"localhost:6379"`
- Check firewall/network settings

**LangCache authentication errors:**
- Verify all three constants are updated in `client.go`: `LangCacheURL`, `LangCacheAPIKey`, `LangCacheCacheID`
- Make sure they're not still set to placeholder values like `"your-api-key-here"`
- Test with curl (replace with actual values):
  ```bash
  curl -H "Authorization: Bearer YOUR_API_KEY" \
    YOUR_LANGCACHE_URL/v1/caches/YOUR_CACHE_ID/entries/search \
    -d '{"prompt":"test"}'
  ```

**Slow responses:**
- First request is always slow (~5s) - this is normal (LLM inference)
- Subsequent cached requests should be <0.1s
- Check Ollama is running: `curl http://localhost:11434/api/tags`

---

## 📝 Post-Workshop

### Follow-up Materials

Share with students:

1. Workshop repository link
2. Additional resources:
   - [Go HTTP Client Best Practices](https://pkg.go.dev/net/http)
   - [Redis Caching Patterns](https://redis.io/docs/manual/patterns/)
   - [LangCache Documentation](https://redis.io/docs/latest/develop/ai/langcache/)
3. Suggested next steps:
   - Build a client for their own API
   - Implement additional caching strategies
   - Add monitoring and metrics

### Feedback

Collect feedback on:
- Pacing (too fast/slow?)
- Difficulty level
- Most/least useful parts
- Suggestions for improvement

---

## 📚 Additional Resources

### For Instructors

- [Go HTTP Package Documentation](https://pkg.go.dev/net/http)
- [go-redis Documentation](https://redis.uptrace.dev/)
- [Redis LangCache API](https://redis.io/docs/latest/develop/ai/langcache/api-examples/)
- [Caching Strategies](https://aws.amazon.com/caching/best-practices/)

### For Students

- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example: HTTP Clients](https://gobyexample.com/http-clients)
- [Redis University](https://university.redis.com/)

---

## 🎉 Success Metrics

A successful workshop means:

- ✅ All students complete at least Session 1 and 2
- ✅ Students understand the trade-offs between caching strategies
- ✅ Students can explain when to use each approach
- ✅ Students have working code they can reference later
- ✅ Students feel confident building HTTP clients in Go

---

**Good luck with the workshop! 🚀**

