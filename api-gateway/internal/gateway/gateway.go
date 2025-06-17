package gateway

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"skool-management/shared"
)

type Gateway struct {
	authServiceURL    string
	schoolServiceURL  string
	studentServiceURL string
	// Circuit breakers for each service
	authCircuitBreaker    *shared.CircuitBreaker
	schoolCircuitBreaker  *shared.CircuitBreaker
	studentCircuitBreaker *shared.CircuitBreaker
}

func New(authURL, schoolURL, studentURL string) *Gateway {
	return &Gateway{
		authServiceURL:    authURL,
		schoolServiceURL:  schoolURL,
		studentServiceURL: studentURL,
		// Initialize circuit breakers for each service
		authCircuitBreaker: shared.NewCircuitBreaker(shared.CircuitBreakerConfig{
			Name:         "auth-service",
			MaxFailures:  5,
			ResetTimeout: 30 * time.Second,
		}),
		schoolCircuitBreaker: shared.NewCircuitBreaker(shared.CircuitBreakerConfig{
			Name:         "school-service",
			MaxFailures:  5,
			ResetTimeout: 30 * time.Second,
		}),
		studentCircuitBreaker: shared.NewCircuitBreaker(shared.CircuitBreakerConfig{
			Name:         "student-service",
			MaxFailures:  5,
			ResetTimeout: 30 * time.Second,
		}),
	}
}

// ProxyRequest proxies requests to target services with circuit breaker protection
func (g *Gateway) ProxyRequest(targetURL string, w http.ResponseWriter, r *http.Request) {
	// Determine which circuit breaker to use based on target URL
	var circuitBreaker *shared.CircuitBreaker
	switch {
	case targetURL == g.authServiceURL:
		circuitBreaker = g.authCircuitBreaker
	case targetURL == g.schoolServiceURL:
		circuitBreaker = g.schoolCircuitBreaker
	case targetURL == g.studentServiceURL:
		circuitBreaker = g.studentCircuitBreaker
	default:
		// Fallback for unknown services
		circuitBreaker = shared.NewCircuitBreaker(shared.CircuitBreakerConfig{
			Name:         "unknown-service",
			MaxFailures:  3,
			ResetTimeout: 30 * time.Second,
		})
	}

	// Execute request with circuit breaker protection
	err := circuitBreaker.Execute(func() error {
		return g.makeProxyRequest(targetURL, w, r)
	})

	if err != nil {
		if err.Error() == "circuit breaker is OPEN" {
			shared.WriteErrorResponse(w, http.StatusServiceUnavailable, "CIRCUIT_BREAKER_OPEN",
				"Service is temporarily unavailable due to circuit breaker")
			return
		}
		// Error already handled in makeProxyRequest
		return
	}
}

// makeProxyRequest performs the actual HTTP request
func (g *Gateway) makeProxyRequest(targetURL string, w http.ResponseWriter, r *http.Request) error {
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
		return err
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
		return err
	}
	defer resp.Body.Close()

	// Check for server errors to trigger circuit breaker
	if resp.StatusCode >= 500 {
		shared.WriteErrorResponse(w, resp.StatusCode, "SERVICE_ERROR", "Target service returned an error")
		return err
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(w, resp.Body)
	return err
}

// Getter methods for service URLs
func (g *Gateway) GetAuthServiceURL() string {
	return g.authServiceURL
}

func (g *Gateway) GetSchoolServiceURL() string {
	return g.schoolServiceURL
}

func (g *Gateway) GetStudentServiceURL() string {
	return g.studentServiceURL
}

// Getter methods for circuit breakers
func (g *Gateway) GetAuthCircuitBreaker() *shared.CircuitBreaker {
	return g.authCircuitBreaker
}

func (g *Gateway) GetSchoolCircuitBreaker() *shared.CircuitBreaker {
	return g.schoolCircuitBreaker
}

func (g *Gateway) GetStudentCircuitBreaker() *shared.CircuitBreaker {
	return g.studentCircuitBreaker
}
