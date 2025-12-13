---
author: Redis Team Sofia
title: Redis AI Workshop
date: 13 December 2025
---

## What are we going to do today?

- Session 1: Create a command line app that interacts with RAG system
- Session 2: Optimize for speed (see how fast it feels)
- Session 3: Implement semantic caching 

--- 

## Rules

- Have fun! 
- Make friends
- Build software

<aside class="notes">
</aside>

---

## grep '^staff' /etc/groups

Say hi to the team!

---

## Requirements 

- modern operating system (preferrably unix-like but WSL also works)
- docker 
- dev tools (git, make, curl, etc)
- dev env (nodejs, python, golang, etc) + editor

---

## Structure

For each session:

- 5 mins intro 
- 45 mins coding session
- 5 mins wrap up

we encourage exploration, asking questions, pair programming.

---

## Session 1: Create a command line app that interacts with RAG system

---

## Say hi to AISHE

- github.com/gotha/aishe
- architecture overview
- dockerfile
- docker-compose

---

### AISHE architecture 

```
┌────────────────────────────────────────────────────────────────────┐
│                         AISHE                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                                                             │   │
│  │  • POST /api/v1/ask  - Answer questions                     │   │
│  └────────────────────────────┬────────────────────────────────┘   │
│                               ▼                                    │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                   RAG Pipeline (rag_pipeline.py)            │   │
│  │            1. Retrieve Wikipedia Articles                   │   │
│  │            2. Context Preparation                           │   │ 
│  │            3. Generation Phase                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│         │                                             │            │
│         ▼                                             ▼            │
│  ┌──────────────────────┐                      ┌────────────┐      │
│  │  Wikipedia (DB)      │                      │   Ollama   │      │
│  └──────────────────────┘                      └────────────┘      │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```
---

## Session 1: Goal

- write a CLI that connects to AISHE and asks questions
- ./workshop/session-1/README.md
- lets see some code

---

### Stretch goal / homework - AISHE Modern architecture 

```
┌────────────────────────────────────────────────────────────────────┐
│                         AISHE                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                                                             │   │
│  │  • POST /api/v1/ask  - Answer questions                     │   │
│  └────────────────────────────┬────────────────────────────────┘   │
│                               ▼                                    │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                   RAG Pipeline (rag_pipeline.py)            │   │
│  │            1. send the prompt to Ollama with tools          │   │
│  │            2. parse ollama response and call tools          │   │ 
│  │            3. Generation loop                               │   │
│  └─────────────────────────────────────────────────────────────┘   │
│         │ (mcp)                                       │            │
│         ▼                                             ▼            │
│  ┌────────────────────────────┐                ┌────────────┐      │
│  │  Wikipedia / Other sources │                │   Ollama   │      │
│  └────────────────────────────┘                └────────────┘      │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```


## Session 1: Wrap up

- did we manage to get an answer?
- what we have learned ?

---

## Problem 

- LLMs are slow 
- LLMs are "expensive"

---

## Lets add some caching - Redis

- In memory store 
- Scalable 
- Shared state

---

## Session 2: Goal

- lets add some cache
- ./workshop/session-2/README.md
- get an account in Redis Cloud - redis.io
- show me the code

---

## Cache hit vs semantic cache hit

- classic cache hit or miss vs semantic cache hit or miss
    - vector search - Redis Search vs Redis VectorSets
- MeanCache paper - https://arxiv.org/abs/2403.02694
    - langcache model
- trade-offs - higher cache hit rate vs slower cache lookup

---

### Stretch goal / homework - implement semantic cache yourself

- Redis as vector store and  
    - `all-minilm:latest` via ollama
    - or langcache model - https://huggingface.co/redis/langcache-embed-v2
- vectorize prompt
- search for similar prompts
    - determine if we have a semantic cache hit
- store vectorized prompt and response

---

## Session 2: Wrap up

- how much faster it is?
- what we have learned ?

---

## Session 3: Goal

- introducing LangCache
- ./workshop/session-3/README.md
- show me the code

---

## Session 3: Wrap up

- what we have learned ?
- trade-offs
- managed vs self-hosted

---

## Thank you!
