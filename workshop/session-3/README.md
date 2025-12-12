# Session 3: Semantic Caching with LangCache

## Objective

Enhance the CLI client from Session 2 by implementing semantic caching using vector embeddings and similarity search. This allows the cache to return results for questions that are semantically similar, even if the exact wording differs.

## Prerequisites

Before starting this session, ensure you have:

1. **Completed Session 2**: Understanding of basic Redis caching

2. **AISHE Server Running**: The AISHE server must be running on `http://localhost:8000`
   ```bash
   docker-compose up -d aishe
   ```

3. [LangCache access](https://cloud.redis.io/#/langcache): You'll need access to Langcache - the semantic caching service by Redis:
   - Vector embedding generation
   - Similarity search capabilities
   - Redis-backed storage

## Architecture Diagram

```
┌─────────────┐
│    User     │
└──────┬──────┘
       │ Question
       ▼
┌─────────────────────────────────────────────┐
│         CLI Application                     │
│                                             │
│  ┌───────────────────────────────────────┐  │
│  │  1. Search in LangCache               │  │
│  └───────────────────────────────────────┘  │
│         │                │                  │
│    Cache Hit        Cache Miss              │
│         │                │                  │
│         │                ▼                  │
│         │  ┌───────────────────────────┐    │
│         │  │ 2. Call AISHE API         │    │
│         │  └───────────────────────────┘    │
│         │                │                  │
│         │                ▼                  │
│         │  ┌───────────────────────────┐    │
│         │  │ 3. Save to LangCache      │    │
│         │  │    (Embedding + Response) │    │
│         │  └───────────────────────────┘    │
│         │                │                  │
│         └────────┬───────┘                  │
│                  ▼                          │
│  ┌───────────────────────────────────────┐  │
│  │  4. Display Response                  │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

## Implementation Approaches

This session can be implemented in multiple languages and frameworks. Choose the approach that best fits your stack:

### Officially supported Libraries 

- javascript / typescript - https://www.npmjs.com/package/@redis-ai/langcache?activeTab=readme
- python - https://pypi.org/project/langcache/

## Implementation Overview

### Semantic Caching Concept

Unlike Session 2's exact-match caching, semantic caching matches questions based on **meaning** rather than exact text:

**Session 2 (Exact Match)**:
- "What is Python?" ✓ matches "What is Python?"
- "What is Python?" ✗ does NOT match "Tell me about Python"
- "What is Python?" ✗ does NOT match "Explain Python programming language"

**Session 3 (Semantic Match)**:
- "What is Python?" ✓ matches "What is Python?"
- "What is Python?" ✓ matches "Tell me about Python" (similar meaning)
- "What is Python?" ✓ matches "Explain Python programming language" (similar meaning)

### How Semantic Caching Works

```
1. User asks a question
2. Convert question to vector embedding (e.g., 1536-dimensional vector)
3. Search cache for similar embeddings using similarity threshold
4. If similar question found (similarity > threshold):
   - Return cached response (cache hit)
5. If no similar question found:
   - Call AISHE API
   - Store question embedding + response in cache
   - Return API response
```


## Testing Your Implementation

### Test 1: Exact Match (Should Hit Cache)

```bash
# First query
./your-cli "What is Redis?"

# Exact same query
./your-cli "What is Redis?"
```

Expected: Cache hit on second query

### Test 2: Semantic Match (Should Hit Cache)

```bash
# First query
./your-cli "What is Redis?"

# Semantically similar queries
./your-cli "Tell me about Redis"
./your-cli "Explain Redis database"
./your-cli "What does Redis do?"
```

Expected: All should hit cache (if similarity > threshold)

### Test 3: Different Questions (Should Miss Cache)

```bash
# First query
./your-cli "What is Redis?"

# Unrelated query
./your-cli "What is the capital of France?"
```

Expected: Cache miss on second query

### Test 4: Paraphrased Questions

```bash
# Original
./your-cli "How does machine learning work?"

# Paraphrases
./your-cli "Explain the concept of machine learning"
./your-cli "What is the mechanism behind ML?"
./your-cli "Can you describe how ML functions?"
```

Expected: Cache hits for paraphrases (depending on threshold)

## Performance Comparison

### Session 1: No Cache
- Every question: API call required
- Response time: 2-5 seconds

### Session 2: Exact Match Cache
- Exact same question: Cache hit (~50ms)
- Slightly different wording: Cache miss (API call)
- Cache hit rate: ~20-30% (typical usage)

### Session 3: Semantic Cache
- Exact same question: Cache hit (~50-100ms)
- Similar question: Cache hit (~50-100ms)
- Different question: Cache miss (API call)
- Cache hit rate: **60-80%** (typical usage)
- **2-4x better cache hit rate** than exact matching

### Trade-offs

**Advantages**:
- Much higher cache hit rate
- Better user experience (more questions answered from cache)
- Reduced API costs and server load

**Disadvantages**:
- Slightly slower cache lookup (embedding + similarity search)
- Requires embedding service/model
- More complex implementation
- Potential for false positives (wrong answer for similar question)

## Common Issues

### Low Cache Hit Rate

**Symptoms**: Most queries result in cache misses despite similar questions

**Possible causes**:
1. Similarity threshold too high (too strict)
2. Embedding model not capturing semantic meaning well
3. Questions too diverse

**Solutions**:
```bash
# Lower the similarity threshold
threshold = 0.80  # from 0.90

# Test with known similar questions
./your-cli "What is Python?"
./your-cli "Tell me about Python"
# Should both hit cache
```

### False Positives (Wrong Cached Answers)

**Symptoms**: Cache returns answer for Question A when asked Question B

**Possible causes**:
1. Similarity threshold too low (too permissive)
2. Questions appear similar but have different intent

**Example**:
```
Q1: "What is Python?" → Answer about programming language
Q2: "What is a python?" → Should answer about the snake, but gets cached programming answer
```

**Solutions**:
```bash
# Increase similarity threshold
threshold = 0.90  # from 0.80

# Add intent detection
# Check if question is about programming vs. animals
```

### Slow Cache Lookups

**Symptoms**: Cache hits take longer than expected

**Possible causes**:
1. Large number of cached entries
2. Inefficient similarity search
3. Network latency to embedding service

**Solutions**:
- Use vector indexing (HNSW, IVF)
- Limit cache size
- Use local embedding models
- Implement cache partitioning
