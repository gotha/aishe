# Session 2 - Go Solution

This is the complete solution for Session 2: Adding Redis caching to improve performance and reduce API calls.

## Prerequisites

- Go 1.21 or higher installed
- AISHE server running on `http://localhost:8000`
- Redis server running on `localhost:6379`

## Running the Solution

1. Navigate to this directory:
   ```bash
   cd workshop/session-2/go/solution
   ```

2. Install dependencies (if not already done):
   ```bash
   go mod download
   ```

3. Run the program with a question:
   ```bash
   go run main.go "What is the capital of France?"
   ```

   Or build and run:
   ```bash
   go build -o aishe-client
   ./aishe-client "What is the capital of France?"
   ```

## Example Output

### First Run (Cache Miss)
```
Asking: What is the capital of France?
✗ Not in cache, calling AISHE API...
Waiting for response...

✓ Response saved to cache

======================================================================
ANSWER:
======================================================================
The capital of France is Paris.

======================================================================
SOURCES:
======================================================================
[1] Paris - Wikipedia
    https://en.wikipedia.org/wiki/Paris

======================================================================
Processing time: 2.45 seconds
======================================================================
```

### Second Run (Cache Hit)
```
Asking: What is the capital of France?
✓ Found in cache! (no API call needed)

======================================================================
ANSWER:
======================================================================
The capital of France is Paris.

======================================================================
SOURCES:
======================================================================
[1] Paris - Wikipedia
    https://en.wikipedia.org/wiki/Paris

======================================================================
Source: Cache
======================================================================
```

## What This Solution Demonstrates

- Connecting to Redis using `github.com/redis/go-redis/v9`
- Generating consistent cache keys using SHA-256 hashing
- Normalizing questions (lowercase, trimmed) for better cache hits
- Storing and retrieving JSON data in Redis
- Setting cache expiration (24 hours)
- Handling cache misses gracefully
- Using Go contexts for Redis operations

## Key Components

- **getCacheKey()**: Normalizes questions and generates SHA-256 hash-based cache keys
- **getFromCache()**: Retrieves cached responses from Redis
- **saveToCache()**: Stores responses in Redis with 24-hour expiration
- **Redis Client**: Configured to connect to localhost:6379
- **Cache namespace**: Uses `aishe:question:{hash}` format for keys

## Cache Behavior

- Cache keys are generated from normalized questions (lowercase, trimmed)
- Questions like "What is Python?" and "what is python?" will hit the same cache entry
- Cached responses expire after 24 hours
- Cache misses trigger API calls, and responses are automatically cached

