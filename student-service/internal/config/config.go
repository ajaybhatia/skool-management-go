package config

import "os"

type Config struct {
	HTTPPort          string
	GRPCPort          string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	SchoolServiceGRPC string
}

func Load() *Config {
	return &Config{
		HTTPPort:          getEnv("HTTP_PORT", "8083"),
		GRPCPort:          getEnv("GRPC_PORT", "50052"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "studentuser"),
		DBPassword:        getEnv("DB_PASSWORD", "studentpass"),
		DBName:            getEnv("DB_NAME", "studentdb"),
		JWTSecret:         getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		SchoolServiceGRPC: getEnv("SCHOOL_SERVICE_GRPC", "localhost:50051"),
	}
}

func (c *Config) GetDSN() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser +
		" password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
