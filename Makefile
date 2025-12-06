.PHONY: help docker-setup docker-up docker-down docker-logs docker-test docker-clean docker-rebuild

help:
	@echo "AISHE - Docker Commands"
	@echo "======================="
	@echo ""
	@echo "Setup & Start:"
	@echo "  make docker-setup    - Setup and start all services (recommended first run)"
	@echo "  make docker-up       - Start all services"
	@echo ""
	@echo "Management:"
	@echo "  make docker-down     - Stop all services"
	@echo "  make docker-restart  - Restart all services"
	@echo "  make docker-logs     - View logs from all services"
	@echo "  make docker-logs-f   - Follow logs from all services"
	@echo ""
	@echo "Testing:"
	@echo "  make docker-test     - Run tests against Docker services"
	@echo ""
	@echo "Maintenance:"
	@echo "  make docker-rebuild  - Rebuild AISHE container"
	@echo "  make docker-clean    - Stop and remove all containers and volumes"
	@echo ""
	@echo "Individual Services:"
	@echo "  make logs-aishe      - View AISHE logs"
	@echo "  make logs-ollama     - View Ollama logs"
	@echo "  make logs-redis      - View Redis logs"
	@echo ""

docker-setup:
	@echo "🚀 Setting up AISHE Docker environment..."
	./docker-setup.sh

docker-up:
	@echo "📦 Starting services..."
	docker-compose up -d
	@echo "✅ Services started"
	@echo "   AISHE API: http://localhost:8000"
	@echo "   API Docs:  http://localhost:8000/docs"

docker-down:
	@echo "🛑 Stopping services..."
	docker-compose down
	@echo "✅ Services stopped"

docker-restart:
	@echo "🔄 Restarting services..."
	docker-compose restart
	@echo "✅ Services restarted"

docker-logs:
	docker-compose logs

docker-logs-f:
	docker-compose logs -f

logs-aishe:
	docker-compose logs -f aishe

logs-ollama:
	docker-compose logs -f ollama

logs-redis:
	docker-compose logs -f redis

docker-test:
	@echo "🧪 Running tests..."
	./test-docker.sh

docker-rebuild:
	@echo "🔨 Rebuilding AISHE container..."
	docker-compose build aishe
	docker-compose up -d aishe
	@echo "✅ AISHE container rebuilt"

docker-clean:
	@echo "⚠️  This will remove all containers and volumes (data will be lost)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		echo "✅ Cleaned up"; \
	else \
		echo "❌ Cancelled"; \
	fi

# API testing
test-health:
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8000/health | python3 -m json.tool

test-ask:
	@echo "Testing ask endpoint..."
	@curl -s -X POST http://localhost:8000/api/v1/ask \
		-H "Content-Type: application/json" \
		-d '{"question": "What is Python?"}' | python3 -m json.tool

