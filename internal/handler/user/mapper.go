package user

import (
	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
	"github.com/vladgrskkh/onerep-auth/internal/handler/user/dto"
)

// toProfileResponse maps the domain read model to the HTTP response.
func toProfileResponse(p userdomain.UserProfile) dto.UserProfileResponse {
	resp := dto.UserProfileResponse{
		ID:          p.ID,
		DisplayName: p.DisplayName,
		Gender:      string(p.Gender),
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
