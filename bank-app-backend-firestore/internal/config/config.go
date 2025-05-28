package config

import (
	"os"
	"strconv"
)

// Config - Application configuration
type Config struct {
	Port              int
	FirebaseProjectID string
	FirestoreEmulator string
	AuthEmulator      string
	JWTSecret         string
	UserID            string
}

// New - Create a new configuration
func New() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))

	return &Config{
		Port:              port,
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", "seventh-league-405315"),
		FirestoreEmulator: getEnv("FIRESTORE_EMULATOR_HOST", "localhost:8091"),
		AuthEmulator:      getEnv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099"),
		JWTSecret:         getEnv("JWT_SECRET", "your-very-secret-jwt-key-change-in-production"),
		UserID:            getEnv("UNIQUE_USER_ID", "demo_user"),
	}
}

// getEnv - Get environment variable or default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
