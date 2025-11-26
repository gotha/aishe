# AIshe - Wikipedia RAG Question Answering System

AIshe is an AI assistant that can help you get the answers you need using Wikipedia as a knowledge base and local LLM for answer generation.

## Overview

This project implements a Retrieval-Augmented Generation (RAG) system that:
- Searches Wikipedia for relevant articles based on your questions
- Retrieves and processes article content
- Generates accurate answers using a local LLM (via Ollama)
- Provides source citations for transparency

## Prerequisites

[nix](https://nixos.org/download/#download-nix)

## Quick Start

### Setup

```sh
nix develop
# or 
# direnv allow .
```

### Start Ollama Server

Make sure Ollama is running in a separate terminal:

```bash
ollama serve
```

### Run the CLI Tool


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

## License

BSD

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
Fair warning: this is mostly vibe coded
