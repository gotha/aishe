# AIshe JavaScript Workshop

Hands-on workshop for learning HTTP clients, caching strategies, and semantic caching with JavaScript.

**NOTE**: the workshop module uses TypeScript. TL;DR: TypeScript = JavaScript + type guards. 
If you don't know or want to use TypeScript, you can simply write your code in 
JavaScript without performing any type checks or guards.

## Prerequisites

- [Node](https://nodejs.org/en) v25.2.1+
- [npm](https://www.npmjs.com/) v11.6.2+
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) (for AIshe server, Redis, Ollama)
- Running AISHE stack on http://localhost:8000 or another port
- Familiarity with JavaScript / TypeScript
- Basic understanding of HTPP and REST APIs
- Knowledge of caching is considered a plus

## Workshop Overview

This workshop teaches you how to build JavaScript/TypeScript clients for AISHE
(AI Search & Help Engine) HTTP API + integrates exact & semantic caching.

AISHE is a cutting-edge AI engine which receives a question from the user. AISHE will
then search Wikipedia and create a summary answer of the information it sources.
Finally, AISHE returns the summarized answer + top 3 sources it used.

In session 1, you'll setup your own local AISHE server. Plus you'll implement a 
basic HTTP client and interactive CLI to interact with AISHE.

In session 2, you'll implement exact caching with basic Redis SET/GET and benchmark
a x1000+ increase in speed. You'll also see the limitations of exact caching.

In session 3, you'll implement semantic caching with LangCache - the newest
proprietary caching technology which stores LLM answers by meaning.

## How to do this workshop

This workshop is organized into three progressive sessions. The directory tree is like so:

```
workshop
├── session1-basic
│   ├── solution
│   │   ├── client.ts
│   │   └── main.ts
│   └── starter
│       ├── client.ts
│       ├── criteria.md
│       └── main.ts
├── session2-redis
│   ├── solution
│   │   ├── client.ts
│   │   └── main.ts
│   └── starter
│       ├── client.ts
│       ├── criteria.md
│       └── main.ts
└── session3-langcache
    ├── solution
    │   ├── client.ts
    │   └── main.ts
    └── starter
        ├── client.ts
        ├── criteria.md
        └── main.ts
```

Inside each session directory you'll find two folders:
- `starter`: contains skeleton code with TODO guidelines for implementation
- `solution`: contains a reference implementation of this session

Each session has three files:
- `criteria.md`: criteria for success with references and help; start each session by reading this file
- `client.ts`: progressively implement AIsheHTTPClient
- `main.ts`: build a simple CLI to interact with AISHE

Your objective is to implement `main.ts` and `client.ts` within the session's time limit.

If you succeed in time, use your own implementation for each next session.

If you don't implement it in time, reference the ready-to-use implementation in `solution` folders.

### Session 1: Basic HTTP Client

**Location:** `workshop/session1-basic/`

Build a simple HTTP client for AISHE.

**Topics:**
- Health checks
- Error handling
- JSON encoding/decoding
- Making HTTP requests & processing HTTP responses

**Run:**
```bash
cd clients/js
npm run session1
```

### Session 2: Redis Cache Client

**Location:** `workshop/session2-redis/`

Add exact-match caching with Redis.

**Topics:**
- Redis integration with node-redis
- Cache key generation (SHA-256)
- Cache HIT/MISS detection
- Basics of caching with Redis GET/SEt

**Run:**
```bash
cd clients/js
npm run session2
```

### Session 3: LangCache Client

**Location:** `workshop/session3-langcache/`

Implement semantic caching with LangCache.

**Topics:**
- Semantic caching
- LangCache REST API
- Similarity thresholds
- Comparing caching strategies

**Run:**
```bash
cd clients/js
export LANGCACHE_API_KEY="..."
export LANGCACHE_CACHE_ID="..."
export LANGCACHE_SERVER_URL="..."
npm run session3
```

## Learning Objectives

By the end of this workshop, you'll understand:

1. ✅ Understand HTTP client patterns in TypeScript
2. ✅ Implement caching strategies for performance
3. ✅ Compare exact-match vs semantic caching
4. ✅ Build interactive CLI applications
5. ✅ Make informed decisions about caching approaches

## Performance Comparison

Typical performance you'll observe:

| Client Type | First Request | Cached (Exact Match) | Cached (Semantic) |
|-------------|---------------|----------------------|-------------------|
| Basic       | ~3-10s        | N/A                  | N/A               |
| Redis       | ~3-10s        | ~0.001s              | N/A               |
| LangCache   | ~3-10s        | ~0.1-0.3s            | ~0.1-0.3s         |

**Key Insights:**

- Exact match caching with Redis is up to x1000+ faster
- LangCache enables x10-x100 faster caching with semantic similarity
- Both caching strategies provide a massive speedup over no caching (+ cut LLM costs)

## Contributing

This workshop is part of the AISHE project. See the main project [README](../../README.md) for contribution guidelines.

## License

Same as the main AISHE project.
