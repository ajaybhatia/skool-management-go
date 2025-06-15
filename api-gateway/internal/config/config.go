package config

import "os"

type Config struct {
	Port              string
	AuthServiceURL    string
	SchoolServiceURL  string
	StudentServiceURL string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		SchoolServiceURL: getEnv("SCHOOL_SERVICE_URL", "http://localhost:8082"),
		StudentServiceURL: getEnv("STUDENT_SERVICE_URL", "http://localhost:8083"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
