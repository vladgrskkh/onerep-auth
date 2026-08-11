package user

import (
	"github.com/google/uuid"

	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
	"github.com/vladgrskkh/onerep-auth/internal/handler/user/dto"
	serviceuser "github.com/vladgrskkh/onerep-auth/internal/service/user"
)

// toProfileResponse maps the domain profile to the HTTP response.
func toProfileResponse(p *userdomain.UserProfile) dto.UserProfileResponse {
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

// toUpdateProfile maps the HTTP request to the service update command.
func toUpdateProfile(req dto.UpdateProfileRequest, userID uuid.UUID) serviceuser.UpdateProfileCommand {
	return serviceuser.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: req.DisplayName,
		Gender:      req.Gender,
		BirthDate:   req.BirthDate,
		AvatarURL:   req.AvatarURL,
	}
}
