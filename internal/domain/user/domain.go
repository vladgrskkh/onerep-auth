package user

import (
	"time"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
)

// UserProfile is the domain type for the user profile. It serves both as the
// read model and as the update command: nil fields are not updated.
type UserProfile struct {
	ID          uuid.UUID
	DisplayName *string
	AvatarURL   *string
	Gender      *authdomain.Gender
	Email       *string
	BirthDate   *time.Time
	CreatedAt   time.Time
}

// ToProfile builds the profile read model. Email and birth date are private
// and only visible to the profile owner.
func ToProfile(u authdomain.User, requesterID uuid.UUID) UserProfile {
	displayName := u.DisplayName
	gender := u.Gender
	profile := UserProfile{
		ID:          u.ID,
		DisplayName: &displayName,
		AvatarURL:   u.AvatarURL,
		Gender:      &gender,
		CreatedAt:   u.CreatedAt,
	}
	if u.ID == requesterID {
		email := u.Email
		profile.Email = &email
		profile.BirthDate = u.BirthDate
	}
	return profile
}
