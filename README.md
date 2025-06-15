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

## 🛠️ Technologies

- **Language**: Go 1.21 (Standard Library + minimal dependencies)
- **Authentication**: JWT tokens with refresh mechanism
- **Databases**: MongoDB (Auth), PostgreSQL (School & Student)
- **Communication**: HTTP REST APIs + gRPC (inter-service)
- **Containerization**: Docker & Docker Compose
- **Security**: bcrypt password hashing, JWT tokens

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

- **[API Documentation](API.md)** - Comprehensive API reference
- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment instructions
- **[Examples](examples/README.md)** - API usage examples and test scripts

## 🔧 Development

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional, for convenient commands)

### Development Workflow

```bash
# 1. Clone and setup
git clone <repository-url>
cd skool-management
make setup

# 2. Make changes to services
# Edit files in auth-service/, school-service/, student-service/, or api-gateway/

# 3. Test changes
make restart    # Restart services
make test       # Run tests
make logs       # Check logs

# 4. Individual service development
make dev-auth   # Run auth service locally (connects to Docker databases)
```

### Adding New Features

1. **New API Endpoint**: Add to respective service's main.go
2. **Database Changes**: Add migration files in service/migrations/
3. **gRPC Changes**: Update proto files and regenerate
4. **Authentication**: Use shared JWT middleware

## 🔐 Security Features

- **JWT Authentication**: Access and refresh token mechanism
- **Password Hashing**: bcrypt with salt
- **Rate Limiting**: Built into API Gateway
- **CORS Protection**: Configurable cross-origin policies
- **Database Isolation**: Separate databases per service
- **Secret Management**: Environment-based configuration

## 🌐 API Features

### Authentication Flow

1. User registers via `/auth/signup`
2. User logs in via `/auth/login` → receives access + refresh tokens
3. Access token used for API calls (15min expiry)
4. Refresh token used to get new access tokens (7 day expiry)

### Inter-Service Communication

- Student Service validates school existence via gRPC to School Service
- All services share JWT secret for token validation
- API Gateway handles routing and authentication

### Error Handling

- Consistent error response format across all services
- Proper HTTP status codes
- Detailed error messages for development

## 🚀 Production Deployment

### Docker Compose (Simple)

```bash
# Copy environment template
cp .env.example .env

# Update production values in .env
# Start production stack
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Kubernetes (Advanced)

See [DEPLOYMENT.md](DEPLOYMENT.md) for complete Kubernetes deployment instructions.

## 📈 Monitoring

### Health Checks

Each service exposes `/health` endpoint:

```bash
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # School Service
curl http://localhost:8083/health  # Student Service
```

### Metrics

Services can be extended with Prometheus metrics for monitoring.

### Logging

Structured logging with service identification and request tracing.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make changes and test (`make test`)
4. Commit changes (`git commit -m 'Add amazing feature'`)
5. Push to branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙋‍♂️ Support

- Check [API.md](API.md) for comprehensive API documentation
- Review [examples/](examples/) for usage examples
- Check [DEPLOYMENT.md](DEPLOYMENT.md) for deployment guidance
- Use `make help` for available commands

## ⭐ Features Overview

✅ JWT Authentication with refresh tokens
✅ Microservices architecture
✅ gRPC inter-service communication
✅ Docker containerization
✅ Database per service pattern
✅ API Gateway with rate limiting
✅ Comprehensive error handling
✅ Health checks and monitoring ready
✅ Production deployment guides
✅ Automated testing suite
✅ Development tooling (Makefile, scripts)
✅ API documentation and examples
