package user

import (
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

const (
	errCodeInvalidUserID = "INVALID_USER_ID"
	errMsgInvalidUserID  = "invalid user id"
	errUserInvalidUserID = "The user ID is invalid"
)

func invalidUserIDDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidUserID,
		Message:     errMsgInvalidUserID,
		UserMessage: errUserInvalidUserID,
	}
}

const (
	errCodeForbidden = "FORBIDDEN"
	errMsgForbidden  = "cannot update another user's profile"
	errUserForbidden = "You can only update your own profile"
)

func forbiddenDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeForbidden,
		Message:     errMsgForbidden,
		UserMessage: errUserForbidden,
	}
}

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

const (
	errCodeUserNotFound = "USER_NOT_FOUND"
	errMsgUserNotFound  = "user not found"
	errUserUserNotFound = "User not found"
)

func userNotFoundDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeUserNotFound,
		Message:     errMsgUserNotFound,
		UserMessage: errUserUserNotFound,
	}
}

const (
	errCodeInvalidBirthDate = "INVALID_BIRTH_DATE"
	errMsgInvalidBirthDate  = "birth date must be in YYYY-MM-DD format"
	errUserInvalidBirthDate = "Birth date must be in YYYY-MM-DD format"
)

func invalidBirthDateDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidBirthDate,
		Message:     errMsgInvalidBirthDate,
		UserMessage: errUserInvalidBirthDate,
	}
}

const (
	errCodeInvalidGender = "INVALID_GENDER"
	errMsgInvalidGender  = "invalid gender"
	errUserInvalidGender = "Please select a valid gender"
)

func invalidGenderDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidGender,
		Message:     errMsgInvalidGender,
		UserMessage: errUserInvalidGender,
	}
}

const (
	errCodeInvalidDisplayName = "INVALID_DISPLAY_NAME"
	errMsgInvalidDisplayName  = "display name must be between 1 and 100 characters"
	errUserInvalidDisplayName = "Display name must be between 1 and 100 characters"
)

func invalidDisplayNameDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidDisplayName,
		Message:     errMsgInvalidDisplayName,
		UserMessage: errUserInvalidDisplayName,
	}
}
