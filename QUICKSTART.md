# 🚀 Quick Start Guide

## Prerequisites

- **Docker & Docker Compose**: Latest version
- **Git**: For cloning the repository
- **curl**: For testing APIs
- **jq**: For JSON formatting (optional but recommended)

## 1. Clone & Setup

```bash
# Clone the repository
git clone <repository-url>
cd skool-management

# Quick setup (builds and starts everything)
make setup
```

## 2. Verify Installation

```bash
# Check if all services are healthy
make health

# Or run comprehensive check
make check
```

## 3. Test the API

```bash
# Run automated tests
make test

# Or test manually
./scripts/test.sh
```

## 4. Access the System

Once running, you can access:

- **API Gateway**: http://localhost:8080
- **API Documentation**: See [API.md](API.md)
- **Service Status**: `make info`

## 5. Quick API Test

```bash
# 1. Register a user
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@school.com",
    "password": "admin123",
    "first_name": "Admin",
    "last_name": "User",
    "role": "admin"
  }'

# 2. Login to get token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@school.com",
    "password": "admin123"
  }'

# Save the access_token from response and use it for authenticated requests
```

## 6. Management Commands

```bash
make help          # Show all available commands
make start         # Start services
make stop          # Stop services
make restart       # Restart services
make logs          # View all logs
make test          # Run API tests
make clean-db      # Clean database records only
make clean         # Clean up everything (containers & volumes)
make build         # Build services without starting
make info          # Show service information
```

## 🔍 Monitoring & Health Checks

### Check System Health

```bash
# Check all services and circuit breakers
curl http://localhost:8080/health | jq .

# Quick health check via make command
make health
```

### Monitor Circuit Breakers

```bash
# View circuit breaker states
curl http://localhost:8080/health | jq '.data.*.circuit_breaker'

# Example response:
# {
#   "failure_count": 0,
#   "state": "CLOSED"
# }
```

### Circuit Breaker States

- **CLOSED**: Normal operation ✅
- **HALF_OPEN**: Testing recovery 🟡
- **OPEN**: Service protection active 🔴

## 🗑️ Database Management

### Clean Test Data

```bash
# Remove all records but keep database structure
make clean-db

# Useful for:
# ✅ Starting fresh tests
# ✅ Removing demo data
# ✅ Debugging with clean state
```

### Full Reset

```bash
# Complete cleanup (containers, volumes, images)
make clean
make setup
```

## 🆘 Troubleshooting

### Services won't start

```bash
# Check Docker is running
docker info

# Clean up and restart
make clean
make setup
```

### Port conflicts

Edit `docker-compose.yml` to change ports if needed.

### Database connection issues

```bash
# Check database logs
make logs-mongo     # MongoDB logs
make logs-postgres  # PostgreSQL logs
```

### Need help?

- Check [API.md](API.md) for API documentation
- Check [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment
- Run `make health` to diagnose issues

## 🎯 Next Steps

1. **Explore APIs**: See [examples/README.md](examples/README.md)
2. **Production Deploy**: See [DEPLOYMENT.md](DEPLOYMENT.md)
3. **Development**: Each service can be run individually with `make dev-<service>`

Happy coding! 🎉
