#!/bin/bash

# Build and start all services
echo "🚀 Starting School Management Microservices..."

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is not available. Please install Docker Compose first."
    exit 1
fi

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker compose down -v

# Build and start all services
echo "🔨 Building and starting services..."
docker compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check service health
echo "🔍 Checking service health..."

# Check API Gateway
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ API Gateway is healthy"
else
    echo "❌ API Gateway is not responding"
fi

# Check Auth Service
if curl -s http://localhost:8081/health > /dev/null; then
    echo "✅ Auth Service is healthy"
else
    echo "❌ Auth Service is not responding"
fi

# Check School Service
if curl -s http://localhost:8082/health > /dev/null; then
    echo "✅ School Service is healthy"
else
    echo "❌ School Service is not responding"
fi

# Check Student Service
if curl -s http://localhost:8083/health > /dev/null; then
    echo "✅ Student Service is healthy"
else
    echo "❌ Student Service is not responding"
fi

echo ""
echo "🎉 School Management System is now running!"
echo ""
echo "📚 Available endpoints:"
echo "  🌐 API Gateway:     http://localhost:8080"
echo "  📖 API Docs:       http://localhost:8080/docs"
echo "  ❤️  Health Check:   http://localhost:8080/health"
echo "  🔐 Auth Service:    http://localhost:8081"
echo "  🏫 School Service:  http://localhost:8082"
echo "  👥 Student Service: http://localhost:8083"
echo ""
echo "🔍 To view logs: docker compose logs -f [service-name]"
echo "🛑 To stop: docker compose down"
echo "📋 To view this info again: ./scripts/info.sh"
