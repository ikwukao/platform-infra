// Package config provides application configuration for Platform-Infra.
package config

import "os"

// Config contains runtime configuration for the Platform-Infra controller.
type Config struct {
	ServerPort  string
	LogLevel    string
	DatabaseURL string
	RedisURL    string
}

// Load reads configuration from environment variables and applies defaults.
func Load() Config {
	return Config{
		ServerPort:  getEnv("PLATFORM_PORT", "9091"),
		LogLevel:    getEnv("PLATFORM_LOG_LEVEL", "info"),
		DatabaseURL: getEnv("PLATFORM_DATABASE_URL", "postgres://platform:platform@localhost:5433/platform?sslmode=disable"),
		RedisURL:    getEnv("PLATFORM_REDIS_URL", "redis://localhost:6380"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
