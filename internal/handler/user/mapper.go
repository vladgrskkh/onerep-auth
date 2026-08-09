package user

import (
	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/handler/user/dto"
)

// toProfile builds the profile response. Email and birth date are private and
// only visible to the profile owner.
func toProfile(u authdomain.User, requesterID uuid.UUID) dto.UserProfile {
	profile := dto.UserProfile{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		Gender:      string(u.Gender),
		CreatedAt:   u.CreatedAt,
	}
	if u.ID == requesterID {
		profile.Email = &u.Email
		profile.BirthDate = u.BirthDate
	}
	return profile
}
