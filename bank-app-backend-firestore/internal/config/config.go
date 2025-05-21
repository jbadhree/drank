package config

import (
	"log"
	"os"
)

type Config struct {
	DBType string // "postgres" or "firestore"
	// Postgres config
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	// Firestore config
	FirestoreProjectID string
	FirestoreCredsPath string
	// Common
	Port      string
	JWTSecret string
}

func New() *Config {

	return &Config{
		DBType:             getEnv("DB_TYPE", "firestore"),
		FirestoreProjectID: getEnv("FIRESTORE_PROJECT_ID", ""),
		FirestoreCredsPath: getEnv("FIRESTORE_CREDS_PATH", ""),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		Port:               getEnv("PORT", "8080"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	log.Println("key:", key, "value:", value)
	if value == "" {
		return defaultValue
	}
	return value
}
