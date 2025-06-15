#!/bin/bash

# Test the API endpoints
echo "🧪 Testing School Management API..."

API_BASE="http://localhost:8080"

# Test API Gateway health
echo "1️⃣  Testing API Gateway health..."
curl -s "$API_BASE/health" | jq '.' || echo "Health check failed"
echo ""

# Test user registration
echo "2️⃣  Testing user registration..."
SIGNUP_RESPONSE=$(curl -s -X POST "$API_BASE/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User",
    "role": "admin"
  }')
echo "$SIGNUP_RESPONSE" | jq '.' || echo "Signup failed"
echo ""

# Test user login
echo "3️⃣  Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
echo "$LOGIN_RESPONSE" | jq '.' || echo "Login failed"

# Extract token from login response
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token // empty')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Failed to get access token. Cannot proceed with authenticated tests."
  exit 1
fi

echo "✅ Got access token: ${TOKEN:0:20}..."
echo ""

# Test school creation
echo "4️⃣  Testing school creation..."
SCHOOL_RESPONSE=$(curl -s -X POST "$API_BASE/schools" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "registration_number": "REG-THS-001",
    "name": "Test High School",
    "address": "123 Education St, Learning City",
    "phone": "+1-555-0123",
    "email": "contact@testhighschool.edu"
  }')
echo "$SCHOOL_RESPONSE" | jq '.' || echo "School creation failed"

# Extract school ID
SCHOOL_ID=$(echo "$SCHOOL_RESPONSE" | jq -r '.data.id // empty')
echo "✅ Created school with ID: $SCHOOL_ID"
echo ""

# Test student creation
echo "5️⃣  Testing student creation..."
STUDENT_RESPONSE=$(curl -s -X POST "$API_BASE/students" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"roll_number\": \"001\",
    \"first_name\": \"John\",
    \"last_name\": \"Doe\",
    \"email\": \"john.doe@student.com\",
    \"phone\": \"+1-555-0456\",
    \"date_of_birth\": \"2005-03-15\",
    \"address\": \"456 Student Ave, Learning City\",
    \"school_id\": $SCHOOL_ID,
    \"status\": \"active\"
  }")
echo "$STUDENT_RESPONSE" | jq '.' || echo "Student creation failed"
echo ""

# Test getting all schools
echo "6️⃣  Testing get all schools..."
curl -s -X GET "$API_BASE/schools" \
  -H "Authorization: Bearer $TOKEN" | jq '.' || echo "Get schools failed"
echo ""

# Test getting all students
echo "7️⃣  Testing get all students..."
curl -s -X GET "$API_BASE/students" \
  -H "Authorization: Bearer $TOKEN" | jq '.' || echo "Get students failed"
echo ""

# Test getting students by school
echo "8️⃣  Testing get students by school..."
curl -s -X GET "$API_BASE/students/school/$SCHOOL_ID" \
  -H "Authorization: Bearer $TOKEN" | jq '.' || echo "Get students by school failed"
echo ""

echo "🎉 API testing completed!"
