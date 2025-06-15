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
}

func New(authURL, schoolURL, studentURL string) *Gateway {
	return &Gateway{
		authServiceURL:    authURL,
		schoolServiceURL:  schoolURL,
		studentServiceURL: studentURL,
	}
}

// ProxyRequest proxies requests to target services
func (g *Gateway) ProxyRequest(targetURL string, w http.ResponseWriter, r *http.Request) {
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

func (g *Gateway) GetAuthServiceURL() string {
	return g.authServiceURL
}

func (g *Gateway) GetSchoolServiceURL() string {
	return g.schoolServiceURL
}

func (g *Gateway) GetStudentServiceURL() string {
	return g.studentServiceURL
}
