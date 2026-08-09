package dto

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`
}
