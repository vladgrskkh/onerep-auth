package user

import (
	"time"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
)

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
