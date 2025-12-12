# Session 3 - Go Solution

This is the complete solution for Session 3: Implementing semantic caching with LangCache for intelligent question matching.

## Prerequisites

- Go 1.21 or higher installed
- AISHE server running (default: `http://localhost:8000`)
- Redis Cloud LangCache account with credentials

## Setup

1. Copy the `.env.example` file to `.env` and fill in your credentials:
   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file with your actual credentials:
   ```bash
   AISHE_URL=http://localhost:8000
   SERVER_URL=your-instance.redis.cloud
   CACHE_ID=your-cache-id
   API_KEY=your-api-key
   SIMILARITY_THRESHOLD=0.8
   ```

   **Configuration Notes**:
   - `AISHE_URL`: The AISHE server URL (default: `http://localhost:8000`)
   - `SERVER_URL`: LangCache hostname (without `https://`). The code will automatically add the `https://` prefix.
   - `CACHE_ID`: Your LangCache cache ID
   - `API_KEY`: Your LangCache API key
   - `SIMILARITY_THRESHOLD`: Controls how similar questions need to be for a cache hit (0.0 to 1.0). Higher values require closer matches. Default is 0.8.

## Running the Solution

1. Navigate to this directory:
   ```bash
   cd workshop/session-3/go/solution
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

✓ Response saved to semantic cache

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

======================================================================
Execution time: 2.48 seconds
======================================================================
```

### Second Run with Similar Question (Semantic Cache Hit)
```
Asking: What's the capital city of France?
✓ Found in semantic cache! (no API call needed)
  Similarity score: 0.9234

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
Source: Semantic Cache (LangCache)
Similarity score: 0.9234
Original processing time: 2.48 seconds
======================================================================

======================================================================
Execution time: 0.12 seconds
======================================================================
```

## What This Solution Demonstrates

- Using LangCache API for semantic similarity search
- Making authenticated HTTP requests with Bearer tokens
- Loading environment variables from `.env` files using `github.com/joho/godotenv`
- Semantic matching with configurable similarity threshold (0.8)
- Storing and retrieving responses based on meaning, not exact text
- Displaying similarity scores when cache hits occur (if provided by the API)
- Handling LangCache API responses with proper error checking

## Key Components

- **LangCacheClient**: Custom client struct for managing LangCache API interactions
- **getFromCache()**: Searches for semantically similar questions using LangCache search API
- **saveToCache()**: Stores question-response pairs in LangCache
- **Similarity Threshold**: Set to 0.8 to allow semantic matches while avoiding false positives
- **Environment Variables**: Secure credential management using `.env` file

## Semantic Caching Behavior

Unlike traditional caching (Session 2), semantic caching matches questions by **meaning**:

- "What is the capital of France?" ✓
- "What's the capital city of France?" ✓ (semantic match)
- "Capital of France?" ✓ (semantic match)

All these variations will hit the same cache entry because they have similar semantic meaning!

## LangCache API Endpoints Used

1. **Search**: `POST /cache/{cache_id}/search`
   - Finds semantically similar cached entries
   - Uses similarity threshold to control matching strictness

2. **Set**: `POST /cache/{cache_id}/set`
   - Stores new question-response pairs
   - Automatically generates embeddings for semantic search

## Performance Metrics

- **Processing time**: Time taken by the AISHE API to process the question (shown on cache miss)
- **Original processing time**: The original API processing time from when the response was first cached (shown on cache hit)
- **Execution time**: Total time from receiving the question to displaying the answer (always shown)
- **Similarity score**: When a cache hit occurs, shows how similar the cached question is to your query (0.0 to 1.0, where 1.0 is identical). Only displayed if the LangCache API returns this value.
- Semantic cache hits are faster (~0.12s) compared to API calls (~2.48s), though slightly slower than traditional Redis cache due to similarity search overhead

## Troubleshooting

If you see credential errors:
- Verify your `.env` file exists in this directory
- Check that `SERVER_URL`, `CACHE_ID`, and `API_KEY` are set correctly
- Ensure `SERVER_URL` is just the hostname (e.g., `your-instance.redis.cloud`)
- Verify your LangCache API key is valid and has proper permissions

