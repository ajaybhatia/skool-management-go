# School Management API Examples

This directory contains example API requests for testing the School Management system.

## Prerequisites

1. Make sure the services are running:

   ```bash
   ./scripts/start.sh
   ```

2. Install curl and jq for testing:

   ```bash
   # macOS
   brew install curl jq
   ```

3. (Optional) Start with clean databases:
   ```bash
   # Clean all existing data for fresh testing
   make clean-db
   ```

## Quick Test

Run the automated test script:

```bash
chmod +x scripts/test.sh
./scripts/test.sh
```

## Manual Testing

### 1. User Registration

```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@school.com",
    "password": "admin123",
    "first_name": "Admin",
    "last_name": "User",
    "role": "admin"
  }'
```

### 2. User Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@school.com",
    "password": "admin123"
  }'
```

Save the `access_token` from the response for authenticated requests.

### 3. Create School

```bash
curl -X POST http://localhost:8080/schools \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "registration_number": "REG-SE-001",
    "name": "Springfield Elementary",
    "address": "742 Evergreen Terrace, Springfield",
    "phone": "+1-555-0199",
    "email": "info@springfield-elementary.edu"
  }'
```

### 4. Get All Schools

```bash
curl -X GET http://localhost:8080/schools \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. Create Student

```bash
curl -X POST http://localhost:8080/students \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "roll_number": "001",
    "first_name": "Bart",
    "last_name": "Simpson",
    "email": "bart.simpson@student.com",
    "phone": "+1-555-0200",
    "date_of_birth": "2010-04-01",
    "address": "742 Evergreen Terrace, Springfield",
    "school_id": 1,
    "status": "active"
  }'
```

### 6. Get Students by School

```bash
curl -X GET http://localhost:8080/students/school/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Important Field Requirements

### School Registration Number

- **Required**: Every school must have a unique `registration_number`
- **Constraint**: No two schools can have the same registration number
- **Format**: Recommended format like "REG-SCHOOL-001", but any string is accepted
- **Example**: "REG-SE-001", "SCHOOL-123", "REG-RIVERSIDE-001"

### Student Roll Number

- **Required**: Every student must have a `roll_number`
- **Constraint**: Roll numbers must be unique within each school (but can be the same across different schools)
- **Format**: Any string format is accepted
- **Examples**:
  - Student with roll_number "001" in School 1 ✅
  - Student with roll_number "001" in School 2 ✅ (allowed, different schools)
  - Another student with roll_number "001" in School 1 ❌ (duplicate in same school)

### Error Handling

- **Duplicate school registration_number**: Returns 409 Conflict with error code "REGISTRATION_NUMBER_EXISTS"
- **Duplicate student roll_number in same school**: Returns 409 Conflict with error code "ROLL_NUMBER_EXISTS"
- **Missing required fields**: Returns 400 Bad Request with validation error

## Response Format

All API responses follow this format:

### Success Response

```json
{
  "message": "Operation successful",
  "data": {
    // Response data
  }
}
```

### Error Response

```json
{
  "error": "ERROR_CODE",
  "message": "Human readable error message"
}
```

## Authentication

All endpoints except auth endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

Tokens expire after 15 minutes. Use the refresh endpoint to get a new token:

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Database Management

### Clean Test Data

If you want to start with fresh databases for testing:

```bash
# Remove all users, schools, and students (keeps structure)
make clean-db
```

This is useful when:

- Starting new test scenarios
- Cleaning up after failed tests
- Preparing demo environments
- Debugging with known clean state

### Test Data Recreation

After cleaning, you can recreate test data:

```bash
# Run the automated test (creates sample data)
./scripts/test.sh
```

Or create data manually using the API endpoints above.
