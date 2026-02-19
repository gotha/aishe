# Session 1 - Java Implementation

A simple command-line Java application that connects to the AISHE web server, sends a question from the user, and prints the answer.

## Prerequisites

- Java 17 or higher
- Maven 3.6 or higher
- Running AISHE server (default: `http://localhost:8000`)

## Setup

### 1. Install the AISHE Java Client Library

First, install the AISHE Java client library to your local Maven repository:

```bash
cd deps/aishe-java
mvn clean install
```

This will make the `aishe-client` library available for use in the workshop sessions.

### 2. Build the Session 1 Application

Navigate to the session-1 java directory:

```bash
cd workshop/session-1/java
```

Compile the application using Maven:

```bash
mvn clean compile
```

## Running the Application

### Option 1: Using Maven Exec Plugin

```bash
mvn exec:java -Dexec.args="What is the capital of France?"
```

### Option 2: Compile and Run Directly

```bash
# Compile
mvn clean compile

# Run with classpath
mvn exec:java -Dexec.mainClass="Main" -Dexec.args="What is the capital of France?"
```

### Option 3: Build JAR and Run

```bash
# Package the application
mvn clean package

# Run the JAR (requires dependencies in classpath)
java -cp target/session-1-1.0.0.jar:~/.m2/repository/com/aishe/aishe-client/1.0.0/aishe-client-1.0.0.jar Main "What is the capital of France?"
```

## Environment Variables

You can customize the AISHE server URL using the `AISHE_URL` environment variable:

```bash
export AISHE_URL=http://localhost:8000
mvn exec:java -Dexec.args="What is the capital of France?"
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
Total execution time: 2.50 seconds
======================================================================
```

## Implementation Details

The application:
- Uses the AISHE Java Client Library (`aishe-client`) for API communication
- Accepts questions as command-line arguments
- Displays the answer, sources, and timing information
- Handles errors appropriately using the client library's exception classes

## API Endpoint

- **URL**: `POST http://localhost:8000/api/v1/ask`
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

