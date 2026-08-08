package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

func NewTokenManager(privateKeyPEM string, ttl time.Duration) (*TokenManager, error) {
	key, err := jwtlib.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	return &TokenManager{privateKey: key, PublicKey: &key.PublicKey, ttl: ttl}, nil
}

func (m *TokenManager) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	return GetUserIDFromToken(tokenString, m.PublicKeyPEM())
}

func (m *TokenManager) GenerateAccessToken(user domain.User) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwtlib.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.ttl)),
		},
		Email: user.Email,
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}

func (m *TokenManager) PublicKeyPEM() string {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(m.PublicKey)
	if err != nil {
		return ""
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}))
}

func (m *TokenManager) IssueTokenPair(user domain.User) (TokenPair, error) {
	accessToken, err := m.GenerateAccessToken(user)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(m.ttl.Seconds()),
	}, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateKeyPair() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func GetUserIDFromToken(tokenString, publicKeyPEM string) (uuid.UUID, error) {
	key, err := jwtlib.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return uuid.Nil, err
	}

	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return uuid.Parse(claims.Subject)
}
