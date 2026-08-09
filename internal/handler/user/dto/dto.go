package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`
}

type UserProfile struct {
	ID          uuid.UUID  `json:"id"`
	DisplayName string     `json:"display_name"`
	AvatarURL   *string    `json:"avatar_url"`
	Gender      string     `json:"gender"`
	Email       *string    `json:"email,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
