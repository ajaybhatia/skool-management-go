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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

type Student struct {
	ID             int       `json:"id" db:"id"`
	RollNumber     string    `json:"roll_number" db:"roll_number"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	Email          string    `json:"email" db:"email"`
	Phone          string    `json:"phone" db:"phone"`
	DateOfBirth    string    `json:"date_of_birth" db:"date_of_birth"`
	Address        string    `json:"address" db:"address"`
	SchoolID       int       `json:"school_id" db:"school_id"`
	SchoolName     string    `json:"school_name,omitempty"`
	EnrollmentDate string    `json:"enrollment_date" db:"enrollment_date"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateStudentRequest struct {
	RollNumber     string `json:"roll_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DateOfBirth    string `json:"date_of_birth"`
	Address        string `json:"address"`
	SchoolID       int    `json:"school_id"`
	EnrollmentDate string `json:"enrollment_date"`
	Status         string `json:"status"`
}

type UpdateStudentRequest struct {
	RollNumber     string `json:"roll_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DateOfBirth    string `json:"date_of_birth"`
	Address        string `json:"address"`
	SchoolID       int    `json:"school_id"`
	EnrollmentDate string `json:"enrollment_date"`
	Status         string `json:"status"`
}

type StudentService struct {
	db                *sql.DB
	schoolServiceConn *grpc.ClientConn
}

func NewStudentService(db *sql.DB, schoolServiceConn *grpc.ClientConn) *StudentService {
	return &StudentService{
		db:                db,
		schoolServiceConn: schoolServiceConn,
	}
}

// Helper function to validate school existence via gRPC
func (s *StudentService) validateSchool(schoolID int) (bool, string, error) {
	if s.schoolServiceConn == nil {
		return true, "", nil // Skip validation if gRPC connection is not available
	}

	client := NewSchoolServiceClient(s.schoolServiceConn)
	resp, err := client.ValidateSchool(context.Background(), &ValidateSchoolRequest{
		Id: strconv.Itoa(schoolID),
	})
	if err != nil {
		return false, "", err
	}

	return resp.Exists, resp.Name, nil
}

func (s *StudentService) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var req CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RollNumber == "" || req.FirstName == "" || req.LastName == "" || req.SchoolID == 0 {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Roll number, first name, last name, and school ID are required")
		return
	}

	// Validate school exists
	schoolExists, schoolName, err := s.validateSchool(req.SchoolID)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school validation", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to validate school")
		return
	}

	if !schoolExists {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SCHOOL", "School does not exist")
		return
	}

	// Set defaults
	if req.Status == "" {
		req.Status = "active"
	}
	if req.EnrollmentDate == "" {
		req.EnrollmentDate = time.Now().Format("2006-01-02")
	}

	query := `
		INSERT INTO students (roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
	`

	now := time.Now()
	var student Student
	err = s.db.QueryRow(query, req.RollNumber, req.FirstName, req.LastName, req.Email, req.Phone, req.DateOfBirth,
		req.Address, req.SchoolID, req.EnrollmentDate, req.Status, now, now).Scan(
		&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
		&student.DateOfBirth, &student.Address, &student.SchoolID, &student.EnrollmentDate,
		&student.Status, &student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		shared.LogError("STUDENT_SERVICE", "create student", err)
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "students_roll_number_school_id_key") {
				shared.WriteErrorResponse(w, http.StatusConflict, "ROLL_NUMBER_EXISTS", "Student with this roll number already exists in this school")
			} else {
				shared.WriteErrorResponse(w, http.StatusConflict, "EMAIL_EXISTS", "Student with this email already exists")
			}
			return
		}
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create student")
		return
	}

	student.SchoolName = schoolName
	shared.WriteSuccessResponse(w, http.StatusCreated, "Student created successfully", student)
}

func (s *StudentService) GetStudents(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
		FROM students
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "get students", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get students")
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(
			&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
			&student.DateOfBirth, &student.Address, &student.SchoolID, &student.EnrollmentDate,
			&student.Status, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			shared.LogError("STUDENT_SERVICE", "scan student", err)
			continue
		}

		// Get school name via gRPC (optional, for performance you might want to do this in batch)
		_, schoolName, _ := s.validateSchool(student.SchoolID)
		student.SchoolName = schoolName

		students = append(students, student)
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Students retrieved successfully", students)
}

func (s *StudentService) GetStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid student ID")
		return
	}

	student, err := s.getStudentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteErrorResponse(w, http.StatusNotFound, "STUDENT_NOT_FOUND", "Student not found")
			return
		}
		shared.LogError("STUDENT_SERVICE", "get student", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get student")
		return
	}

	// Get school name via gRPC
	_, schoolName, _ := s.validateSchool(student.SchoolID)
	student.SchoolName = schoolName

	shared.WriteSuccessResponse(w, http.StatusOK, "Student retrieved successfully", student)
}

