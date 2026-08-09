package jwt

import (
	"crypto/rsa"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwtlib.RegisteredClaims

	Email string `json:"email"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type TokenManager struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	ttl        time.Duration
}
