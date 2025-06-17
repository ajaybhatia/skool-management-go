# School Management Microservices

A comprehensive microservices architecture for school management built with Go, featuring JWT authentication, gRPC communication, and Docker containerization.

## 🚀 Quick Start

**New to this project?** Check out our [Quick Start Guide](QUICKSTART.md) for the fastest way to get up and running!

```bash
# One-command setup
make setup

# Check system health
make health

# Run tests
make test
```

## 🏗️ Architecture

```
┌─────────────────┐
│   API Gateway   │ :8080
│   (Rate Limit)  │
└─────────┬───────┘
          │
    ┌─────┼─────┐
    │     │     │
    ▼     ▼     ▼
┌─────┐ ┌─────┐ ┌─────┐
│Auth │ │School│ │Student│
│:8081│ │:8082 │ │:8083│
└─────┘ └─────┘ └─────┘
   │       │       │
   │       └───gRPC─┘
   ▼       ▼       ▼
┌─────┐ ┌─────┐ ┌─────┐
│Mongo│ │PG-1 │ │PG-2 │
│DB   │ │     │ │     │
└─────┘ └─────┘ └─────┘
```

## 🚀 Services

### 1. **Auth Service** (Port 8081)

- JWT-based authentication with refresh tokens
- User registration and login
- MongoDB for user data storage
- Password hashing with bcrypt

### 2. **School Service** (Port 8082)

- CRUD operations for schools
- PostgreSQL database
- gRPC server for inter-service communication
- JWT middleware for authentication

### 3. **Student Service** (Port 8083)

- CRUD operations for students
- School validation via gRPC
- PostgreSQL database (separate from schools)
- Relationship management with schools

### 4. **API Gateway** (Port 8080)

- Central entry point for all requests
- Rate limiting and request routing
- JWT token validation
- Load balancing and service discovery
- **Circuit breaker pattern** for service resilience

## 🛠️ Technologies

- **Language**: Go 1.21 (Standard Library + minimal dependencies)
- **Authentication**: JWT tokens with refresh mechanism
- **Databases**: MongoDB (Auth), PostgreSQL (School & Student)
- **Communication**: HTTP REST APIs + gRPC (inter-service)
- **Containerization**: Docker & Docker Compose
- **Security**: bcrypt password hashing, JWT tokens
- **Resilience**: Circuit breaker pattern for fault tolerance

## 📁 Project Structure

```
skool-management/
├── api-gateway/          # API Gateway service
├── auth-service/         # Authentication service
├── school-service/       # School management service
├── student-service/      # Student management service
├── shared/              # Shared utilities (JWT, utils)
├── proto/               # Protocol Buffers definitions
├── examples/            # API usage examples
├── scripts/             # Automation scripts
├── API.md              # Comprehensive API documentation
├── CIRCUIT_BREAKER.md  # Circuit breaker pattern documentation
├── DEPLOYMENT.md       # Deployment guide
├── docker-compose.yml  # Docker orchestration
└── Makefile           # Development commands
```

## ⚡ Quick Start

### Option 1: Complete Setup (Recommended)

```bash
# Complete development setup
make setup
```

### Option 2: Manual Setup

```bash
# Start all services
make start

# Check service health
make health

# View logs
make logs
```

### Option 3: Individual Services

```bash
# Run individual services for development
make dev-auth     # Auth service
make dev-school   # School service
make dev-student  # Student service
make dev-gateway  # API Gateway
```

## 🧪 Testing

### Automated Testing

```bash
# Run comprehensive API tests
make test
```

### Manual Testing

```bash
# 1. Register a user
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@school.com","password":"admin123","first_name":"Admin","last_name":"User"}'

# 2. Login and get tokens
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@school.com","password":"admin123"}'

# 3. Use the access_token for authenticated requests
export TOKEN="your_access_token_here"

# 4. Create a school
curl -X POST http://localhost:8080/schools \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Springfield Elementary","address":"742 Evergreen Terrace","phone":"+1-555-0199","email":"info@springfield.edu"}'
```

## 📊 Available Commands

```bash
make help          # Show all available commands
make setup         # Complete development setup
make start         # Start all services
make stop          # Stop all services
make restart       # Restart all services
make test          # Run API tests
make logs          # Show all service logs
make clean         # Clean up containers and volumes
make clean-db      # Clean all database records (keep structure)
make build         # Build all services without starting
make info          # Show service information
```

## 🔗 Service URLs

Once running, services are available at:

- **API Gateway**: http://localhost:8080
- **Auth Service**: http://localhost:8081
- **School Service**: http://localhost:8082
- **Student Service**: http://localhost:8083

## 🗄️ Database Management

### Database Access

```bash
# MongoDB (Auth service)
make db-mongo

# PostgreSQL (School service)
make db-school

# PostgreSQL (Student service)
make db-student
```

### Database Cleanup

```bash
# Clean all database records while keeping structure
make clean-db

# What it does:
# ✅ Removes all users from MongoDB
# ✅ Removes all schools from PostgreSQL
# ✅ Removes all students from PostgreSQL
# ✅ Resets auto-increment sequences to 1
# ✅ Keeps database tables and structure intact
```

### Database Reset Options

```bash
make clean-db      # Clean records only (recommended for development)
make clean         # Full reset - removes containers and volumes
```

## 📚 Documentation

- **[Quick Start Guide](QUICKSTART.md)** - Get up and running in minutes
- **[API Documentation](API.md)** - Comprehensive API reference
- **[Circuit Breaker Pattern](CIRCUIT_BREAKER.md)** - Resilience and fault tolerance
- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment instructions
