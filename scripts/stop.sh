#!/bin/bash

# Stop all services
echo "🛑 Stopping School Management Microservices..."

docker compose down -v

echo "✅ All services stopped and volumes removed"
