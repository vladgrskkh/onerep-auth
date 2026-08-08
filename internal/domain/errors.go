package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrOAuthAccountExists = errors.New("oauth account already linked")
	ErrTokenNotFound      = errors.New("refresh token not found or expired")
)

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}
