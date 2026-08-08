package dto

type RegisterRequest struct {
	Email       string `json:"email"        validate:"required,email"`
	Password    string `json:"password"     validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required,min=1"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Gender      *string `json:"gender,omitempty"       validate:"omitempty,oneof=male female other"`
	BirthDate   *string `json:"birth_date,omitempty"`
}
