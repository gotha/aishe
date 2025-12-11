# Session 1 - Go Solution

This is the complete solution for Session 1: Building a basic CLI client that calls the AISHE API.

## Prerequisites

- Go 1.21 or higher installed
- AISHE server running on `http://localhost:8000`

## Running the Solution

1. Navigate to this directory:
   ```bash
   cd workshop/session-1/go/solution
   ```

2. Run the program with a question:
   ```bash
   go run main.go "What is the capital of France?"
   ```

   Or build and run:
   ```bash
   go build -o aishe-client
   ./aishe-client "What is the capital of France?"
   ```

## Example Output

```
Asking: What is the capital of France?
Waiting for response...

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

----------------------------------------------------------------------
Execution time: 2.47 seconds
----------------------------------------------------------------------
```

## What This Solution Demonstrates

- Making HTTP POST requests to the AISHE API
- Marshaling Go structs to JSON for requests
- Unmarshaling JSON responses to Go structs
- Proper error handling for network and HTTP errors
- Setting timeouts for HTTP clients
- Formatting output for better readability

## Key Components

- **Request struct**: Represents the API request payload
- **Response struct**: Represents the API response with answer, sources, and processing time
- **Source struct**: Represents individual source citations
- **HTTP Client**: Configured with 120-second timeout for long-running queries
- **Error handling**: Graceful handling of connection errors, timeouts, and HTTP errors

## Performance Metrics

- **Processing time**: Time taken by the AISHE API to process the question
- **Execution time**: Total time from receiving the question to displaying the answer (includes network overhead, parsing, etc.)

