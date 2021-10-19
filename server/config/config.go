package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type FrontendConfig struct {
	PathToDist string
	IndexPath  string
}

type HttpServerConfig struct {
	Host         string
	Port         int
	ShutdownWait time.Duration
}

type Config struct {
	HttpServer HttpServerConfig
	Frontend   FrontendConfig
}

func New() *Config {
	return &Config{
		HttpServer: HttpServerConfig{
			Host:         getEnvAsString("SERVER_HOST", ""),
			Port:         getEnvAsInt("SERVER_PORT", 80),
			ShutdownWait: getEnvAsDuration("SHUTDOWN_WAIT", time.Second*15),
		},
		Frontend: FrontendConfig{
			PathToDist: getEnvAsString("FRONTEND_PATH", "frontend/build"),
			IndexPath:  getEnvAsString("FRONTEND_INDEX", "index.html"),
		},
	}
}

func getEnvAsString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	if valueStr := getEnvAsString(key, ""); valueStr != "" {
		if res, err := strconv.ParseBool(valueStr); err == nil {
			return res
		}
		log.Printf("Error to parse environment variable: %s\n", key)
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valueStr := getEnvAsString(key, ""); valueStr != "" {
		if res, err := strconv.Atoi(valueStr); err == nil {
			return res
		}
		log.Printf("Error to parse environment variable: %s\n", key)
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valueStr := getEnvAsString(key, ""); valueStr != "" {
		if res, err := time.ParseDuration(valueStr); err == nil {
			return res
		}
		log.Printf("Error to parse environment variable: %s\n", key)
	}
	return defaultVal
}
