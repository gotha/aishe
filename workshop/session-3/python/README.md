# Session 3: Semantic Caching with LangCache - Python Setup

This directory contains the Python implementation for Session 3 of the AISHE workshop, which implements semantic caching using Redis LangCache.

## Prerequisites

Before starting, ensure you have:

1. **Completed Sessions 1 & 2**: Understanding of basic CLI and simple caching
2. **AISHE Server Running**: The AISHE server must be accessible
3. **Redis Stack Running**: Redis with vector search capabilities
   - Use Redis Stack (not regular Redis)
   - Supports vector similarity search
4. **Python 3.12+**: Required for this project

## Project Setup

### 1. Sync Dependencies

```bash
# Install all dependencies including LangCache
uv sync
```

This will install:
- `langcache>=0.11.1` - Semantic caching library for Redis
- `python-dotenv>=1.2.1` - Environment variable management
- `requests>=2.32.5` - HTTP client for API calls

### 2. Activate Virtual Environment

**On macOS/Linux:**
```bash
source .venv/bin/activate
```

**On Windows:**
```powershell
.venv\Scripts\activate
```


### 4. Configure Environment

Create a `.env` file in this directory:

```bash
# .env file
AISHE_SERVER_URL=http://localhost:8000
REDIS_URL=redis://localhost:6379
LANGCACHE_SIMILARITY_THRESHOLD=0.85
```

## Running the CLI with Semantic Caching

```bash
python main.py "What is Redis?"
```

### Inspect LangCache in Redis

```bash
# Connect to Redis CLI
redis-cli

# List LangCache indexes
FT._LIST

# Search for cached embeddings
FT.SEARCH langcache:index "*"

# Clear cache
FLUSHDB
```

## Key Concepts

### Semantic Caching Strategy

- **Embeddings**: Convert questions to vector representations
- **Similarity Search**: Find semantically similar questions
- **Threshold**: Configurable similarity score (0.0-1.0)
- **Hit**: Similar question found → return cached answer
- **Miss**: No similar question → call API, store with embedding

### How It Works

1. **Question arrives** → "What is Redis?"
2. **Generate embedding** → [0.123, 0.456, 0.789, ...]
3. **Search cache** → Find similar embeddings
4. **Check threshold** → Similarity > 0.85?
5. **Cache hit** → Return cached answer
6. **Cache miss** → Call API, cache with embedding

## Troubleshooting

### Redis Stack Not Running

```
Error: Redis module 'search' not loaded
```

**Solution:** Use Redis Stack, not regular Redis:
```bash
docker run -d -p 6379:6379 redis/redis-stack:latest
```

### LangCache Import Error

```
ModuleNotFoundError: No module named 'langcache'
```

**Solution:**
```bash
uv sync
source .venv/bin/activate
```

### Low Cache Hit Rate

**Adjust similarity threshold:**
```bash
# More lenient (more hits, less accurate)
export LANGCACHE_SIMILARITY_THRESHOLD=0.75

# More strict (fewer hits, more accurate)
export LANGCACHE_SIMILARITY_THRESHOLD=0.95
```

### Embedding Model Issues

**Ensure Ollama has the embedding model:**
```bash
ollama pull all-minilm:latest
ollama list
```

## Performance Comparison

### Session 1 (No Cache)
- Every request calls API
- Response time: ~2-5 seconds
- Cost: High (every request)

### Session 2 (Simple Cache)
- Only exact matches cached
- Cache hit rate: ~20-30%
- Response time: <10ms (cached)

### Session 3 (Semantic Cache)
- Semantic matches cached
- Cache hit rate: ~70-90%
- Response time: <50ms (cached)
- Cost savings: 70-90%