func (s *StudentService) GetStudentsBySchool(w http.ResponseWriter, r *http.Request) {
	schoolIDStr := strings.TrimPrefix(r.URL.Path, "/students/school/")
	schoolID, err := strconv.Atoi(schoolIDStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SCHOOL_ID", "Invalid school ID")
		return
	}

	// Validate school exists
	schoolExists, schoolName, err := s.validateSchool(schoolID)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school validation", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to validate school")
		return
	}

	if !schoolExists {
		shared.WriteErrorResponse(w, http.StatusNotFound, "SCHOOL_NOT_FOUND", "School not found")
		return
	}

	query := `
		SELECT id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
		FROM students
		WHERE school_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, schoolID)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "get students by school", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get students")
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(
			&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
			&student.DateOfBirth, &student.Address, &student.SchoolID, &student.EnrollmentDate,
			&student.Status, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			shared.LogError("STUDENT_SERVICE", "scan student", err)
			continue
		}
		student.SchoolName = schoolName
		students = append(students, student)
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Students retrieved successfully", students)
}

func (s *StudentService) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid student ID")
		return
	}

	var req UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RollNumber == "" || req.FirstName == "" || req.LastName == "" || req.SchoolID == 0 {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Roll number, first name, last name, and school ID are required")
		return
	}

	// Validate school exists
	schoolExists, schoolName, err := s.validateSchool(req.SchoolID)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school validation", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to validate school")
		return
	}

	if !schoolExists {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_SCHOOL", "School does not exist")
		return
	}

	query := `
		UPDATE students
		SET roll_number = $1, first_name = $2, last_name = $3, email = $4, phone = $5, date_of_birth = $6,
		    address = $7, school_id = $8, enrollment_date = $9, status = $10, updated_at = $11
		WHERE id = $12
		RETURNING id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
	`

	var student Student
	err = s.db.QueryRow(query, req.RollNumber, req.FirstName, req.LastName, req.Email, req.Phone, req.DateOfBirth,
		req.Address, req.SchoolID, req.EnrollmentDate, req.Status, time.Now(), id).Scan(
		&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
		&student.DateOfBirth, &student.Address, &student.SchoolID, &student.EnrollmentDate,
		&student.Status, &student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteErrorResponse(w, http.StatusNotFound, "STUDENT_NOT_FOUND", "Student not found")
			return
		}
		shared.LogError("STUDENT_SERVICE", "update student", err)
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "students_roll_number_school_id_key") {
				shared.WriteErrorResponse(w, http.StatusConflict, "ROLL_NUMBER_EXISTS", "Student with this roll number already exists in this school")
			} else {
				shared.WriteErrorResponse(w, http.StatusConflict, "EMAIL_EXISTS", "Student with this email already exists")
			}
			return
		}
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update student")
		return
	}

	student.SchoolName = schoolName
	shared.WriteSuccessResponse(w, http.StatusOK, "Student updated successfully", student)
}

func (s *StudentService) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid student ID")
		return
	}

	query := `DELETE FROM students WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "delete student", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete student")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		shared.WriteErrorResponse(w, http.StatusNotFound, "STUDENT_NOT_FOUND", "Student not found")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Student deleted successfully", nil)
}

func (s *StudentService) getStudentByID(id int) (*Student, error) {
	query := `
		SELECT id, roll_number, first_name, last_name, email, phone, date_of_birth, address, school_id, enrollment_date, status, created_at, updated_at
		FROM students
		WHERE id = $1
	`

	var student Student
	err := s.db.QueryRow(query, id).Scan(
		&student.ID, &student.RollNumber, &student.FirstName, &student.LastName, &student.Email, &student.Phone,
		&student.DateOfBirth, &student.Address, &student.SchoolID, &student.EnrollmentDate,
		&student.Status, &student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
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
		dbUser = "studentuser"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "studentpass"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "studentdb"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8083"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	schoolServiceGRPC := os.Getenv("SCHOOL_SERVICE_GRPC")
	if schoolServiceGRPC == "" {
		schoolServiceGRPC = "localhost:50051"
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

	// Connect to School Service gRPC
	var schoolServiceConn *grpc.ClientConn
	schoolServiceConn, err = grpc.Dial(schoolServiceGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school service gRPC connection", err)
		shared.LogInfo("STUDENT_SERVICE", "Continuing without school service validation")
	} else {
		defer schoolServiceConn.Close()
	}

	// Create student service
	studentService := NewStudentService(db, schoolServiceConn)

	// Start gRPC server (placeholder for future student-specific gRPC endpoints)
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		shared.LogInfo("STUDENT_SERVICE", fmt.Sprintf("Starting gRPC server on port %s", grpcPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Setup HTTP routes
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			authMiddleware(studentService.GetStudents)(w, r)
		case "POST":
			authMiddleware(studentService.CreateStudent)(w, r)
		default:
			shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		}
	})

	http.HandleFunc("/students/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a school-specific endpoint
		if strings.HasPrefix(r.URL.Path, "/students/school/") {
			if r.Method == "GET" {
				authMiddleware(studentService.GetStudentsBySchool)(w, r)
				return
			}
		} else {
			// Regular student endpoints
			switch r.Method {
			case "GET":
				authMiddleware(studentService.GetStudent)(w, r)
			case "PUT":
				authMiddleware(studentService.UpdateStudent)(w, r)
			case "DELETE":
				authMiddleware(studentService.DeleteStudent)(w, r)
			default:
				shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			}
		}
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteSuccessResponse(w, http.StatusOK, "Student service is healthy", nil)
	})

	shared.LogInfo("STUDENT_SERVICE", fmt.Sprintf("Starting HTTP server on port %s", httpPort))
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
