# AIshe - Wikipedia RAG Question Answering System

AIshe is an AI assistant that can help you get the answers you need using Wikipedia as a knowledge base and local LLM for answer generation.

## 🚀 Quick Start

**New to AISHE?** Start here: **[QUICKSTART.md](QUICKSTART.md)**

The quickstart guide will get you up and running in 5 minutes with:
- Docker setup
- API testing

## Overview

This project implements a Retrieval-Augmented Generation (RAG) system that:
- Searches Wikipedia for relevant articles based on your questions
- Retrieves and processes article content
- Generates accurate answers using a local LLM (via Ollama)
- Provides source citations for transparency

## Prerequisites

### Option 1: Docker (Recommended)

- Docker (20.10+)
- Docker Compose (1.29+)

See [QUICKSTART.md](QUICKSTART.md) or [DOCKER.md](DOCKER.md) for setup instructions.

### Option 2: Nix Development Environment

- [Nix](https://nixos.org/download/#download-nix)
- Ollama (for local LLM)

## Setup

### Docker Setup (Recommended)

Run the automated setup script:

```bash
./docker-setup.sh
```

This will start all services (AISHE, Ollama, Redis) in Docker containers. See [QUICKSTART.md](QUICKSTART.md) for details.

### Local Development Setup

#### 1. Setup Development Environment

```sh
nix develop
# or
# direnv allow .
```

#### 2. Start Ollama Server

Make sure Ollama is running in a separate terminal:

```bash
ollama serve
ollama pull llama3.2:3b
ollama ls
```

## Usage

### CLI Tool

```bash
python src/cli.py
```

You'll see a prompt where you can ask questions:

```
======================================================================
Wikipedia RAG Question Answering System
======================================================================
Ask questions and get answers based on Wikipedia articles.
Type 'quit' or 'exit' to stop.
======================================================================

Your question: What is Python programming language?
```

Example output:

```
──────────────────────────────────────────────────────────────────────
ANSWER:
──────────────────────────────────────────────────────────────────────
Python is a high-level, interpreted programming language created by
Guido van Rossum and first released in 1991. It emphasizes code
readability with its use of significant indentation and supports
multiple programming paradigms including procedural, object-oriented,
and functional programming.

──────────────────────────────────────────────────────────────────────
SOURCES:
──────────────────────────────────────────────────────────────────────
[1] Python (programming language)
    https://en.wikipedia.org/wiki/Python_(programming_language)
[2] History of Python
    https://en.wikipedia.org/wiki/History_of_Python
──────────────────────────────────────────────────────────────────────
```

## API Server

AISHE also provides a FastAPI-based HTTP server:

### Start the Server (Docker)

```bash
./docker-setup.sh
```

The API will be available at http://localhost:8000

### Start the Server (Local)

```bash
python -m uvicorn src.server:app --host 0.0.0.0 --port 8000
```

### API Documentation

- Interactive docs: http://localhost:8000/docs
- Health check: http://localhost:8000/health
- Ask endpoint: `POST /api/v1/ask`

### Example API Usage

```bash
curl -X POST http://localhost:8000/api/v1/ask \
  -H "Content-Type: application/json" \
  -d '{"question": "What is the capital of France?"}'
```

## License

BSD

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
Fair warning: this is mostly vibe coded
