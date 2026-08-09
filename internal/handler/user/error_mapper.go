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
	errMsgInvalidBirthDate  = "invalid birth date"
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

func errInvalidBirthDate() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidBirthDate,
		Message:     errMsgInvalidBirthDate,
		UserMessage: errUserInvalidBirthDate,
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
	default:
		return http.StatusInternalServerError, handler.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}
}
