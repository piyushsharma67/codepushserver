package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server configuration
	Port     int    `json:"port"`
	Host     string `json:"host"`
	DBType   string `json:"db_type"`
	DBHost   string `json:"db_host"`
	DBPort   int    `json:"db_port"`
	DBUser   string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName   string `json:"db_name"`
	JWTKey   string `json:"jwt_key"`
	JWTExpiry int    `json:"jwt_expiry"`
}

func NewConfig() *Config {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	return &Config{
		Port:     8080,
		Host:     "localhost",
		DBType:   getEnv("DB_TYPE", "postgres"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   port,
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "codepush"),
		JWTKey:   "your-secret-key",
		JWTExpiry: 24, // hours
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
} 