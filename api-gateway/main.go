package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"skool-management/shared"
)

type APIGateway struct {
	authServiceURL    string
	schoolServiceURL  string
	studentServiceURL string
}

type ProxyResponse struct {
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
}

func NewAPIGateway(authURL, schoolURL, studentURL string) *APIGateway {
	return &APIGateway{
		authServiceURL:    authURL,
		schoolServiceURL:  schoolURL,
		studentServiceURL: studentURL,
	}
}

// Middleware for CORS
func (gw *APIGateway) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Middleware for logging
func (gw *APIGateway) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		shared.LogInfo("API_GATEWAY", fmt.Sprintf("%s %s - Started", r.Method, r.URL.Path))

		next(w, r)

		shared.LogInfo("API_GATEWAY", fmt.Sprintf("%s %s - Completed in %v", r.Method, r.URL.Path, time.Since(start)))
	}
}

// Middleware for rate limiting (basic implementation)
func (gw *APIGateway) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simple rate limiting - in production, use Redis or more sophisticated solution
		// For now, just pass through
		next(w, r)
	}
}

// Middleware for JWT authentication
func (gw *APIGateway) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			shared.WriteErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			return
		}

		// Validate token with auth service
		validateReq, err := http.NewRequest("GET", gw.authServiceURL+"/validate", nil)
		if err != nil {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to create validation request")
			return
		}
		validateReq.Header.Set("Authorization", authHeader)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(validateReq)
		if err != nil {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to validate token")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Token is invalid
			body, _ := io.ReadAll(resp.Body)
			var errorResp map[string]interface{}
			json.Unmarshal(body, &errorResp)

			if errorResp["error"] != nil {
				shared.WriteErrorResponse(w, http.StatusUnauthorized, errorResp["error"].(string), errorResp["message"].(string))
			} else {
				shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			}
			return
		}

		// Token is valid, proceed to next handler
		next(w, r)
	}
}

// Generic proxy function
func (gw *APIGateway) proxyRequest(targetURL string, w http.ResponseWriter, r *http.Request) {
	// Read request body
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body.Close()
	}

	// Create new request
	fullURL := targetURL + r.URL.Path
	if r.URL.RawQuery != "" {
		fullURL += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, fullURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "PROXY_ERROR", "Failed to create proxy request")
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		shared.LogError("API_GATEWAY", "proxy request", err)
		shared.WriteErrorResponse(w, http.StatusBadGateway, "SERVICE_UNAVAILABLE", "Target service is unavailable")
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}

// Auth service routes
func (gw *APIGateway) handleAuth(w http.ResponseWriter, r *http.Request) {
	// Remove /auth prefix from path
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/auth")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}
	gw.proxyRequest(gw.authServiceURL, w, r)
}

// School service routes
func (gw *APIGateway) handleSchools(w http.ResponseWriter, r *http.Request) {
	// Remove /schools prefix and proxy to school service
	gw.proxyRequest(gw.schoolServiceURL, w, r)
}

// Student service routes
func (gw *APIGateway) handleStudents(w http.ResponseWriter, r *http.Request) {
	// Remove /students prefix and proxy to student service
	gw.proxyRequest(gw.studentServiceURL, w, r)
}

// Health check endpoint
func (gw *APIGateway) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Check health of all services
	services := map[string]string{
		"auth":    gw.authServiceURL + "/health",
		"school":  gw.schoolServiceURL + "/health",
		"student": gw.studentServiceURL + "/health",
	}

	healthStatus := make(map[string]interface{})
	allHealthy := true

	client := &http.Client{Timeout: 5 * time.Second}

	for service, url := range services {
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			healthStatus[service] = map[string]interface{}{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			allHealthy = false
		} else {
			healthStatus[service] = map[string]interface{}{
				"status": "healthy",
			}
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	healthStatus["gateway"] = map[string]interface{}{
		"status": "healthy",
	}

	if allHealthy {
		shared.WriteSuccessResponse(w, http.StatusOK, "All services are healthy", healthStatus)
	} else {
		shared.WriteErrorResponse(w, http.StatusServiceUnavailable, "PARTIAL_OUTAGE", "Some services are unhealthy")
		// Still include the health status in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "PARTIAL_OUTAGE",
			"message": "Some services are unhealthy",
			"data":    healthStatus,
		})
	}
}

