package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitzero"`
	Gender      *string `json:"gender,omitzero"`
	BirthDate   *string `json:"birth_date,omitzero"`
}

type UserProfileResponse struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url,omitzero"`
	Gender      string    `json:"gender"`
	Email       string    `json:"email,omitzero"`
	BirthDate   time.Time `json:"birth_date,omitzero"`
	CreatedAt   time.Time `json:"created_at"`
}
