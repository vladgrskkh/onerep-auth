package auth

import (
	"errors"
	"net/http"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

func mapError(err error) (int, handler.ErrorDetail) {
	switch {
	case errors.Is(err, authdomain.ErrInvalidCredentials):
		return http.StatusUnauthorized, handler.ErrorDetail{
			Code:        "INVALID_CREDENTIALS",
			Message:     "invalid credentials",
			UserMessage: "The email or password you entered is incorrect",
			Type:        handler.ErrorTypeUser,
		}
	case errors.Is(err, authdomain.ErrEmailAlreadyExists):
		return http.StatusConflict, handler.ErrorDetail{
			Code:        "EMAIL_ALREADY_EXISTS",
			Message:     "email already exists",
			UserMessage: "An account with this email already exists",
			Type:        handler.ErrorTypeUser,
		}
	case errors.Is(err, authdomain.ErrInvalidEmail):
		return http.StatusBadRequest, handler.ErrorDetail{
			Code:        "INVALID_EMAIL",
			Message:     "invalid email format",
			UserMessage: "Please enter a valid email address",
			Type:        handler.ErrorTypeUser,
		}
	case errors.Is(err, authdomain.ErrInvalidPassword):
		return http.StatusBadRequest, handler.ErrorDetail{
			Code:        "INVALID_PASSWORD",
			Message:     "password must be at least 8 characters",
			UserMessage: "Password must be at least 8 characters long",
			Type:        handler.ErrorTypeUser,
		}
	case errors.Is(err, authdomain.ErrUserNotFound):
		return http.StatusNotFound, handler.ErrorDetail{
			Code:        "USER_NOT_FOUND",
			Message:     "user not found",
			UserMessage: "User not found",
			Type:        handler.ErrorTypeUser,
		}
	case errors.Is(err, authdomain.ErrTokenNotFound):
		return http.StatusUnauthorized, handler.ErrorDetail{
			Code:        "TOKEN_NOT_FOUND",
			Message:     "refresh token not found or expired",
			UserMessage: "Your session has expired, please log in again",
			Type:        handler.ErrorTypeUser,
		}
	default:
		return http.StatusInternalServerError, handler.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
			Type:    handler.ErrorTypeSystem,
		}
	}
}
