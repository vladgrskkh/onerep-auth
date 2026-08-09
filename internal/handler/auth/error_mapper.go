package auth

import (
	"errors"
	"net/http"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

const (
	errCodeInvalidRequestBody = "INVALID_REQUEST_BODY"
	errMsgInvalidRequestBody  = "invalid request body"
	errUserInvalidRequestBody = "The request body is invalid"
)

func invalidRequestBodyDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidRequestBody,
		Message:     errMsgInvalidRequestBody,
		UserMessage: errUserInvalidRequestBody,
	}
}

func mapError(err error) (int, handler.ErrorDetail) {
	switch {
	case errors.Is(err, authdomain.ErrInvalidCredentials):
		return http.StatusUnauthorized, handler.ErrorDetail{
			Code:        "INVALID_CREDENTIALS",
			Message:     err.Error(),
			UserMessage: "The email or password you entered is incorrect",
		}
	case errors.Is(err, authdomain.ErrEmailAlreadyExists):
		return http.StatusConflict, handler.ErrorDetail{
			Code:        "EMAIL_ALREADY_EXISTS",
			Message:     err.Error(),
			UserMessage: "An account with this email already exists",
		}
	case errors.Is(err, authdomain.ErrInvalidEmail):
		return http.StatusBadRequest, handler.ErrorDetail{
			Code:        "INVALID_EMAIL",
			Message:     err.Error(),
			UserMessage: "Please enter a valid email address",
		}
	case errors.Is(err, authdomain.ErrInvalidPassword):
		return http.StatusBadRequest, handler.ErrorDetail{
			Code:        "INVALID_PASSWORD",
			Message:     err.Error(),
			UserMessage: "Password must be at least 8 characters long",
		}
	case errors.Is(err, authdomain.ErrUserNotFound):
		return http.StatusNotFound, handler.ErrorDetail{
			Code:        "USER_NOT_FOUND",
			Message:     err.Error(),
			UserMessage: "User not found",
		}
	case errors.Is(err, authdomain.ErrTokenNotFound):
		return http.StatusUnauthorized, handler.ErrorDetail{
			Code:        "TOKEN_NOT_FOUND",
			Message:     err.Error(),
			UserMessage: "Your session has expired, please log in again",
		}
	default:
		return http.StatusInternalServerError, handler.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}
}
