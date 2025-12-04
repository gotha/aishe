# AISHE Quick Start Guide

Get up and running with AISHE in 5 minutes!

## What is AISHE?

AISHE is a Wikipedia-based RAG (Retrieval-Augmented Generation) question answering system that uses local LLMs via Ollama to provide accurate answers with source citations.

## Prerequisites

- Docker (20.10+)
- Docker Compose (1.29+)
- At least 8GB free disk space

## Step 1: Start Docker Services

Run the automated setup script:

```bash
./docker-setup.sh
```

This will:
- ✅ Start Redis (caching)
- ✅ Start Ollama (LLM inference)
- ✅ Start AISHE API server
- ✅ Pull the `llama3.2:3b` model
- ✅ Verify all services are healthy

**Expected output:**
```
✅ All services started successfully!

Service URLs:
  AISHE API:  http://localhost:8000
  API Docs:   http://localhost:8000/docs
  Ollama:     http://localhost:11434
  Redis:      localhost:6379

Quick test:
  curl http://localhost:8000/health
```

## Step 2: Test the API

### Health Check

```bash
curl http://localhost:8000/health
```

Expected response:
```json
{
  "status": "healthy",
  "ollama_accessible": true
}
```

### Ask a Question

```bash
curl -X POST http://localhost:8000/api/v1/ask \
  -H "Content-Type: application/json" \
  -d '{"question": "What is the capital of France?"}'
```

## Step 3: Try the Go Workshop

The workshop teaches you to build three different Go clients with progressive complexity.

### Navigate to Workshop

```bash
cd clients/go/workshop
```

### Session 1: Basic HTTP Client

Learn HTTP client basics without caching:

```bash
cd session1-basic/starter
go run .
```

**What you'll learn:**
- HTTP client fundamentals
- JSON marshaling/unmarshaling
- Error handling

### Session 2: Redis Cache Client

Add exact-match caching with Redis:

```bash
cd ../session2-redis/starter
go run .
```

**What you'll learn:**
- Cache-aside pattern
- Redis integration
- Performance measurement
- ~18,000x speedup on cache hits!

### Session 3: LangCache Client (Advanced)

Implement semantic caching:

```bash
cd ../session3-langcache/starter
# Update LangCache constants in client.go first
go run .
```

**What you'll learn:**
- Semantic caching
- REST API integration
- Benchmarking all three approaches

### Run the Benchmark

Compare all three approaches side-by-side:

```bash
cd session3-langcache/solution
go run . benchmark
```

This shows:
- Basic client: No caching
- Redis client: Exact-match caching
- LangCache: Semantic caching (requires credentials)

## Workshop Documentation

- **`README.md`** - Student guide with detailed instructions
- **`INSTRUCTOR_GUIDE.md`** - Teaching guide with timing and tips
- **`WORKSHOP_STRUCTURE.md`** - Technical documentation

## Managing Docker Services

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f aishe
docker-compose logs -f ollama
docker-compose logs -f redis
```

### Stop Services

```bash
# Stop all (keeps data)
docker-compose stop

# Stop and remove containers (keeps data)
docker-compose down

# Remove everything including data
docker-compose down -v
```

### Restart Services

```bash
docker-compose restart
```

## Troubleshooting

### Docker Not Running

```bash
# macOS
open -a Docker

# Wait for Docker to start, then run setup again
./docker-setup.sh
```

### Services Not Healthy

```bash
# Check service status
docker-compose ps

# Check logs
docker-compose logs aishe
docker-compose logs ollama

# Restart services
docker-compose restart
```

### Port Already in Use

Edit `docker-compose.yml` to change ports:

```yaml
services:
  aishe:
    ports:
      - "8001:8000"  # Use port 8001 instead of 8000
```

### Ollama Model Not Pulled

```bash
# Pull the model manually
docker exec aishe-ollama ollama pull llama3.2:3b

# Verify
docker exec aishe-ollama ollama list
```

## Next Steps

1. **Explore the API**: Visit http://localhost:8000/docs for interactive API documentation
2. **Complete the Workshop**: Work through all three sessions in `clients/go/workshop/`
3. **Read the Docs**:
   - [README.md](README.md) - Full project documentation
   - [DOCKER.md](DOCKER.md) - Detailed Docker setup guide
   - [clients/go/README.md](clients/go/README.md) - Go clients documentation

## Quick Reference

### Service URLs

- **AISHE API**: http://localhost:8000
- **API Docs**: http://localhost:8000/docs
- **Ollama**: http://localhost:11434
- **Redis**: localhost:6379

### Useful Commands

```bash
# Test API
curl http://localhost:8000/health

# Ask a question
curl -X POST http://localhost:8000/api/v1/ask \
  -H "Content-Type: application/json" \
  -d '{"question": "Your question here"}'

# Check Ollama models
docker exec aishe-ollama ollama list

# Test Redis
docker exec aishe-redis redis-cli ping

# View all logs
docker-compose logs -f

# Restart everything
docker-compose restart
```

## Workshop Quick Start

```bash
# Session 1: Basic client
cd clients/go/workshop/session1-basic/starter
go run .

# Session 2: Redis cache
cd ../session2-redis/starter
go run .

# Session 3: LangCache (update constants first)
cd ../session3-langcache/starter
# Edit client.go to add LangCache credentials
go run .

# Run benchmark
cd ../solution
go run . benchmark
```

---

**You're all set! 🚀**

For detailed information, see:
- [README.md](README.md) - Full documentation
- [DOCKER.md](DOCKER.md) - Docker details
- [clients/go/workshop/README.md](clients/go/workshop/README.md) - Workshop guide

