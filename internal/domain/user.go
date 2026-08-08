package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
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

func NewUser(email, password, displayName string) (User, error) {
	if !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}
	if utf8.RuneCountInString(password) < 8 {
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

type UserProfile struct {
	ID          uuid.UUID  `json:"id"`
	DisplayName string     `json:"display_name"`
	AvatarURL   *string    `json:"avatar_url"`
	Gender      Gender     `json:"gender"`
	Email       *string    `json:"email,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (u *User) ToProfile(requesterID uuid.UUID) UserProfile {
	profile := UserProfile{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		Gender:      u.Gender,
		CreatedAt:   u.CreatedAt,
	}
	if u.ID == requesterID {
		profile.Email = &u.Email
		profile.BirthDate = u.BirthDate
	}
	return profile
}
