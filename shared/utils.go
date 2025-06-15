package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, error, message string) {
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}
	WriteJSONResponse(w, statusCode, response)
}

func WriteSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Message: message,
		Data:    data,
	}
	WriteJSONResponse(w, statusCode, response)
}

func LogError(service, operation string, err error) {
	fmt.Printf("[%s] Error in %s: %v\n", service, operation, err)
}

func LogInfo(service, message string) {
	fmt.Printf("[%s] %s\n", service, message)
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
