#!/bin/bash
set -e

echo "🧪 Testing AISHE Docker Setup"
echo "=============================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local name=$1
    local url=$2
    local expected_status=${3:-200}
    
    echo -n "Testing $name... "
    
    status=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")
    
    if [ "$status" = "$expected_status" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $status)"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $status, expected $expected_status)"
        return 1
    fi
}

# Check if services are running
echo "1. Checking if Docker services are running..."
if ! docker-compose ps | grep -q "aishe-server.*Up"; then
    echo -e "${RED}✗ AISHE server is not running${NC}"
    echo "Run: ./docker-setup.sh"
    exit 1
fi
echo -e "${GREEN}✓ Services are running${NC}"
echo ""

# Test Redis
echo "2. Testing Redis..."
if docker exec aishe-redis redis-cli ping | grep -q "PONG"; then
    echo -e "${GREEN}✓ Redis is responding${NC}"
else
    echo -e "${RED}✗ Redis is not responding${NC}"
    exit 1
fi
echo ""

# Test Ollama
echo "3. Testing Ollama..."
if curl -s http://localhost:11434/api/tags > /dev/null; then
    echo -e "${GREEN}✓ Ollama is responding${NC}"

    # Check if model is pulled
    if docker exec aishe-ollama ollama list | grep -q "llama3.2:3b"; then
        echo -e "${GREEN}✓ Model llama3.2:3b is available${NC}"
    else
        echo -e "${YELLOW}⚠ Model llama3.2:3b not found${NC}"
        echo "Run: docker exec aishe-ollama ollama pull llama3.2:3b"
    fi
else
    echo -e "${RED}✗ Ollama is not responding${NC}"
    exit 1
fi
echo ""

# Test AISHE API
echo "4. Testing AISHE API..."
test_endpoint "Health endpoint" "http://localhost:8000/health"
test_endpoint "Root endpoint" "http://localhost:8000/"
echo ""

# Test question answering
echo "5. Testing question answering..."
echo -n "Asking a question... "

response=$(curl -s -X POST http://localhost:8000/api/v1/ask \
    -H "Content-Type: application/json" \
    -d '{"question": "What is Python?"}' 2>/dev/null)

if echo "$response" | grep -q '"answer"'; then
    echo -e "${GREEN}✓ PASS${NC}"
    echo ""
    echo "Response preview:"
    echo "$response" | python3 -m json.tool 2>/dev/null | head -20
else
    echo -e "${RED}✗ FAIL${NC}"
    echo "Response: $response"
    exit 1
fi
echo ""

# Test Go workshop clients
echo "6. Testing Go workshop clients (optional)..."
if command -v go &> /dev/null; then
    echo "Testing Session 1 (Basic client)..."
    cd clients/go/workshop/session1-basic/solution
    if timeout 30 go run . > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Session 1 works${NC}"
    else
        echo -e "${YELLOW}⚠ Session 1 test failed or timed out${NC}"
    fi
    cd - > /dev/null

    echo "Testing Session 2 (Redis cache client)..."
    cd clients/go/workshop/session2-redis/solution
    if timeout 30 go run . > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Session 2 works${NC}"
    else
        echo -e "${YELLOW}⚠ Session 2 test failed or timed out${NC}"
    fi
    cd - > /dev/null
else
    echo -e "${YELLOW}⚠ Go not installed, skipping workshop tests${NC}"
fi
echo ""

echo "=============================="
echo -e "${GREEN}✅ All tests passed!${NC}"
echo ""
echo "📊 Service URLs:"
echo "  - AISHE API:  http://localhost:8000"
echo "  - API Docs:   http://localhost:8000/docs"
echo "  - Ollama:     http://localhost:11434"
echo "  - Redis:      localhost:6379"
echo ""

