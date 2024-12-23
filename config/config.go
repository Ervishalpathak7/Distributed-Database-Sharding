package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
	"strconv"
)

// DatabaseConfig holds the database connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// ServerConfig holds the server settings
type ServerConfig struct {
	Host string
	Port int
}

// Config holds the database and server settings
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// LoadEnv loads environment variables from a .env file (if present) and returns a Config struct
func LoadEnv() (*Config, error) {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err

	}
	return &Config{
		Database: getDatabaseConfig(),
		Server:   getServerConfig(),
	} , nil
}


// getDatabaseConfig retrieves the database configuration from environment variables
func getDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		Database: getEnv("DB_NAME", "postgres"),
	}
}

// getEnv retrieves an environment variable or defaults to the provided fallback value
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getServerConfig retrieves the server configuration from environment variables
func getServerConfig() ServerConfig {
	return ServerConfig{
		Host: getEnv("HOST", "localhost"),
		Port: getEnvAsInt("PORT", 8080),
	}
}

// getEnvAsInt retrieves an environment variable as an integer or defaults to the provided fallback value
func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		// Convert string to integer and return
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Unable to convert %s to int, using fallback value %d", key, fallback)
	}
	return fallback
}
