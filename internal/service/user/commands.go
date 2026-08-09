package user

import "github.com/google/uuid"

// UpdateProfileCommand carries the profile fields to update. Nil pointers
// mean the field was not provided and stays unchanged.
type UpdateProfileCommand struct {
	UserID      uuid.UUID
	DisplayName *string
	Gender      *string
	BirthDate   *string
	AvatarURL   *string
}
