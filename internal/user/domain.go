package user

import (
	"time"

	"github.com/vladgrskkh/onerep-auth/internal/auth"
)

type UpdateProfileInput struct {
	DisplayName *string      `json:"display_name,omitempty"`
	Gender      *auth.Gender `json:"gender,omitempty"`
	BirthDate   *time.Time   `json:"birth_date,omitempty"`
	AvatarURL   *string      `json:"avatar_url,omitempty"`
}
