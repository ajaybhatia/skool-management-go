package main

import (
	"fmt"
	"log"
	"net/http"

	"skool-management/api-gateway/internal/config"
	"skool-management/api-gateway/internal/gateway"
	"skool-management/api-gateway/internal/handlers"
	"skool-management/api-gateway/internal/middleware"
	"skool-management/shared"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create gateway
	gw := gateway.New(cfg.AuthServiceURL, cfg.SchoolServiceURL, cfg.StudentServiceURL)

	// Create middleware
	mw := middleware.New(cfg.AuthServiceURL)

	// Create handlers
	h := handlers.New(gw)

	// Setup routes with middleware chain
	mux := http.NewServeMux()

	// Health and documentation endpoints
	mux.HandleFunc("/health", mw.CORS(mw.Logging(h.HandleHealth)))
	mux.HandleFunc("/docs", mw.CORS(mw.Logging(h.HandleDocs)))

	// Service routes
	mux.HandleFunc("/auth/", mw.CORS(mw.Logging(mw.RateLimit(h.HandleAuth))))
	mux.HandleFunc("/schools", mw.CORS(mw.Logging(mw.Auth(mw.RateLimit(h.HandleSchools)))))
	mux.HandleFunc("/schools/", mw.CORS(mw.Logging(mw.Auth(mw.RateLimit(h.HandleSchools)))))
	mux.HandleFunc("/students", mw.CORS(mw.Logging(mw.Auth(mw.RateLimit(h.HandleStudents)))))
	mux.HandleFunc("/students/", mw.CORS(mw.Logging(mw.Auth(mw.RateLimit(h.HandleStudents)))))

	// Root endpoint
	mux.HandleFunc("/", mw.CORS(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			h.HandleDocs(w, r)
		} else {
			shared.WriteErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
		}
	}))

	shared.LogInfo("API_GATEWAY", fmt.Sprintf("Starting API Gateway on port %s", cfg.Port))
	shared.LogInfo("API_GATEWAY", fmt.Sprintf("Auth Service: %s", cfg.AuthServiceURL))
	shared.LogInfo("API_GATEWAY", fmt.Sprintf("School Service: %s", cfg.SchoolServiceURL))
	shared.LogInfo("API_GATEWAY", fmt.Sprintf("Student Service: %s", cfg.StudentServiceURL))
	shared.LogInfo("API_GATEWAY", "Available endpoints:")
	shared.LogInfo("API_GATEWAY", "  GET  / - API Documentation")
	shared.LogInfo("API_GATEWAY", "  GET  /docs - API Documentation")
	shared.LogInfo("API_GATEWAY", "  GET  /health - Health Check")
	shared.LogInfo("API_GATEWAY", "  POST /auth/signup - User Registration")
	shared.LogInfo("API_GATEWAY", "  POST /auth/login - User Login")
	shared.LogInfo("API_GATEWAY", "  POST /auth/refresh - Refresh Token")
	shared.LogInfo("API_GATEWAY", "  *    /schools/* - School Management")
	shared.LogInfo("API_GATEWAY", "  *    /students/* - Student Management")

	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
