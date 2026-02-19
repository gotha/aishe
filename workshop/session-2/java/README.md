# Session 2 - Java Implementation with Redis Cache

A command-line Java application that connects to the AISHE web server with Redis caching to store and retrieve previous question-answer pairs, reducing API calls and improving response times.

## Prerequisites

- Java 17 or higher
- Maven 3.6 or higher
- Running AISHE server (default: `http://localhost:8000`)
- Running Redis server on `localhost:6379`

## Setup

### 1. Install the AISHE Java Client Library

First, install the AISHE Java client library to your local Maven repository:

```bash
cd deps/aishe-java
mvn clean install
```

This will make the `aishe-client` library available for use in the workshop sessions.

### 2. Start Required Services

```bash
# From the project root directory
docker-compose up -d aishe redis
```

### 3. Build the Session 2 Application

Navigate to the session-2 java directory:

```bash
cd workshop/session-2/java
```

Compile the application using Maven:

```bash
mvn clean compile
```

(Optional) Build a standalone JAR:

```bash
mvn clean package
```

## Running the Application

### Option 1: Using Maven Exec Plugin

```bash
mvn exec:java -Dexec.args="What is the capital of France?"
```

### Option 2: Using the Standalone JAR

```bash
# After building with 'mvn clean package'
java -jar target/aishe-client.jar "What is the capital of France?"
```

## Environment Variables

You can customize the server URLs using environment variables:

```bash
export AISHE_API_URL=http://localhost:8000
export REDIS_HOST=localhost
export REDIS_PORT=6379

mvn exec:java -Dexec.args="What is the capital of France?"
```

## Example Output

### First Query (Cache Miss)

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

### Second Query (Cache Hit)

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
Source: Redis Cache
======================================================================
```

Notice: No processing time, instant response from cache!

## Implementation Details

The application:
- Uses the AISHE Java Client Library (`aishe-client`) for API communication
- Accepts questions as command-line arguments
- Generates cache keys using the client library's `Utils.generateCacheKey()` method
- Checks Redis cache before making API calls
- Stores responses in Redis with 24-hour expiration (86400 seconds)
- Displays the answer, sources, and timing/cache information
- Handles errors appropriately using the client library's exception classes

### Cache Key Generation

1. Normalize the question (lowercase, trim whitespace)
2. Generate SHA-256 hash of normalized question
3. Prefix with namespace: `aishe:question:{hash}`

This ensures "What is Python?" and "what is python?" produce the same cache key.

## Testing Cache Functionality

### Verify Cache in Redis

```bash
# Connect to Redis
docker exec -it aishe-redis redis-cli

# List all AISHE cache keys
KEYS aishe:question:*

# View a specific cached value
GET aishe:question:{hash}

# Check TTL (time to live)
TTL aishe:question:{hash}

# Exit Redis CLI
exit
```

### Clear Cache

```bash
# Clear all AISHE cache keys
docker exec -it aishe-redis redis-cli --scan --pattern "aishe:question:*" | xargs docker exec -i aishe-redis redis-cli DEL

# Or flush entire Redis database (use with caution)
docker exec -it aishe-redis redis-cli FLUSHDB
```

## Performance Benefits

- **Cache Hit**: ~10-50 milliseconds (instant response)
- **Cache Miss**: ~2-5 seconds (API call required)
- **Speed Improvement**: 50-500x faster for cached responses

## Dependencies

- **aishe-client 1.0.0**: AISHE Java Client Library (includes Gson for JSON parsing)
- **Jedis 5.1.0**: Redis client for Java

## Troubleshooting

### Redis Connection Error

```
Error: Could not connect to Redis at localhost:6379
```

**Solution**: Ensure Redis is running:
```bash
docker-compose up -d redis
docker-compose ps  # Check if redis container is running
```

### AISHE Server Connection Error

```
Error: HTTP error code: 500
```

**Solution**: Ensure AISHE server is running:
```bash
docker-compose up -d aishe
docker-compose ps  # Check if aishe container is running
```

