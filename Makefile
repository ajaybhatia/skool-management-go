# School Management Microservices Makefile

.PHONY: help setup start stop restart logs test clean build

# Default target
help:
	@echo "School Management Microservices"
	@echo "Available commands:"
	@echo "  setup     - Complete development setup (builds and starts all services)"
	@echo "  start     - Start all services"
	@echo "  stop      - Stop all services"
	@echo "  restart   - Restart all services"
	@echo "  logs      - Show logs from all services"
	@echo "  test      - Run API tests"
	@echo "  clean     - Clean up containers and volumes"
	@echo "  clean-db  - Clean all database records (keep structure)"
	@echo "  build     - Build all services without starting"
	@echo "  info      - Show service information"
	@echo ""
	@echo "Individual service commands:"
	@echo "  logs-auth     - Show auth service logs"
	@echo "  logs-school   - Show school service logs"
	@echo "  logs-student  - Show student service logs"
	@echo "  logs-gateway  - Show API gateway logs"

# Complete development setup
setup:
	@echo "🚀 Setting up School Management Microservices..."
	./scripts/dev-setup.sh

# Start services
start:
	@echo "▶️  Starting services..."
	./scripts/start.sh

# Stop services
stop:
	@echo "⏹️  Stopping services..."
	./scripts/stop.sh

# Restart services
restart: stop start

# Show logs
logs:
	docker compose logs -f

# Show logs for individual services
logs-auth:
	docker compose logs -f auth-service

logs-school:
	docker compose logs -f school-service

logs-student:
	docker compose logs -f student-service

logs-gateway:
	docker compose logs -f api-gateway

# Run tests
test:
	@echo "🧪 Running API tests..."
	./scripts/test.sh

# Clean up
clean:
	@echo "🧹 Cleaning up..."
	docker compose down -v
	docker system prune -f
	docker volume prune -f

# Clean all database records (keep structure)
clean-db:
	@echo "🗑️  Cleaning all database records..."
	./scripts/clean-db.sh

# Build without starting
build:
	@echo "🔨 Building services..."
	docker compose build

# Show service info
info:
	@echo "📋 Service Information:"
	./scripts/info.sh

# Development helpers
dev-auth:
	@echo "🔧 Running auth service in development mode..."
	cd auth-service && go run main.go

dev-school:
	@echo "🔧 Running school service in development mode..."
	cd school-service && go run main.go grpc.go

dev-student:
	@echo "🔧 Running student service in development mode..."
	cd student-service && go run main.go grpc.go

dev-gateway:
	@echo "🔧 Running API gateway in development mode..."
	cd api-gateway && go run main.go

# Database helpers
db-mongo:
	@echo "📊 Connecting to MongoDB..."
	docker exec -it school_mongodb mongosh mongodb://admin:password@localhost:27017/authdb

db-school:
	@echo "📊 Connecting to School PostgreSQL..."
	docker exec -it school_postgres_school psql -U schooluser -d schooldb

db-student:
	@echo "📊 Connecting to Student PostgreSQL..."
	docker exec -it school_postgres_student psql -U studentuser -d studentdb

# Go module management
deps:
	@echo "📦 Managing Go dependencies..."
	go mod tidy
	go mod download

# Generate protobuf (for future use)
proto:
	@echo "🔄 Generating protobuf files..."
	@echo "Note: This requires protoc and protoc-gen-go to be installed"
	protoc --go_out=. --go-grpc_out=. proto/school.proto

# Health check
health:
	@echo "🏥 Checking service health..."
	./scripts/health-check.sh

# Complete health check (alternative to basic health)
check:
	@echo "🔍 Running comprehensive health check..."
	./scripts/health-check.sh
