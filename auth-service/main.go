package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"skool-management/auth-service/internal/config"
	"skool-management/auth-service/internal/handlers"
	"skool-management/auth-service/internal/repository"
	"skool-management/auth-service/internal/service"
	"skool-management/shared"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
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
		cfg.JWTSecret,
		cfg.JWTRefreshSecret,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager)
	authHandlers := handlers.NewAuthHandlers(authService)

	// Setup routes
	http.HandleFunc("/signup", authHandlers.Signup)
	http.HandleFunc("/login", authHandlers.Login)
	http.HandleFunc("/refresh", authHandlers.Refresh)
	http.HandleFunc("/validate", authHandlers.ValidateToken)
	http.HandleFunc("/health", authHandlers.Health)

	shared.LogInfo("AUTH_SERVICE", fmt.Sprintf("Starting auth service on port %s", cfg.Port))
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
