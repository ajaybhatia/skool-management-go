package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"skool-management/shared"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

type School struct {
	ID                 int       `json:"id" db:"id"`
	RegistrationNumber string    `json:"registration_number" db:"registration_number"`
	Name               string    `json:"name" db:"name"`
	Address            string    `json:"address" db:"address"`
	Phone              string    `json:"phone" db:"phone"`
	Email              string    `json:"email" db:"email"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSchoolRequest struct {
	RegistrationNumber string `json:"registration_number"`
	Name               string `json:"name"`
	Address            string `json:"address"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
}

type UpdateSchoolRequest struct {
	RegistrationNumber string `json:"registration_number"`
	Name               string `json:"name"`
	Address            string `json:"address"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
}

type SchoolService struct {
	db *sql.DB
}

func NewSchoolService(db *sql.DB) *SchoolService {
	return &SchoolService{db: db}
}

// HTTP Handlers
func (s *SchoolService) CreateSchool(w http.ResponseWriter, r *http.Request) {
	var req CreateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Name == "" {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "School name is required")
		return
	}

	if req.RegistrationNumber == "" {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "School registration number is required")
		return
	}

	query := `
		INSERT INTO schools (registration_number, name, address, phone, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, registration_number, name, address, phone, email, created_at, updated_at
	`

	now := time.Now()
	var school School
	err := s.db.QueryRow(query, req.RegistrationNumber, req.Name, req.Address, req.Phone, req.Email, now, now).Scan(
		&school.ID, &school.RegistrationNumber, &school.Name, &school.Address, &school.Phone, &school.Email,
		&school.CreatedAt, &school.UpdatedAt,
	)

	if err != nil {
		shared.LogError("SCHOOL_SERVICE", "create school", err)
		if strings.Contains(err.Error(), "duplicate key value") {
			shared.WriteErrorResponse(w, http.StatusConflict, "REGISTRATION_EXISTS", "School with this registration number already exists")
			return
		}
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create school")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusCreated, "School created successfully", school)
}

func (s *SchoolService) GetSchools(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, registration_number, name, address, phone, email, created_at, updated_at
		FROM schools
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		shared.LogError("SCHOOL_SERVICE", "get schools", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get schools")
		return
	}
	defer rows.Close()

	var schools []School
	for rows.Next() {
		var school School
		err := rows.Scan(
			&school.ID, &school.RegistrationNumber, &school.Name, &school.Address, &school.Phone, &school.Email,
			&school.CreatedAt, &school.UpdatedAt,
		)
		if err != nil {
			shared.LogError("SCHOOL_SERVICE", "scan school", err)
			continue
		}
		schools = append(schools, school)
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Schools retrieved successfully", schools)
}

func (s *SchoolService) GetSchool(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/schools/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid school ID")
		return
	}

	school, err := s.getSchoolByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", "School not found")
			return
		}
		shared.LogError("SCHOOL_SERVICE", "get school", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get school")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "School retrieved successfully", school)
}

func (s *SchoolService) UpdateSchool(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/schools/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid school ID")
		return
	}

	var req UpdateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RegistrationNumber == "" || req.Name == "" {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Registration number and name are required")
		return
	}

	query := `
		UPDATE schools
		SET registration_number = $1, name = $2, address = $3, phone = $4, email = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, registration_number, name, address, phone, email, created_at, updated_at
	`

	var school School
	err = s.db.QueryRow(query, req.RegistrationNumber, req.Name, req.Address, req.Phone, req.Email, time.Now(), id).Scan(
		&school.ID, &school.RegistrationNumber, &school.Name, &school.Address, &school.Phone, &school.Email,
		&school.CreatedAt, &school.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", "School not found")
			return
		}
		shared.LogError("SCHOOL_SERVICE", "update school", err)
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "schools_registration_number_key") {
				shared.WriteErrorResponse(w, http.StatusConflict, "REGISTRATION_NUMBER_EXISTS", "School with this registration number already exists")
			} else {
				shared.WriteErrorResponse(w, http.StatusConflict, "EMAIL_EXISTS", "School with this email already exists")
			}
			return
		}
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update school")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "School updated successfully", school)
}

func (s *SchoolService) DeleteSchool(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/schools/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid school ID")
		return
	}

	query := `DELETE FROM schools WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		shared.LogError("SCHOOL_SERVICE", "delete school", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete school")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", "School not found")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "School deleted successfully", nil)
}