// API documentation endpoint
func (gw *APIGateway) handleDocs(w http.ResponseWriter, r *http.Request) {
	docs := map[string]interface{}{
		"title":       "School Management API Gateway",
		"version":     "1.0.0",
		"description": "API Gateway for School Management Microservices",
		"endpoints": map[string]interface{}{
			"auth": map[string]interface{}{
				"signup": map[string]string{
					"method":      "POST",
					"path":        "/auth/signup",
					"description": "Register a new user",
				},
				"login": map[string]string{
					"method":      "POST",
					"path":        "/auth/login",
					"description": "Login user and get JWT tokens",
				},
				"refresh": map[string]string{
					"method":      "POST",
					"path":        "/auth/refresh",
					"description": "Refresh JWT access token",
				},
				"validate": map[string]string{
					"method":      "POST",
					"path":        "/auth/validate",
					"description": "Validate JWT token",
				},
			},
			"schools": map[string]interface{}{
				"list": map[string]string{
					"method":      "GET",
					"path":        "/schools",
					"description": "Get all schools",
					"auth":        "required",
				},
				"create": map[string]string{
					"method":      "POST",
					"path":        "/schools",
					"description": "Create new school",
					"auth":        "required",
				},
				"get": map[string]string{
					"method":      "GET",
					"path":        "/schools/{id}",
					"description": "Get school by ID",
					"auth":        "required",
				},
				"update": map[string]string{
					"method":      "PUT",
					"path":        "/schools/{id}",
					"description": "Update school",
					"auth":        "required",
				},
				"delete": map[string]string{
					"method":      "DELETE",
					"path":        "/schools/{id}",
					"description": "Delete school",
					"auth":        "required",
				},
			},
			"students": map[string]interface{}{
				"list": map[string]string{
					"method":      "GET",
					"path":        "/students",
					"description": "Get all students",
					"auth":        "required",
				},
				"create": map[string]string{
					"method":      "POST",
					"path":        "/students",
					"description": "Create new student",
					"auth":        "required",
				},
				"get": map[string]string{
					"method":      "GET",
					"path":        "/students/{id}",
					"description": "Get student by ID",
					"auth":        "required",
				},
				"update": map[string]string{
					"method":      "PUT",
					"path":        "/students/{id}",
					"description": "Update student",
					"auth":        "required",
				},
				"delete": map[string]string{
					"method":      "DELETE",
					"path":        "/students/{id}",
					"description": "Delete student",
					"auth":        "required",
				},
				"by_school": map[string]string{
					"method":      "GET",
					"path":        "/students/school/{school_id}",
					"description": "Get students by school ID",
					"auth":        "required",
				},
			},
		},
		"authentication": map[string]interface{}{
			"type":        "Bearer Token (JWT)",
			"header":      "Authorization: Bearer <token>",
			"description": "Include JWT token in Authorization header for protected endpoints",
		},
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "API Documentation", docs)
}

func main() {
	// Get environment variables
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8081"
	}

	schoolServiceURL := os.Getenv("SCHOOL_SERVICE_URL")
	if schoolServiceURL == "" {
		schoolServiceURL = "http://localhost:8082"
	}

	studentServiceURL := os.Getenv("STUDENT_SERVICE_URL")
	if studentServiceURL == "" {
		studentServiceURL = "http://localhost:8083"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create API Gateway
	gateway := NewAPIGateway(authServiceURL, schoolServiceURL, studentServiceURL)

	// Setup routes with middleware chain
	mux := http.NewServeMux()

	// Health and documentation endpoints
	mux.HandleFunc("/health", gateway.corsMiddleware(gateway.loggingMiddleware(gateway.handleHealth)))
	mux.HandleFunc("/docs", gateway.corsMiddleware(gateway.loggingMiddleware(gateway.handleDocs)))

	// Service routes
	mux.HandleFunc("/auth/", gateway.corsMiddleware(gateway.loggingMiddleware(gateway.rateLimitMiddleware(gateway.handleAuth))))
	mux.HandleFunc("/schools", gateway.corsMiddleware(gateway.loggingMiddleware(gateway.authMiddleware(gateway.rateLimitMiddleware(gateway.handleSchools)))))
	mux.HandleFunc("/schools/", gateway.corsMiddleware(gateway.loggingMiddleware(gateway.authMiddleware(gateway.rateLimitMiddleware(gateway.handleSchools)))))
	mux.HandleFunc("/students", gateway.corsMiddleware(gateway.loggingMiddleware(gateway.authMiddleware(gateway.rateLimitMiddleware(gateway.handleStudents)))))
	mux.HandleFunc("/students/", gateway.corsMiddleware(gateway.loggingMiddleware(gateway.authMiddleware(gateway.rateLimitMiddleware(gateway.handleStudents)))))

	// Root endpoint
	mux.HandleFunc("/", gateway.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			gateway.handleDocs(w, r)
		} else {
			shared.WriteErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
		}
	}))

	shared.LogInfo("API_GATEWAY", fmt.Sprintf("Starting API Gateway on port %s", port))
	shared.LogInfo("API_GATEWAY", fmt.Sprintf("Auth Service: %s", authServiceURL))
	shared.LogInfo("API_GATEWAY", fmt.Sprintf("School Service: %s", schoolServiceURL))
	shared.LogInfo("API_GATEWAY", fmt.Sprintf("Student Service: %s", studentServiceURL))
	shared.LogInfo("API_GATEWAY", "Available endpoints:")
	shared.LogInfo("API_GATEWAY", "  GET  / - API Documentation")
	shared.LogInfo("API_GATEWAY", "  GET  /docs - API Documentation")
	shared.LogInfo("API_GATEWAY", "  GET  /health - Health Check")
	shared.LogInfo("API_GATEWAY", "  POST /auth/signup - User Registration")
	shared.LogInfo("API_GATEWAY", "  POST /auth/login - User Login")
	shared.LogInfo("API_GATEWAY", "  POST /auth/refresh - Refresh Token")
	shared.LogInfo("API_GATEWAY", "  *    /schools/* - School Management")
	shared.LogInfo("API_GATEWAY", "  *    /students/* - Student Management")

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
