package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func Load() (*Config, error) {
	godotenv.Load()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBName:     getEnv("DB_NAME", "dev_book"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		RedisHost:     getEnv("REDIS_HOST", ""),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}

	if cfg.JWTSecret == "" {
		return nil, os.ErrInvalid
	}

	return cfg, nil
}

func (c *Config) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
}

func (c *Config) RedisEnabled() bool {
	return c.RedisHost != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
