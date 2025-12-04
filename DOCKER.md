# AISHE Docker Setup

This guide explains how to run AISHE using Docker Compose with all required services (Ollama, Redis, and AISHE API).

## Prerequisites

- Docker (20.10+)
- Docker Compose (1.29+)
- At least 8GB of free disk space (for Ollama models)
- (Optional) NVIDIA GPU with nvidia-docker for GPU acceleration

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

### 3. Test Go Clients

#### Basic Client

```bash
cd clients/go/basic/example
AISHE_API_URL=http://localhost:8000 go run main.go
```

#### Redis Cache Client

```bash
cd clients/go/redis_cache/example
AISHE_API_URL=http://localhost:8000 REDIS_ADDR=localhost:6379 go run main.go
```

#### LangCache Client

Note: LangCache client requires a LangCache service (not included in this Docker setup).

```bash
cd clients/go/langcache/example
export AISHE_API_URL=http://localhost:8000
export LANGCACHE_URL=https://api.langcache.redis.io
export LANGCACHE_API_KEY=your-api-key
export LANGCACHE_CACHE_ID=your-cache-id
go run main.go
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

## GPU Support

The docker-compose.yml includes GPU support for Ollama. If you have an NVIDIA GPU:

1. Install [nvidia-docker](https://github.com/NVIDIA/nvidia-docker)
2. The GPU will be automatically used

If you **don't** have a GPU, comment out the `deploy` section in `docker-compose.yml`:

```yaml
# ollama:
#   ...
#   deploy:
#     resources:
#       reservations:
#         devices:
#           - driver: nvidia
#             count: all
#             capabilities: [gpu]
```

## Environment Variables

You can customize the setup by setting environment variables in `docker-compose.yml`:

### AISHE Service

- `AISHE_SERVER_HOST`: Server host (default: `0.0.0.0`)
- `AISHE_SERVER_PORT`: Server port (default: `8000`)
- `AISHE_OLLAMA_MODEL`: Ollama model to use (default: `llama3.2:3b`)
- `OLLAMA_HOST`: Ollama service URL (default: `http://ollama:11434`)
- `REDIS_ADDR`: Redis address (default: `redis:6379`)

### Example: Using a Different Model

Edit `docker-compose.yml`:

```yaml
aishe:
  environment:
    - AISHE_OLLAMA_MODEL=llama3.2:1b  # Smaller, faster model
```

Then pull the new model:

```bash
docker exec aishe-ollama ollama pull llama3.2:1b
docker-compose restart aishe
```

## Troubleshooting

### Ollama Not Ready

If AISHE shows "ollama_accessible": false:

```bash
# Check Ollama logs
docker-compose logs ollama

# Verify Ollama is running
docker exec aishe-ollama curl http://localhost:11434/api/tags

# Check if model is pulled
docker exec aishe-ollama ollama list
```

### AISHE Container Crashes

```bash
# Check logs
docker-compose logs aishe

# Common issues:
# 1. Ollama not ready - wait longer or check Ollama logs
# 2. Model not pulled - run: docker exec aishe-ollama ollama pull llama3.2:3b
```

### Redis Connection Issues

```bash
# Check Redis is running
docker-compose ps redis

# Test Redis connection
docker exec aishe-redis redis-cli ping
# Should return: PONG
```

### Port Already in Use

If ports 8000, 6379, or 11434 are already in use, edit `docker-compose.yml`:

```yaml
services:
  aishe:
    ports:
      - "8001:8000"  # Use port 8001 instead
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

## Performance Tips

1. **Use GPU**: Significantly faster inference with NVIDIA GPU
2. **Smaller Model**: Use `llama3.2:1b` for faster responses (lower quality)
3. **Increase Resources**: Allocate more CPU/RAM to Docker
4. **SSD Storage**: Store volumes on SSD for better performance

## Next Steps

- Read the [main README](README.md) for more information about AISHE
- Check out the [Go clients documentation](clients/go/README.md)
- Explore the API at http://localhost:8000/docs (FastAPI auto-generated docs)

