# Session 1: Basic CLI Client

## Objective

Build a command line application that connects to the AISHE web server, sends a question from the user, and prints the answer.

## Prerequisites

Before starting this session, ensure you have:

1. **AISHE Server Running**: The AISHE server must be running on `http://localhost:8000`
   ```bash
   # From the project root directory
   docker-compose up -d aishe
   ```

## Implementation Overview

### Command-Line Interface

- Accepts questions in the command line

### API Communication


- **Endpoint**: `POST http://localhost:8000/api/v1/ask`
- **Request Format**:
  ```json
  {
    "question": "Your question here"
  }
  ```
- **Response Format**:
  ```json
  {
    "answer": "The generated answer",
    "sources": [
      {
        "number": 1,
        "title": "Wikipedia Article Title",
        "url": "https://en.wikipedia.org/wiki/..."
      }
    ],
    "processing_time": 2.45
  }
  ```

You can test the AISHE API directly using curl:

```bash
curl -X POST http://localhost:8000/api/v1/ask \
  -H "Content-Type: application/json" \
  -d '{"question": "What is the capital of France?"}'
```
