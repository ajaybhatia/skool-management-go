package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"skool-management/school-service/internal/config"
	schoolGrpc "skool-management/school-service/internal/grpc"
	"skool-management/school-service/internal/handlers"
	"skool-management/school-service/internal/middleware"
	"skool-management/school-service/internal/repository"
	"skool-management/school-service/internal/service"
	"skool-management/shared"

	"google.golang.org/grpc"
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

	// Initialize layers
	schoolRepo := repository.NewSchoolRepository(db)
	schoolService := service.NewSchoolService(schoolRepo)
	schoolHandlers := handlers.NewSchoolHandlers(schoolService)
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}

		grpcServer := grpc.NewServer()
		grpcSchoolServer := schoolGrpc.NewGRPCSchoolServer(schoolService)
		RegisterSchoolServiceServer(grpcServer, grpcSchoolServer)
		reflection.Register(grpcServer)

		shared.LogInfo("SCHOOL_SERVICE", fmt.Sprintf("Starting gRPC server on port %s", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Setup HTTP routes
	http.HandleFunc("/schools", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			authMiddleware(schoolHandlers.GetSchools)(w, r)
		case "POST":
			authMiddleware(schoolHandlers.CreateSchool)(w, r)
		default:
			shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		}
	})

	http.HandleFunc("/schools/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			authMiddleware(schoolHandlers.GetSchool)(w, r)
		case "PUT":
			authMiddleware(schoolHandlers.UpdateSchool)(w, r)
		case "DELETE":
			authMiddleware(schoolHandlers.DeleteSchool)(w, r)
		default:
			shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		}
	})

	// Health check
	http.HandleFunc("/health", schoolHandlers.Health)

	shared.LogInfo("SCHOOL_SERVICE", fmt.Sprintf("Starting HTTP server on port %s", cfg.HTTPPort))
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, nil))
}
