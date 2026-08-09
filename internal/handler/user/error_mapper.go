package user

import (
	"errors"
	"net/http"

	"github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

const (
	errCodeInvalidUserID = "INVALID_USER_ID"
	errMsgInvalidUserID  = "invalid user id"
	errUserInvalidUserID = "The user ID is invalid"

	errCodeForbidden = "FORBIDDEN"
	errMsgForbidden  = "cannot update another user's profile"
	errUserForbidden = "You can only update your own profile"

	errCodeInvalidRequestBody = "INVALID_REQUEST_BODY"
	errMsgInvalidRequestBody  = "invalid request body"
	errUserInvalidRequestBody = "The request body is invalid"

	errCodeInvalidBirthDate = "INVALID_BIRTH_DATE"
	errUserInvalidBirthDate = "Birth date must be in YYYY-MM-DD format"
)

func invalidUserIDDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidUserID,
		Message:     errMsgInvalidUserID,
		UserMessage: errUserInvalidUserID,
	}
}

func forbiddenDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeForbidden,
		Message:     errMsgForbidden,
		UserMessage: errUserForbidden,
	}
}

func invalidRequestBodyDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidRequestBody,
		Message:     errMsgInvalidRequestBody,
		UserMessage: errUserInvalidRequestBody,
	}
}

func mapError(err error) (int, handler.ErrorDetail) {
	switch {
	case errors.Is(err, auth.ErrUserNotFound):
		return http.StatusNotFound, handler.ErrorDetail{
			Code:        "USER_NOT_FOUND",
			Message:     err.Error(),
			UserMessage: "User not found",
		}
	case errors.Is(err, auth.ErrInvalidBirthDate):
		return http.StatusBadRequest, handler.ErrorDetail{
			Code:        errCodeInvalidBirthDate,
			Message:     err.Error(),
			UserMessage: errUserInvalidBirthDate,
		}
	case errors.Is(err, auth.ErrInvalidGender):
		return http.StatusBadRequest, handler.ErrorDetail{
			Code:        "INVALID_GENDER",
			Message:     err.Error(),
			UserMessage: "Please select a valid gender",
		}
	case errors.Is(err, auth.ErrInvalidDisplayName):
		return http.StatusBadRequest, handler.ErrorDetail{
			Code:        "INVALID_DISPLAY_NAME",
			Message:     err.Error(),
			UserMessage: "Display name must be between 1 and 100 characters",
		}
	default:
		return http.StatusInternalServerError, handler.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}
}
