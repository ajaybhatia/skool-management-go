#!/bin/bash

# School Management Microservices - Database Cleanup Script
# This script removes all records from all databases while keeping the structure

set -e

echo "🗑️  Cleaning all database records..."

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

# Check if services are running
print_step "Checking if services are running..."
if ! docker compose ps | grep -q "Up"; then
    print_error "Services are not running. Please start them first with 'make start'"
    exit 1
fi

print_status "Services are running ✓"

# Clean MongoDB (Auth Service)
print_step "Cleaning MongoDB (Auth Service)..."
docker compose exec mongodb mongosh --eval "
use authdb;
db.users.deleteMany({});
print('Deleted', db.users.countDocuments(), 'users from auth database');
" mongodb://admin:password@localhost:27017/authdb?authSource=admin > /dev/null 2>&1

if [ $? -eq 0 ]; then
    print_status "MongoDB cleaned ✓"
else
    print_warning "Failed to clean MongoDB (it might be empty already)"
fi

# Clean PostgreSQL School Database
print_step "Cleaning PostgreSQL School Database..."
docker compose exec postgres_school psql -U schooluser -d schooldb -c "
DELETE FROM schools;
SELECT 'Deleted all records from schools table' as result;
" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    print_status "School database cleaned ✓"
else
    print_warning "Failed to clean school database (it might be empty already)"
fi

# Clean PostgreSQL Student Database
print_step "Cleaning PostgreSQL Student Database..."
docker compose exec postgres_student psql -U studentuser -d studentdb -c "
DELETE FROM students;
SELECT 'Deleted all records from students table' as result;
" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    print_status "Student database cleaned ✓"
else
    print_warning "Failed to clean student database (it might be empty already)"
fi

# Reset auto-increment sequences for PostgreSQL
print_step "Resetting auto-increment sequences..."
docker compose exec postgres_school psql -U schooluser -d schooldb -c "
ALTER SEQUENCE schools_id_seq RESTART WITH 1;
SELECT 'Reset schools ID sequence' as result;
" > /dev/null 2>&1

docker compose exec postgres_student psql -U studentuser -d studentdb -c "
ALTER SEQUENCE students_id_seq RESTART WITH 1;
SELECT 'Reset students ID sequence' as result;
" > /dev/null 2>&1

print_status "Sequences reset ✓"

echo ""
echo "✅ All database records have been cleaned!"
echo ""
echo "📊 Database Status:"
echo "   📄 MongoDB (Auth):      All users deleted"
echo "   🏫 PostgreSQL (School): All schools deleted, ID sequence reset to 1"
echo "   👥 PostgreSQL (Student): All students deleted, ID sequence reset to 1"
echo ""
echo "💡 The database structures and tables remain intact."
echo "   You can now create fresh data for testing."
echo ""
echo "🧪 To populate with test data:"
echo "   ./scripts/test.sh"
echo ""
