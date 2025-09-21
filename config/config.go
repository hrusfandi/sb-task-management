package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "task_management"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		Port:       getEnv("PORT", "8080"),
	}

	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDatabaseURL() string {
	return "host=" + AppConfig.DBHost + " user=" + AppConfig.DBUser + " password=" + AppConfig.DBPassword + " dbname=" + AppConfig.DBName + " port=" + AppConfig.DBPort + " sslmode=disable"
}