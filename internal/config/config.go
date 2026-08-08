package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	Port        string

	JWTPrivateKeyPEM string
	JWTPublicKeyPEM  string
	JWTTokenTTL      time.Duration

	RefreshTokenTTL time.Duration

	GoogleClientID    string
	GoogleRedirectURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://gym:gym_pass@localhost:5432/gym?sslmode=disable"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379/0"),
		Port:              getEnv("AUTH_PORT", "8080"),
		JWTPrivateKeyPEM:  getEnv("JWT_PRIVATE_KEY", ""),
		JWTPublicKeyPEM:   getEnv("JWT_PUBLIC_KEY", ""),
		JWTTokenTTL:       getEnvDuration("JWT_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL:   getEnvDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/v1/auth/oauth/google/callback"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return defaultVal
}
