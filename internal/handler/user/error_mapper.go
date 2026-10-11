package user

import (
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

const (
	errCodeInvalidRequestBody = "INVALID_REQUEST_BODY"
	errMsgInvalidRequestBody  = "invalid request body"
	errUserInvalidRequestBody = "The request body is invalid"

	errCodeInvalidUserID = "INVALID_USER_ID"
	errMsgInvalidUserID  = "invalid user id"
	errUserInvalidUserID = "The user ID is invalid"

	errCodeForbidden = "FORBIDDEN"
	errMsgForbidden  = "cannot update another user's profile"
	errUserForbidden = "You can only update your own profile"

	errCodeUserNotFound = "USER_NOT_FOUND"
	errUserUserNotFound = "User not found"

	errCodeInvalidBirthDate = "INVALID_BIRTH_DATE"
	errUserInvalidBirthDate = "Birth date must be in YYYY-MM-DD format"

	errCodeInvalidGender = "INVALID_GENDER"
	errUserInvalidGender = "Please select a valid gender"

	errCodeInvalidDisplayName = "INVALID_DISPLAY_NAME"
	errUserInvalidDisplayName = "Display name must be between 1 and 100 characters"
)

func invalidRequestBodyDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidRequestBody,
		Message:     errMsgInvalidRequestBody,
		UserMessage: errUserInvalidRequestBody,
	}
}

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

func userNotFoundDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeUserNotFound,
		Message:     err.Error(),
		UserMessage: errUserUserNotFound,
	}
}

func invalidBirthDateDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidBirthDate,
		Message:     err.Error(),
		UserMessage: errUserInvalidBirthDate,
	}
}

func invalidGenderDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidGender,
		Message:     err.Error(),
		UserMessage: errUserInvalidGender,
	}
}

func invalidDisplayNameDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidDisplayName,
		Message:     err.Error(),
		UserMessage: errUserInvalidDisplayName,
	}
}
