package config

import (
	"os"
	"time"
)

type Config struct {
	Port             string
	MongoURI         string
	JWTSecret        string
	JWTRefreshSecret string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8081"),
		MongoURI:         getEnv("MONGODB_URI", "mongodb://admin:password@localhost:27017/authdb?authSource=admin"),
		JWTSecret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-super-secret-refresh-key-change-in-production"),
		AccessTokenTTL:   15 * time.Minute,
		RefreshTokenTTL:  7 * 24 * time.Hour, // 7 days
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
