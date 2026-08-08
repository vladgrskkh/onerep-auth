package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

func TestNewUser(t *testing.T) {
	u, err := auth.NewUser("test@example.com", "password123", "Test User")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", u.Email)
	assert.Equal(t, "Test User", u.DisplayName)
	assert.Equal(t, auth.GenderOther, u.Gender)
	assert.Equal(t, 1, u.Version)
}

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := auth.NewUser("notanemail", "password123", "Test User")
	assert.ErrorIs(t, err, domain.ErrInvalidEmail)
}

func TestNewUser_ShortPassword(t *testing.T) {
	_, err := auth.NewUser("test@example.com", "1234567", "Test User")
	assert.ErrorIs(t, err, domain.ErrInvalidPassword)
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := auth.NewUser("test@example.com", "password123", "")
	assert.Error(t, err)
}

func TestUser_ToProfile_Self(t *testing.T) {
	u, _ := auth.NewUser("a@b.com", "password123", "Alice")
	profile := u.ToProfile(u.ID)
	assert.Equal(t, "a@b.com", *profile.Email)
}

func TestUser_ToProfile_Other(t *testing.T) {
	u, _ := auth.NewUser("a@b.com", "password123", "Alice")
	otherID := uuid.New()
	profile := u.ToProfile(otherID)
	assert.Nil(t, profile.Email)
}
