#!/bin/bash
set -e

echo "🚀 Setting up AISHE Docker environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Error: docker-compose is not installed."
    exit 1
fi

echo "✅ Docker is running"

# Start the services
echo "📦 Starting services (Redis, Ollama, AISHE)..."
docker-compose up -d

# Wait for Ollama to be ready
echo "⏳ Waiting for Ollama to be ready..."
timeout=120
elapsed=0
while ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; do
    if [ $elapsed -ge $timeout ]; then
        echo "❌ Timeout waiting for Ollama to start"
        exit 1
    fi
    sleep 2
    elapsed=$((elapsed + 2))
    echo -n "."
done
echo ""
echo "✅ Ollama is ready"

# Pull the model
echo "📥 Pulling llama3.2:3b model (this may take a while on first run)..."
docker exec aishe-ollama ollama pull llama3.2:3b

echo ""
echo "✅ Setup complete!"
echo ""
echo "📊 Service URLs:"
echo "  - AISHE API:  http://localhost:8000"
echo "  - Ollama:     http://localhost:11434"
echo "  - Redis:      localhost:6379"
echo ""
echo "🧪 Test the API:"
echo "  curl http://localhost:8000/health"
echo ""
echo "📝 View logs:"
echo "  docker-compose logs -f aishe"
echo ""
echo "🛑 Stop services:"
echo "  docker-compose down"
echo ""

