#!/bin/bash

# School Management Microservices - Development Setup
# This script sets up the development environment

set -e

echo "🚀 Setting up School Management Microservices..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if Docker is running
print_step "Checking Docker..."
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi
print_status "Docker is running ✓"

# Check if Docker Compose is available
if ! docker compose version > /dev/null 2>&1; then
    print_error "Docker Compose is not available. Please install Docker Compose and try again."
    exit 1
fi
print_status "Docker Compose is available ✓"

# Clean up any existing containers
print_step "Cleaning up existing containers..."
docker compose down -v 2>/dev/null || true
docker system prune -f > /dev/null 2>&1 || true
print_status "Cleanup completed ✓"

# Initialize Go modules
print_step "Initializing Go modules..."
if [ ! -f "go.sum" ]; then
    go mod tidy
    print_status "Go modules initialized ✓"
else
    print_status "Go modules already initialized ✓"
fi

# Build and start services
print_step "Building and starting services..."
docker compose up --build -d

# Wait for services to be ready
print_step "Waiting for services to be ready..."
sleep 10

# Check service health
print_step "Checking service health..."

services=(
    "http://localhost:8080/health:API Gateway"
    "http://localhost:8081/health:Auth Service"
    "http://localhost:8082/health:School Service"
    "http://localhost:8083/health:Student Service"
)

all_healthy=true

for service in "${services[@]}"; do
    url="${service%%:*}"
    name="${service##*:}"

    if curl -f -s "$url" > /dev/null 2>&1; then
        print_status "$name is healthy ✓"
    else
        print_warning "$name is not responding"
        all_healthy=false
    fi
done

# Show service URLs
echo ""
echo "🎉 School Management Microservices are running!"
echo ""
echo "📋 Service URLs:"
echo "   API Gateway:    http://localhost:8080"
echo "   Auth Service:   http://localhost:8081"
echo "   School Service: http://localhost:8082"
echo "   Student Service: http://localhost:8083"
echo ""
echo "🗄️  Database URLs:"
echo "   MongoDB (Auth):     mongodb://admin:password@localhost:27017/authdb"
echo "   PostgreSQL (School): postgresql://schooluser:schoolpass@localhost:5432/schooldb"
echo "   PostgreSQL (Student): postgresql://studentuser:studentpass@localhost:5433/studentdb"
echo ""
echo "📚 Documentation:"
echo "   API Examples: ./examples/README.md"
echo "   Project README: ./README.md"
echo ""

if [ "$all_healthy" = true ]; then
    echo "✅ All services are healthy and ready to use!"
    echo ""
    echo "🧪 To run tests:"
    echo "   ./scripts/test.sh"
    echo ""
    echo "🛑 To stop services:"
    echo "   ./scripts/stop.sh"
else
    echo "⚠️  Some services are not responding. Check logs with:"
    echo "   docker compose logs"
fi

echo ""
echo "📖 Quick start:"
echo "   1. Register a user: curl -X POST http://localhost:8080/auth/signup -H 'Content-Type: application/json' -d '{\"email\":\"test@example.com\",\"password\":\"password123\",\"first_name\":\"Test\",\"last_name\":\"User\"}'"
echo "   2. Login: curl -X POST http://localhost:8080/auth/login -H 'Content-Type: application/json' -d '{\"email\":\"test@example.com\",\"password\":\"password123\"}'"
echo "   3. Use the access_token for authenticated requests"
echo ""
