package shared

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	SecretKey        string
	RefreshSecretKey string
	TokenDuration    time.Duration
	RefreshDuration  time.Duration
}

func NewJWTManager(secretKey, refreshSecretKey string, tokenDuration, refreshDuration time.Duration) *JWTManager {
	return &JWTManager{
		SecretKey:        secretKey,
		RefreshSecretKey: refreshSecretKey,
		TokenDuration:    tokenDuration,
		RefreshDuration:  refreshDuration,
	}
}

func (manager *JWTManager) GenerateToken(userID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.SecretKey))
}

func (manager *JWTManager) GenerateRefreshToken(userID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.RefreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.RefreshSecretKey))
}

func (manager *JWTManager) VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(manager.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (manager *JWTManager) VerifyRefreshToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(manager.RefreshSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

// ValidateJWT validates a JWT token using the default secret
func ValidateJWT(tokenString string) bool {
	jwtSecret := GetEnv("JWT_SECRET", "your-secret-key")
	manager := NewJWTManager(jwtSecret, jwtSecret, time.Hour, time.Hour*24*7)

	_, err := manager.VerifyToken(tokenString)
	return err == nil
}
