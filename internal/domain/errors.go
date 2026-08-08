package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrOAuthAccountExists = errors.New("oauth account already linked")
	ErrCannotUpdateOther  = errors.New("cannot update another user's profile")
	ErrTokenNotFound      = errors.New("refresh token not found or expired")
)
