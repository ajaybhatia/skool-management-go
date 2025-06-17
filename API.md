# School Management API Documentation

## Overview

The School Management System is a microservices-based application that provides APIs for managing schools, students, and user authentication. All services communicate through an API Gateway that handles authentication, routing, and centralized logging.

## Architecture

```
Client → API Gateway → [Auth Service, School Service, Student Service]
                    ↓
                [MongoDB, PostgreSQL-1, PostgreSQL-2]
```

## Base URL

All API requests should be made to the API Gateway:

```
http://localhost:8080
```

## Authentication

### JWT Token Authentication

Most endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

Tokens expire after 15 minutes. Use the refresh endpoint to get a new token.

## API Endpoints

### Authentication Service

#### POST /auth/signup

Register a new user.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "role": "admin" // optional, defaults to "user"
}
```

**Response:**

```json
{
  "message": "User created successfully",
  "data": {
    "id": "user_id",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "admin",
    "created_at": "2025-06-15T10:00:00Z",
    "updated_at": "2025-06-15T10:00:00Z"
  }
}
```

#### POST /auth/login

Authenticate a user and receive JWT tokens.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "message": "Login successful",
  "data": {
    "user": {
      "id": "user_id",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin"
    },
    "access_token": "jwt_access_token",
    "refresh_token": "jwt_refresh_token"
  }
}
```

#### POST /auth/refresh

Refresh an expired access token.

**Request Body:**

```json
{
  "refresh_token": "your_refresh_token"
}
```

**Response:**

```json
{
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "new_jwt_access_token"
  }
}
```

### School Service

All school endpoints require authentication.

#### GET /schools

Get all schools.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
```

**Response:**

```json
{
  "message": "Schools retrieved successfully",
  "data": [
    {
      "id": 1,
      "registration_number": "REG-SE-001",
      "name": "Springfield Elementary",
      "address": "742 Evergreen Terrace, Springfield",
      "phone": "+1-555-0199",
      "email": "info@springfield-elementary.edu",
      "created_at": "2025-06-15T10:00:00Z",
      "updated_at": "2025-06-15T10:00:00Z"
    }
  ]
}
```

#### POST /schools

Create a new school.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "registration_number": "REG-SE-001",
  "name": "Springfield Elementary",
  "address": "742 Evergreen Terrace, Springfield",
  "phone": "+1-555-0199",
  "email": "info@springfield-elementary.edu"
}
```

**Response:**

```json
{
  "message": "School created successfully",
  "data": {
    "id": 1,
    "registration_number": "REG-SE-001",
    "name": "Springfield Elementary",
    "address": "742 Evergreen Terrace, Springfield",
    "phone": "+1-555-0199",
    "email": "info@springfield-elementary.edu",
    "created_at": "2025-06-15T10:00:00Z",
    "updated_at": "2025-06-15T10:00:00Z"
  }
}
```

#### GET /schools/{id}

Get a specific school by ID.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
```

**Response:**

```json
{
  "message": "School retrieved successfully",
  "data": {
    "id": 1,
    "registration_number": "REG-SE-001",
    "name": "Springfield Elementary",
    "address": "742 Evergreen Terrace, Springfield",
    "phone": "+1-555-0199",
    "email": "info@springfield-elementary.edu",
    "created_at": "2025-06-15T10:00:00Z",
    "updated_at": "2025-06-15T10:00:00Z"
  }
}
```

#### PUT /schools/{id}

Update a school.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "registration_number": "REG-SE-001",
  "name": "Springfield Elementary School",
  "address": "742 Evergreen Terrace, Springfield, IL",
  "phone": "+1-555-0199",
  "email": "contact@springfield-elementary.edu"
}
```

#### DELETE /schools/{id}

Delete a school.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
```

**Response:**

```json
{
  "message": "School deleted successfully"
}
```

### Student Service

All student endpoints require authentication.

#### GET /students

Get all students.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
```

