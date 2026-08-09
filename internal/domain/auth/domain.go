package auth

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrOAuthAccountExists = errors.New("oauth account already linked")
	ErrTokenNotFound      = errors.New("refresh token not found or expired")
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type OAuthProvider string

const (
	OAuthGoogle OAuthProvider = "google"
	OAuthApple  OAuthProvider = "apple"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash *string
	DisplayName  string
	AvatarURL    *string
	Gender       Gender
	BirthDate    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int
}

type OAuthAccount struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       OAuthProvider
	ProviderUserID string
	CreatedAt      time.Time
}

const minPasswordLength = 8

func NewUser(email, password, displayName string) (User, error) {
	if !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}
	if utf8.RuneCountInString(password) < minPasswordLength {
		return User{}, ErrInvalidPassword
	}
	if strings.TrimSpace(displayName) == "" {
		return User{}, ErrInvalidCredentials
	}

	now := time.Now()
	return User{
		ID:          uuid.Must(uuid.NewV7()),
		Email:       strings.ToLower(email),
		DisplayName: displayName,
		Gender:      GenderOther,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}, nil
}

func UserFromOAuth(email, displayName string) *User {
	now := time.Now()
	id := uuid.Must(uuid.NewV7())
	return &User{
		ID:          id,
		Email:       strings.ToLower(email),
		DisplayName: displayName,
		Gender:      GenderOther,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}
}

func NewOAuthAccount(userID uuid.UUID, provider OAuthProvider, providerUserID string) OAuthAccount {
	return OAuthAccount{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		CreatedAt:      time.Now(),
	}
}
