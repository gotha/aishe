# AISHE Docker Setup

This guide explains how to run AISHE using Docker Compose with all required services (Ollama, Redis, and AISHE API).

## Prerequisites

- Docker (20.10+)
- Docker Compose (1.29+)
- At least 8GB of free disk space (for Ollama models)

## Quick Start

### 1. Automated Setup (Recommended)

Run the setup script to start all services and pull the Ollama model:

```bash
./docker-setup.sh
```

This script will:
- Start Redis, Ollama, and AISHE containers
- Wait for services to be ready
- Pull the `llama3.2:3b` model
- Display service URLs and helpful commands

### 2. Manual Setup

If you prefer manual control:

```bash
# Start all services
docker-compose up -d

# Wait for Ollama to be ready (check logs)
docker-compose logs -f ollama

# Pull the model (in another terminal)
docker exec aishe-ollama ollama pull llama3.2:3b

# Check AISHE logs
docker-compose logs -f aishe
```

## Services

The Docker Compose setup includes three services:

### 1. Redis (`aishe-redis`)
- **Port**: 6379
- **Purpose**: Caching for the Redis cache client
- **Data**: Persisted in `redis-data` volume

### 2. Ollama (`aishe-ollama`)
- **Port**: 11434
- **Purpose**: LLM inference
- **Model**: llama3.2:3b (pulled during setup)
- **Data**: Persisted in `ollama-data` volume

### 3. AISHE API (`aishe-server`)
- **Port**: 8000
- **Purpose**: Main API server
- **Endpoints**:
  - `GET /health` - Health check
  - `POST /api/v1/ask` - Ask questions

## Testing the Setup

### 1. Health Check

```bash
curl http://localhost:8000/health
```

Expected response:
```json
{
  "status": "healthy",
  "ollama_accessible": true
}
```

### 2. Ask a Question

```bash
curl -X POST http://localhost:8000/api/v1/ask \
  -H "Content-Type: application/json" \
  -d '{"question": "What is the capital of France?"}'
```

## Managing Services

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f aishe
docker-compose logs -f ollama
docker-compose logs -f redis
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart aishe
```

### Stop Services

```bash
# Stop all services (keeps data)
docker-compose stop

# Stop and remove containers (keeps data)
docker-compose down

# Stop and remove everything including volumes (DELETES DATA)
docker-compose down -v
```

### Rebuild AISHE Container

If you make changes to the code:

```bash
docker-compose build aishe
docker-compose up -d aishe
```

## Data Persistence

Data is persisted in Docker volumes:

- `redis-data`: Redis cache data
- `ollama-data`: Ollama models and data

To view volumes:

```bash
docker volume ls | grep aishe
```

To backup volumes:

```bash
# Backup Ollama models
docker run --rm -v aishe_ollama-data:/data -v $(pwd):/backup alpine tar czf /backup/ollama-backup.tar.gz -C /data .

# Backup Redis data
docker run --rm -v aishe_redis-data:/data -v $(pwd):/backup alpine tar czf /backup/redis-backup.tar.gz -C /data .
```

## Next Steps

- Read the [main README](README.md) for more information about AISHE
- Explore the API at http://localhost:8000/docs (FastAPI auto-generated docs)

