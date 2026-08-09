package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL       string        `env:"DATABASE_URL,notEmpty" envDefault:"postgres://gym:gym_pass@localhost:5432/gym?sslmode=disable"`
	RedisURL          string        `env:"REDIS_URL,notEmpty"    envDefault:"redis://localhost:6379/0"`
	Port              string        `env:"AUTH_PORT"             envDefault:"8080"`
	JWTPrivateKeyPEM  string        `env:"JWT_PRIVATE_KEY"`
	JWTPublicKeyPEM   string        `env:"JWT_PUBLIC_KEY"`
	JWTTokenTTL       time.Duration `env:"JWT_TOKEN_TTL"         envDefault:"15m"`
	RefreshTokenTTL   time.Duration `env:"REFRESH_TOKEN_TTL"     envDefault:"168h"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"   envDefault:"5s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT"      envDefault:"10s"`
	GoogleClientID    string        `env:"GOOGLE_CLIENT_ID"`
	GoogleRedirectURL string        `env:"GOOGLE_REDIRECT_URL"   envDefault:"http://localhost:8080/v1/auth/oauth/google/callback"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
