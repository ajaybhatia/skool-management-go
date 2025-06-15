package handlers

import (
	"encoding/json"
	"net/http"

	"skool-management/auth-service/internal/models"
	"skool-management/auth-service/internal/service"
	"skool-management/shared"
)

type AuthHandlers struct {
	authService *service.AuthService
}

func NewAuthHandlers(authService *service.AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

func (h *AuthHandlers) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	user, err := h.authService.Signup(&req)
	if err != nil {
		switch err.Error() {
		case "all fields are required":
			shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case "user with this email already exists":
			shared.WriteErrorResponse(w, http.StatusConflict, "USER_EXISTS", err.Error())
		default:
			shared.LogError("AUTH_SERVICE", "signup", err)
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create user")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusCreated, "User created successfully", user)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
		} else {
			shared.LogError("AUTH_SERVICE", "login", err)
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Login failed")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Login successful", response)
}

func (h *AuthHandlers) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	accessToken, err := h.authService.RefreshToken(&req)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
		return
	}

	response := map[string]string{
		"access_token": accessToken,
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Token refreshed successfully", response)
}

func (h *AuthHandlers) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
		return
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Token is valid", map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
	})
}

func (h *AuthHandlers) Health(w http.ResponseWriter, r *http.Request) {
	shared.WriteSuccessResponse(w, http.StatusOK, "Auth service is healthy", nil)
}
