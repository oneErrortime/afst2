package config

import (
	"os"
	"strconv"
)

var (
	// Database configuration
	DatabaseURL string

	// JWT configuration
	JWTSecretKey          string
	JWTAlgorithm          string
	JWTAccessTokenExpiry int64
)

func LoadConfig() {
	// Database configuration
	DatabaseURL = getEnv("DATABASE_URL", "library.db")

	// JWT configuration
	JWTSecretKey = getEnv("SECRET_KEY", "your-secret-key-change-in-production")
	JWTAlgorithm = getEnv("ALGORITHM", "HS256")
	accessTokenExpiry, err := strconv.Atoi(getEnv("ACCESS_TOKEN_EXPIRE_MINUTES", "30"))
	if err != nil {
		accessTokenExpiry = 30 // default to 30 minutes
	}
	JWTAccessTokenExpiry = int64(accessTokenExpiry)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}