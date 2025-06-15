package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"skool-management/shared"
)

type Middleware struct {
	authServiceURL string
}

func New(authServiceURL string) *Middleware {
	return &Middleware{
		authServiceURL: authServiceURL,
	}
}

// CORS middleware
func (m *Middleware) CORS(next http.HandlerFunc) http.HandlerFunc {
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

// Logging middleware
func (m *Middleware) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		shared.LogInfo("API_GATEWAY", fmt.Sprintf("%s %s - Started", r.Method, r.URL.Path))

		next(w, r)

		shared.LogInfo("API_GATEWAY", fmt.Sprintf("%s %s - Completed in %v", r.Method, r.URL.Path, time.Since(start)))
	}
}

// Rate limiting middleware (basic implementation)
func (m *Middleware) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simple rate limiting - in production, use Redis or more sophisticated solution
		// For now, just pass through
		next(w, r)
	}
}

// JWT authentication middleware
func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			shared.WriteErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			return
		}

		// Validate token with auth service
		validateReq, err := http.NewRequest("GET", m.authServiceURL+"/validate", nil)
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
