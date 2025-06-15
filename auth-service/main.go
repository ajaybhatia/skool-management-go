package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"skool-management/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"-"`
	FirstName    string             `bson:"first_name" json:"first_name"`
	LastName     string             `bson:"last_name" json:"last_name"`
	Role         string             `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	RefreshToken string             `bson:"refresh_token" json:"-"`
}

type SignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthService struct {
	db         *mongo.Database
	jwtManager *shared.JWTManager
}

func NewAuthService(db *mongo.Database, jwtManager *shared.JWTManager) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "All fields are required")
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Check if user already exists
	collection := s.db.Collection("users")
	var existingUser User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		shared.WriteErrorResponse(w, http.StatusConflict, "USER_EXISTS", "User with this email already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		shared.LogError("AUTH_SERVICE", "password hashing", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process password")
		return
	}

	// Create user
	user := User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		shared.LogError("AUTH_SERVICE", "user creation", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create user")
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	shared.WriteSuccessResponse(w, http.StatusCreated, "User created successfully", user)
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Find user by email
	collection := s.db.Collection("users")
	var user User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		shared.LogError("AUTH_SERVICE", "access token generation", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate access token")
		return
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.Hex(), user.Email)
	if err != nil {
		shared.LogError("AUTH_SERVICE", "refresh token generation", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate refresh token")
		return
	}

	// Store refresh token in database
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{"refresh_token": refreshToken, "updated_at": time.Now()}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		shared.LogError("AUTH_SERVICE", "refresh token storage", err)
	}

	response := LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Login successful", response)
}

func (s *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Verify refresh token
	claims, err := s.jwtManager.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid refresh token")
		return
	}

	// Find user and verify refresh token in database
	collection := s.db.Collection("users")
	userID, _ := primitive.ObjectIDFromHex(claims.UserID)
	var user User
	err = collection.FindOne(context.Background(), bson.M{"_id": userID, "refresh_token": req.RefreshToken}).Decode(&user)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid refresh token")
		return
	}

	// Generate new access token
	accessToken, err := s.jwtManager.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		shared.LogError("AUTH_SERVICE", "access token generation", err)
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate access token")
		return
	}

	response := map[string]string{
		"access_token": accessToken,
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Token refreshed successfully", response)
}

func (s *AuthService) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
		return
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := s.jwtManager.VerifyToken(tokenString)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Token is valid", map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
	})
}

func main() {
	// Get environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@localhost:27017/authdb?authSource=admin"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-in-production"
	}

	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtRefreshSecret == "" {
		jwtRefreshSecret = "your-super-secret-refresh-key-change-in-production"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	// Test connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	db := client.Database("authdb")

	// Create JWT manager
	jwtManager := shared.NewJWTManager(
		jwtSecret,
		jwtRefreshSecret,
		15*time.Minute, // Access token duration
		7*24*time.Hour, // Refresh token duration (7 days)
	)

	// Create auth service
	authService := NewAuthService(db, jwtManager)

	// Setup routes
	http.HandleFunc("/signup", authService.Signup)
	http.HandleFunc("/login", authService.Login)
	http.HandleFunc("/refresh", authService.Refresh)
	http.HandleFunc("/validate", authService.ValidateToken)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteSuccessResponse(w, http.StatusOK, "Auth service is healthy", nil)
	})

	shared.LogInfo("AUTH_SERVICE", fmt.Sprintf("Starting auth service on port %s", port))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
