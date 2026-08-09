package user

import (
	"time"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
)

// UpdateProfileInput is the domain command for updating a user profile.
type UpdateProfileInput struct {
	UserID      uuid.UUID
	DisplayName *string
	Gender      *authdomain.Gender
	BirthDate   *time.Time
	AvatarURL   *string
}

type UserProfile struct {
	ID          uuid.UUID
	DisplayName string
	AvatarURL   *string
	Gender      authdomain.Gender
	Email       *string
	BirthDate   *time.Time
	CreatedAt   time.Time
}

// ToProfile builds the profile read model. Email and birth date are private
// and only visible to the profile owner.
func ToProfile(u authdomain.User, requesterID uuid.UUID) UserProfile {
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
