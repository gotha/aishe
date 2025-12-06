# AISHE Go Workshop

Hands-on workshop for learning HTTP clients, caching strategies, and semantic caching with Go.

## 📚 Documentation

**Main Workshop Documentation:**
- **[Workshop Guide](../../docs/workshop/README.md)** - Complete workshop guide (language-agnostic)
- **[Instructor Guide](../../docs/workshop/INSTRUCTOR_GUIDE.md)** - Teaching notes and solutions
- **[Workshop Structure](../../docs/workshop/WORKSHOP_STRUCTURE.md)** - Technical details and architecture

**Quick Start:**
- See [QUICKSTART.md](../../QUICKSTART.md) in the project root for setup instructions

## 🎓 Workshop Overview

This workshop teaches you how to build Go clients for the AISHE (AI Search & Help Engine) HTTP API, progressing from basic HTTP requests to advanced semantic caching.

**What you'll learn:**
- Building HTTP clients in Go
- Implementing Redis caching for performance
- Using semantic caching with LangCache
- Comparing different caching strategies
- Interactive CLI development

## 📚 Workshop Structure

The workshop is organized into three progressive sessions:

### Session 1: Basic HTTP Client
**Location:** `workshop/session1-basic/`

Learn the fundamentals of HTTP clients in Go.

**Topics:**
- Making HTTP requests with `net/http`
- JSON encoding/decoding
- Error handling
- Health checks

**Files:**
- `starter/` - Starting point with TODOs
- `solution/` - Complete reference implementation

### Session 2: Redis Cache Client
**Location:** `workshop/session2-redis/`

Add exact-match caching with Redis.

**Topics:**
- Redis integration with go-redis
- Cache key generation (SHA-256)
- TTL configuration
- Cache hit/miss detection

**Files:**
- `starter/` - Starting point with TODOs
- `solution/` - Complete reference implementation

### Session 3: LangCache Client
**Location:** `workshop/session3-langcache/`

Implement semantic caching with LangCache.

**Topics:**
- Semantic similarity
- LangCache REST API
- Similarity thresholds
- Comparing caching strategies

**Files:**
- `starter/` - Starting point with TODOs
- `solution/` - Complete reference implementation

## 🏃 Running the Workshop

Each session has an interactive loop where you can:
1. Ask questions
2. See the answer
3. View the source (API, Redis Cache, or LangCache)
4. Check execution time
5. Type "exit" to quit

### Session 1: Basic Client

```bash
cd workshop/session1-basic/starter
go run .
```

### Session 2: Redis Cache

```bash
cd workshop/session2-redis/starter
go run .
```

### Session 3: LangCache

```bash
cd workshop/session3-langcache/starter
# First, update LangCache constants in client.go
go run .
```

## 🎯 Learning Objectives

By the end of this workshop, you will:

1. ✅ Understand HTTP client patterns in Go
2. ✅ Implement caching strategies for performance
3. ✅ Compare exact-match vs semantic caching
4. ✅ Build interactive CLI applications
5. ✅ Make informed decisions about caching approaches

## 🛠️ Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for AISHE server, Redis, Ollama)
- Basic understanding of HTTP and REST APIs
- Familiarity with Go syntax

## 📊 Performance Comparison

Typical performance you'll observe:

| Client Type | First Request | Cached (Exact Match) | Cached (Semantic) |
|-------------|---------------|---------------------|-------------------|
| Basic       | ~5-10s        | N/A                 | N/A               |
| Redis       | ~5-10s        | ~0.001s ✅          | N/A               |
| LangCache   | ~5-10s        | ~0.1-0.5s ✅        | ~0.1-0.5s ✅      |

**Key Insights:**
- Redis is fastest for exact matches
- LangCache enables semantic matching (similar questions hit cache)
- Both provide massive speedup over no caching

## 🤝 Contributing

This workshop is part of the AISHE project. See the main project README for contribution guidelines.

## 📝 License

Same as the main AISHE project.

