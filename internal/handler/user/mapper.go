package user

import (
	"time"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/user/dto"
)

// toProfileResponse maps the domain profile to the HTTP response.
func toProfileResponse(p userdomain.UserProfile) dto.UserProfileResponse {
	resp := dto.UserProfileResponse{
		ID:          p.ID,
		DisplayName: *p.DisplayName,
		Gender:      string(*p.Gender),
		CreatedAt:   p.CreatedAt,
	}
	if p.AvatarURL != nil {
		resp.AvatarURL = *p.AvatarURL
	}
	if p.Email != nil {
		resp.Email = *p.Email
	}
	if p.BirthDate != nil {
		resp.BirthDate = *p.BirthDate
	}
	return resp
}

// toUpdateProfile maps the HTTP request to the domain profile update command.
func toUpdateProfile(req dto.UpdateProfileRequest, userID uuid.UUID) (userdomain.UserProfile, handler.ErrorDetail) {
	profile := userdomain.UserProfile{
		ID:          userID,
		DisplayName: req.DisplayName,
	}
	if req.Gender != nil {
		g := authdomain.Gender(*req.Gender)
		profile.Gender = &g
	}
	if req.BirthDate != nil {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return userdomain.UserProfile{}, errInvalidBirthDate()
		}
		profile.BirthDate = &t
	}
	return profile, handler.ErrorDetail{}
}
