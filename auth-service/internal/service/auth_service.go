package service

import (
	"errors"

	"skool-management/auth-service/internal/models"
	"skool-management/auth-service/internal/repository"
	"skool-management/shared"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *shared.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *shared.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Signup(req *models.SignupRequest) (*models.User, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, errors.New("all fields are required")
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Check if user already exists
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid email or password")
		}
		return nil, errors.New("failed to find user")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Store refresh token in database
	err = s.userRepo.UpdateRefreshToken(user.ID, refreshToken)
	if err != nil {
		shared.LogError("AUTH_SERVICE", "refresh token storage", err)
	}

	response := &models.LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *AuthService) RefreshToken(req *models.RefreshRequest) (string, error) {
	// Verify refresh token
	claims, err := s.jwtManager.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Find user and verify refresh token in database
	userID, _ := primitive.ObjectIDFromHex(claims.UserID)
	user, err := s.userRepo.GetByRefreshToken(req.RefreshToken)
	if err != nil || user.ID != userID {
		return "", errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := s.jwtManager.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}

func (s *AuthService) ValidateToken(token string) (*shared.JWTClaims, error) {
	return s.jwtManager.VerifyToken(token)
}
