package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"skool-management/student-service/internal/config"
	"skool-management/student-service/internal/handlers"
	"skool-management/student-service/internal/middleware"
	"skool-management/student-service/internal/repository"
	"skool-management/student-service/internal/service"
	"skool-management/shared"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", cfg.GetDSN())
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
	schoolServiceConn, err = grpc.Dial(cfg.SchoolServiceGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		shared.LogError("STUDENT_SERVICE", "school service gRPC connection", err)
		shared.LogInfo("STUDENT_SERVICE", "Continuing without school service validation")
	} else {
		defer schoolServiceConn.Close()
	}

	// Initialize layers
	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo, schoolServiceConn)
	studentHandlers := handlers.NewStudentHandlers(studentService)
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Start gRPC server (placeholder for future student-specific gRPC endpoints)
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		shared.LogInfo("STUDENT_SERVICE", fmt.Sprintf("Starting gRPC server on port %s", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Setup HTTP routes
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			authMiddleware(studentHandlers.GetStudents)(w, r)
		case "POST":
			authMiddleware(studentHandlers.CreateStudent)(w, r)
		default:
			shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		}
	})

	http.HandleFunc("/students/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a school-specific endpoint
		if strings.HasPrefix(r.URL.Path, "/students/school/") {
			if r.Method == "GET" {
				authMiddleware(studentHandlers.GetStudentsBySchool)(w, r)
				return
			}
		} else {
			// Regular student endpoints
			switch r.Method {
			case "GET":
				authMiddleware(studentHandlers.GetStudent)(w, r)
			case "PUT":
				authMiddleware(studentHandlers.UpdateStudent)(w, r)
			case "DELETE":
				authMiddleware(studentHandlers.DeleteStudent)(w, r)
			default:
				shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
			}
		}
	})

	// Health check
	http.HandleFunc("/health", studentHandlers.Health)

	shared.LogInfo("STUDENT_SERVICE", fmt.Sprintf("Starting HTTP server on port %s", cfg.HTTPPort))
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, nil))
}