func (s *SchoolService) getSchoolByID(id int) (*School, error) {
	query := `
		SELECT id, registration_number, name, address, phone, email, created_at, updated_at
		FROM schools
		WHERE id = $1
	`

	var school School
	err := s.db.QueryRow(query, id).Scan(
		&school.ID, &school.RegistrationNumber, &school.Name, &school.Address, &school.Phone, &school.Email,
		&school.CreatedAt, &school.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &school, nil
}

// gRPC Server implementation
type GRPCSchoolServer struct {
	service *SchoolService
}

func (g *GRPCSchoolServer) GetSchool(ctx context.Context, req *GetSchoolRequest) (*GetSchoolResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return &GetSchoolResponse{Found: false}, nil
	}

	school, err := g.service.getSchoolByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &GetSchoolResponse{Found: false}, nil
		}
		return nil, err
	}

	protoSchool := &ProtoSchool{
		Id:                 strconv.Itoa(school.ID),
		RegistrationNumber: school.RegistrationNumber,
		Name:               school.Name,
		Address:            school.Address,
		Phone:              school.Phone,
		Email:              school.Email,
		CreatedAt:          school.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          school.UpdatedAt.Format(time.RFC3339),
	}

	return &GetSchoolResponse{
		School: protoSchool,
		Found:  true,
	}, nil
}

func (g *GRPCSchoolServer) ValidateSchool(ctx context.Context, req *ValidateSchoolRequest) (*ValidateSchoolResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return &ValidateSchoolResponse{Exists: false}, nil
	}

	school, err := g.service.getSchoolByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &ValidateSchoolResponse{Exists: false}, nil
		}
		return &ValidateSchoolResponse{Exists: false}, err
	}

	return &ValidateSchoolResponse{
		Exists: true,
		Name:   school.Name,
	}, nil
}

// Middleware to validate JWT token
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			shared.WriteErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Create JWT manager with same secret as auth service
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your-super-secret-jwt-key-change-in-production"
		}

		jwtManager := shared.NewJWTManager(jwtSecret, "", time.Hour, time.Hour)
		_, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		next(w, r)
	}
}

func main() {
	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "schooluser"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "schoolpass"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "schooldb"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8082"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	// Connect to PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create school service
	schoolService := NewSchoolService(db)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}

		grpcServer := grpc.NewServer()
		RegisterSchoolServiceServer(grpcServer, &GRPCSchoolServer{service: schoolService})
		reflection.Register(grpcServer)

		shared.LogInfo("SCHOOL_SERVICE", fmt.Sprintf("Starting gRPC server on port %s", grpcPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Setup HTTP routes
	http.HandleFunc("/schools", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			authMiddleware(schoolService.GetSchools)(w, r)
		case "POST":
			authMiddleware(schoolService.CreateSchool)(w, r)
		default:
			shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		}
	})

	http.HandleFunc("/schools/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			authMiddleware(schoolService.GetSchool)(w, r)
		case "PUT":
			authMiddleware(schoolService.UpdateSchool)(w, r)
		case "DELETE":
			authMiddleware(schoolService.DeleteSchool)(w, r)
		default:
			shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		}
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteSuccessResponse(w, http.StatusOK, "School service is healthy", nil)
	})

	shared.LogInfo("SCHOOL_SERVICE", fmt.Sprintf("Starting HTTP server on port %s", httpPort))
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