**Response:**

```json
{
  "message": "Students retrieved successfully",
  "data": [
    {
      "id": 1,
      "roll_number": "001",
      "first_name": "Bart",
      "last_name": "Simpson",
      "email": "bart.simpson@student.com",
      "phone": "+1-555-0200",
      "date_of_birth": "2010-04-01",
      "address": "742 Evergreen Terrace, Springfield",
      "school_id": 1,
      "school_name": "Springfield Elementary",
      "enrollment_date": "2025-01-15",
      "status": "active",
      "created_at": "2025-06-15T10:00:00Z",
      "updated_at": "2025-06-15T10:00:00Z"
    }
  ]
}
```

#### POST /students

Create a new student.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "roll_number": "001",
  "first_name": "Bart",
  "last_name": "Simpson",
  "email": "bart.simpson@student.com",
  "phone": "+1-555-0200",
  "date_of_birth": "2010-04-01",
  "address": "742 Evergreen Terrace, Springfield",
  "school_id": 1,
  "enrollment_date": "2025-01-15", // optional, defaults to current date
  "status": "active" // optional, defaults to "active"
}
```

**Response:**

```json
{
  "message": "Student created successfully",
  "data": {
    "id": 1,
    "roll_number": "001",
    "first_name": "Bart",
    "last_name": "Simpson",
    "email": "bart.simpson@student.com",
    "phone": "+1-555-0200",
    "date_of_birth": "2010-04-01",
    "address": "742 Evergreen Terrace, Springfield",
    "school_id": 1,
    "enrollment_date": "2025-01-15",
    "status": "active",
    "created_at": "2025-06-15T10:00:00Z",
    "updated_at": "2025-06-15T10:00:00Z"
  }
}
```

#### GET /students/{id}

Get a specific student by ID.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
```

**Response:**

```json
{
  "message": "Student retrieved successfully",
  "data": {
    "id": 1,
    "roll_number": "001",
    "first_name": "Bart",
    "last_name": "Simpson",
    "email": "bart.simpson@student.com",
    "phone": "+1-555-0200",
    "date_of_birth": "2010-04-01",
    "address": "742 Evergreen Terrace, Springfield",
    "school_id": 1,
    "enrollment_date": "2025-01-15",
    "status": "active",
    "created_at": "2025-06-15T10:00:00Z",
    "updated_at": "2025-06-15T10:00:00Z"
  }
}
```

#### PUT /students/{id}

Update a student.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "roll_number": "002",
  "first_name": "Bart",
  "last_name": "Simpson-Updated",
  "email": "bart.simpson@student.com",
  "phone": "+1-555-0201",
  "date_of_birth": "2010-04-01",
  "address": "742 Evergreen Terrace, Springfield",
  "school_id": 1,
  "enrollment_date": "2025-01-15",
  "status": "active"
}
```

#### DELETE /students/{id}

Delete a student.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
```

**Response:**

```json
{
  "message": "Student deleted successfully"
}
```

#### GET /students/school/{school_id}

Get all students for a specific school.

**Headers:**

```
Authorization: Bearer <your_jwt_token>
```

**Response:**

```json
{
  "message": "Students retrieved successfully",
  "data": [
    {
      "id": 1,
      "roll_number": "001",
      "first_name": "Bart",
      "last_name": "Simpson",
      "school_id": 1,
      "school_name": "Springfield Elementary"
      // ... other student fields
    }
  ]
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": "ERROR_CODE",
  "message": "Human readable error message"
}
```

### Common Error Codes

- `INVALID_REQUEST` - Request body is malformed
- `VALIDATION_ERROR` - Required fields are missing
- `UNAUTHORIZED` - Invalid or missing authentication token
- `FORBIDDEN` - User doesn't have permission
- `NOT_FOUND` - Resource not found
- `CONFLICT` - Resource already exists
- `INTERNAL_ERROR` - Server error

### HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## Rate Limiting

The API Gateway implements rate limiting:

- 100 requests per minute per IP address
- 1000 requests per hour per authenticated user

## Circuit Breaker Protection

The system implements circuit breaker patterns for enhanced resilience:

### Circuit Breaker Endpoints

#### GET /health

Returns system health status including circuit breaker states:

**Response:**

```json
{
  "message": "All services are healthy",
  "data": {
    "auth": {
      "circuit_breaker": {
        "failure_count": 0,
        "state": "CLOSED"
      },
      "status": "healthy"
    },
    "gateway": {
      "status": "healthy"
    },
    "school": {
      "circuit_breaker": {
        "failure_count": 0,
        "state": "CLOSED"
      },
      "status": "healthy"
    },
    "student": {
      "circuit_breaker": {
        "failure_count": 0,
        "state": "CLOSED"
      },
      "status": "healthy"
    }
  }
}
```

### Circuit Breaker Behavior

| State         | Behavior         | Response                            |
| ------------- | ---------------- | ----------------------------------- |
| **CLOSED**    | Normal operation | Requests pass through               |
| **HALF_OPEN** | Testing recovery | Limited requests allowed            |
| **OPEN**      | Service failure  | Immediate rejection with 503 status |

### Circuit Breaker Responses

When a circuit breaker is open:

**Status Code:** `503 Service Unavailable`

**Response:**

```json
{
  "error": "CIRCUIT_BREAKER_OPEN",
  "message": "Service is temporarily unavailable due to circuit breaker"
}
```

### Configuration

| Service             | Max Failures | Reset Timeout |
| ------------------- | ------------ | ------------- |
| HTTP Services       | 5            | 60 seconds    |
| gRPC Services       | 3            | 30 seconds    |
| Database Operations | 5            | 60 seconds    |

## Inter-Service Communication

Services communicate using gRPC for internal operations:

- Student Service validates school existence via School Service gRPC
- All services use the same JWT secret for token validation

## Data Validation

### Email Format

Must be a valid email address format.

### Password Requirements

- Minimum 6 characters
- No maximum length limit

### School ID

Must be a positive integer and must exist in the school database.

### Date Format

Dates should be in YYYY-MM-DD format.

## Unique Constraints

### School Registration Numbers

- Each school must have a unique `registration_number`
- Attempting to create or update a school with an existing registration number returns:
  ```json
  {
    "error": "REGISTRATION_NUMBER_EXISTS",
    "message": "School with this registration number already exists"
  }
  ```

### Student Roll Numbers

- Each student must have a unique `roll_number` within their school
- Students in different schools can have the same roll number
- Attempting to create or update a student with an existing roll number in the same school returns:
  ```json
  {
    "error": "ROLL_NUMBER_EXISTS",
    "message": "Student with this roll number already exists in this school"
  }
  ```

### Email Addresses

- Student email addresses must be globally unique across all schools
- Attempting to create or update a student with an existing email returns:
  ```json
  {
    "error": "EMAIL_EXISTS",
    "message": "Student with this email already exists"
  }
  ```

## Development

### Running Locally

```bash
# Start all services
make setup

# Run individual services in development mode
make dev-auth
make dev-school
make dev-student
make dev-gateway
```

### Testing

```bash
# Run API tests
make test

# Check service health
make health
```

### Database Access

```bash
# Connect to MongoDB (Auth)
make db-mongo

# Connect to School PostgreSQL
make db-school

# Connect to Student PostgreSQL
make db-student
```

### Database Management

```bash
# Clean all database records (keeps structure)
make clean-db

# Full reset (containers and volumes)
make clean

# What clean-db does:
# ✅ Removes all users from MongoDB
# ✅ Removes all schools from PostgreSQL
# ✅ Removes all students from PostgreSQL
# ✅ Resets auto-increment sequences
# ✅ Preserves database tables and schemas
```
