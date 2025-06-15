#!/bin/bash

# System Health Check and Setup Verification
echo "🏥 School Management System - Health Check"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
CHECKS_PASSED=0
CHECKS_TOTAL=0

# Function to check and report
check_service() {
    local service_name="$1"
    local url="$2"
    local expected_status="${3:-200}"

    CHECKS_TOTAL=$((CHECKS_TOTAL + 1))

    echo -n "Checking $service_name... "

    if response=$(curl -s -w "%{http_code}" -o /dev/null "$url" 2>/dev/null); then
        if [ "$response" = "$expected_status" ]; then
            echo -e "${GREEN}✓ OK${NC} (HTTP $response)"
            CHECKS_PASSED=$((CHECKS_PASSED + 1))
        else
            echo -e "${YELLOW}⚠ Warning${NC} (HTTP $response)"
        fi
    else
        echo -e "${RED}✗ Failed${NC} (No response)"
    fi
}

# Function to check Docker container
check_container() {
    local container_name="$1"
    local service_name="$2"

    CHECKS_TOTAL=$((CHECKS_TOTAL + 1))

    echo -n "Checking $service_name container... "

    if docker ps --format "table {{.Names}}" | grep -q "$container_name"; then
        echo -e "${GREEN}✓ Running${NC}"
        CHECKS_PASSED=$((CHECKS_PASSED + 1))
    else
        echo -e "${RED}✗ Not running${NC}"
    fi
}

echo ""
echo "📋 Prerequisites Check"
echo "----------------------"

# Check Docker
echo -n "Docker installation... "
if command -v docker >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Installed${NC}"
else
    echo -e "${RED}✗ Not installed${NC}"
fi

# Check Docker Compose
echo -n "Docker Compose... "
if docker compose version >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Available${NC}"
else
    echo -e "${RED}✗ Not available${NC}"
fi

# Check Go
echo -n "Go installation... "
if command -v go >/dev/null 2>&1; then
    go_version=$(go version | cut -d' ' -f3)
    echo -e "${GREEN}✓ $go_version${NC}"
else
    echo -e "${YELLOW}⚠ Not installed${NC} (Optional for Docker setup)"
fi

# Check curl
echo -n "curl installation... "
if command -v curl >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Installed${NC}"
else
    echo -e "${RED}✗ Not installed${NC}"
fi

# Check jq
echo -n "jq installation... "
if command -v jq >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Installed${NC}"
else
    echo -e "${YELLOW}⚠ Not installed${NC} (Recommended for testing)"
fi

echo ""
echo "🐳 Container Status"
echo "-------------------"

# Check containers
check_container "school_api_gateway" "API Gateway"
check_container "school_auth_service" "Auth Service"
check_container "school_school_service" "School Service"
check_container "school_student_service" "Student Service"
check_container "school_mongodb" "MongoDB"
check_container "school_postgres_school" "PostgreSQL (School)"
check_container "school_postgres_student" "PostgreSQL (Student)"

echo ""
echo "🌐 Service Health"
echo "------------------"

# Check service endpoints
check_service "API Gateway Health" "http://localhost:8080/health"
check_service "Auth Service Health" "http://localhost:8081/health"
check_service "School Service Health" "http://localhost:8082/health"
check_service "Student Service Health" "http://localhost:8083/health"

echo ""
echo "🔗 API Endpoints"
echo "----------------"

# Check API endpoints (these might return 401 for auth, which is expected)
check_service "Schools API" "http://localhost:8080/schools" "401"
check_service "Students API" "http://localhost:8080/students" "401"
check_service "Auth Signup" "http://localhost:8080/auth/signup" "400"  # Bad request without body is expected

echo ""
echo "📊 Summary"
echo "----------"

if [ $CHECKS_PASSED -eq $CHECKS_TOTAL ]; then
    echo -e "${GREEN}✅ All checks passed ($CHECKS_PASSED/$CHECKS_TOTAL)${NC}"
    echo -e "${GREEN}🎉 System is healthy and ready to use!${NC}"
    exit 0
elif [ $CHECKS_PASSED -gt $((CHECKS_TOTAL / 2)) ]; then
    echo -e "${YELLOW}⚠️  Most checks passed ($CHECKS_PASSED/$CHECKS_TOTAL)${NC}"
    echo -e "${YELLOW}🔧 System is partially working - check warnings above${NC}"
    exit 1
else
    echo -e "${RED}❌ Multiple checks failed ($CHECKS_PASSED/$CHECKS_TOTAL)${NC}"
    echo -e "${RED}🚨 System needs attention - check errors above${NC}"

    echo ""
    echo "🛠️  Troubleshooting:"
    echo "   1. Make sure Docker is running"
    echo "   2. Run: make start"
    echo "   3. Wait a few minutes for services to start"
    echo "   4. Run this health check again"

    exit 2
fi
