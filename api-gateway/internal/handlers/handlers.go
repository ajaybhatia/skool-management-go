package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"skool-management/api-gateway/internal/gateway"
	"skool-management/shared"
)

type Handlers struct {
	gateway *gateway.Gateway
}

func New(gw *gateway.Gateway) *Handlers {
	return &Handlers{
		gateway: gw,
	}
}

// Auth service routes
func (h *Handlers) HandleAuth(w http.ResponseWriter, r *http.Request) {
	// Remove /auth prefix from path
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/auth")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}
	h.gateway.ProxyRequest(h.gateway.GetAuthServiceURL(), w, r)
}

// School service routes
func (h *Handlers) HandleSchools(w http.ResponseWriter, r *http.Request) {
	// Remove /schools prefix and proxy to school service
	h.gateway.ProxyRequest(h.gateway.GetSchoolServiceURL(), w, r)
}

// Student service routes
func (h *Handlers) HandleStudents(w http.ResponseWriter, r *http.Request) {
	// Remove /students prefix and proxy to student service
	h.gateway.ProxyRequest(h.gateway.GetStudentServiceURL(), w, r)
}

// Health check endpoint
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	// Check health of all services
	services := map[string]string{
		"auth":    h.gateway.GetAuthServiceURL() + "/health",
		"school":  h.gateway.GetSchoolServiceURL() + "/health",
		"student": h.gateway.GetStudentServiceURL() + "/health",
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
func (h *Handlers) HandleDocs(w http.ResponseWriter, r *http.Request) {
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
