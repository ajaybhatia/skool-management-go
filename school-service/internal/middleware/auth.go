package middleware

import (
	"net/http"
	"time"

	"skool-management/shared"
)

func AuthMiddleware(jwtSecret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				shared.WriteErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
				return
			}

			tokenString := authHeader
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}

			// Create JWT manager with same secret as auth service
			jwtManager := shared.NewJWTManager(jwtSecret, "", time.Hour, time.Hour)
			_, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
				return
			}

			next(w, r)
		}
	}
}
