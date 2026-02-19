# Session 3: Java - Semantic Caching with LangCache

Command line client for AISHE with semantic caching using Redis Cloud LangCache.

## Prerequisites

1. **Java 17 or higher**
2. **Maven**
3. **AISHE Client Library installed**:
   ```bash
   cd deps/aishe-java
   mvn clean install
   ```
4. **Redis Cloud LangCache credentials** (API_KEY, CACHE_ID, SERVER_URL)

## Setup

1. Create a `.env` file in this directory with your LangCache credentials:
   ```
   SERVER_URL=your-instance.redis.cloud
   CACHE_ID=your-cache-id
   API_KEY=your-api-key
   SIMILARITY_THRESHOLD=0.8
   AISHE_API_URL=http://localhost:8000
   ```

2. Build the project:
   ```bash
   mvn clean compile
   ```

## Usage

```bash
# Run with Maven
mvn exec:java -Dexec.args="What is Redis?"

# Or build standalone JAR and run
mvn clean package
java -jar target/aishe-semantic-client.jar "What is Redis?"
```

## How It Works

Unlike Session 2's exact-match caching, semantic caching matches questions based on **meaning**:

| Session 2 (Exact Match) | Session 3 (Semantic Match) |
|------------------------|---------------------------|
| "What is Python?" ✓ matches "What is Python?" | "What is Python?" ✓ matches "What is Python?" |
| "What is Python?" ✗ does NOT match "Tell me about Python" | "What is Python?" ✓ matches "Tell me about Python" |

## Testing

### Test 1: Exact Match
```bash
mvn exec:java -Dexec.args="What is Redis?"
mvn exec:java -Dexec.args="What is Redis?"
# Second call should hit cache
```

### Test 2: Semantic Match
```bash
mvn exec:java -Dexec.args="What is Redis?"
mvn exec:java -Dexec.args="Tell me about Redis"
# Second call should hit cache (semantically similar)
```

### Test 3: Different Question
```bash
mvn exec:java -Dexec.args="What is Redis?"
mvn exec:java -Dexec.args="What is the capital of France?"
# Second call should miss cache (different topic)
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_URL` | - | LangCache server URL |
| `CACHE_ID` | - | LangCache cache ID |
| `API_KEY` | - | LangCache API key |
| `SIMILARITY_THRESHOLD` | `0.8` | Similarity threshold (0.0-1.0) |
| `AISHE_API_URL` | `http://localhost:8000` | AISHE server URL |

## Example Output

### Cache Miss (First Query)
```
Asking: What is Redis?
✗ Not in cache, calling AISHE API...
Waiting for response...

✓ Response saved to semantic cache

======================================================================
ANSWER:
======================================================================
Redis is an open-source, in-memory data structure store...

======================================================================
Processing time: 8.45 seconds
======================================================================
```

### Cache Hit (Semantic Match)
```
Asking: Tell me about Redis
✓ Found in semantic cache! (no API call needed)
  Similarity score: 0.9234

======================================================================
ANSWER:
======================================================================
Redis is an open-source, in-memory data structure store...

======================================================================
Source: Semantic Cache (LangCache)
Similarity score: 0.9234
Original processing time: 8.45 seconds
======================================================================

======================================================================
Execution time: 0.15 seconds
======================================================================
```

