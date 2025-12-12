# Session 2: CLI with Redis Caching - Python Setup

## Prerequisites

Before starting, ensure you have:

1. **Completed Session 1**: Understanding of the basic CLI implementation
2. **AISHE Server Running**: The AISHE server must be accessible
3. **Redis Server Running**: A Redis instance must be available
   - Local: `redis://localhost:6379`
   - Or use a remote Redis instance in the cloud 
4. **Python 3.12+**: Required for this project

## Project Setup

### Sync Dependencies

```bash
# Install all dependencies including Redis client
uv sync
```

This will install:
- `redis>=7.1.0` - Redis client for Python
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

### 3. Start Redis (if not already running)

**Using Docker:**
```bash
docker run -d -p 6379:6379 redis:latest
```

**Verify Redis is running:**
```bash
redis-cli ping
# Should return: PONG
```

## Running the CLI with Caching


```bash
export AISHE_SERVER_URL=http://localhost:8000
export REDIS_URL=redis://localhost:6379

# Run the CLI
python main.py "What is Redis?"

# Run the same question again - should be instant (cached)
python main.py "What is Redis?"
```


### Inspect Redis Cache

```bash
# Connect to Redis CLI
redis-cli

# List all keys
KEYS *

# Get a cached value (replace with actual key)
GET "aishe:question:What is the capital of France?"

# Check TTL (time to live)
TTL "aishe:question:What is the capital of France?"

# Clear all cache
FLUSHDB
```

## Key Concepts

### Simple Caching Strategy

- **Cache Key**: Question text (exact match)
- **Cache Value**: Complete API response
- **TTL**: Time-to-live for cache entries (e.g., 1 hour)
- **Hit**: Question found in cache → return immediately
- **Miss**: Question not in cache → call API, store result

## Troubleshooting

### Redis Connection Error

```
redis.exceptions.ConnectionError: Error connecting to Redis
```

**Solution:**
1. Ensure Redis is running: `redis-cli ping`
2. Check Redis URL: `echo $REDIS_URL`
3. Verify port is not blocked: `telnet localhost 6379`

### Cache Not Working

**Check if Redis is storing data:**
```bash
redis-cli
> KEYS *
> GET "your-cache-key"
```

**Clear cache and try again:**
```bash
redis-cli FLUSHDB
```
