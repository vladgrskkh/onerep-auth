package auth

import (
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

const (
	errCodeInvalidRequestBody = "INVALID_REQUEST_BODY"
	errMsgInvalidRequestBody  = "invalid request body"
	errUserInvalidRequestBody = "The request body is invalid"

	errCodeEmailAlreadyExists = "EMAIL_ALREADY_EXISTS"
	errUserEmailAlreadyExists = "An account with this email already exists"

	errCodeInvalidEmail = "INVALID_EMAIL"
	errUserInvalidEmail = "Please enter a valid email address"

	errCodeInvalidPassword = "INVALID_PASSWORD"
	errUserInvalidPassword = "Password must be at least 8 characters long"

	errCodeInvalidDisplayName = "INVALID_DISPLAY_NAME"
	errUserInvalidDisplayName = "Display name must be between 1 and 100 characters"

	errCodeTokenNotFound = "TOKEN_NOT_FOUND"
	errUserTokenNotFound = "Your session has expired, please log in again"
)

func invalidRequestBodyDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidRequestBody,
		Message:     errMsgInvalidRequestBody,
		UserMessage: errUserInvalidRequestBody,
	}
}

func invalidCredentialsDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        "INVALID_CREDENTIALS",
		Message:     err.Error(),
		UserMessage: "The email or password you entered is incorrect",
	}
}

func emailAlreadyExistsDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeEmailAlreadyExists,
		Message:     err.Error(),
		UserMessage: errUserEmailAlreadyExists,
	}
}

func invalidEmailDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidEmail,
		Message:     err.Error(),
		UserMessage: errUserInvalidEmail,
	}
}

func invalidPasswordDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidPassword,
		Message:     err.Error(),
		UserMessage: errUserInvalidPassword,
	}
}

func invalidDisplayNameDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidDisplayName,
		Message:     err.Error(),
		UserMessage: errUserInvalidDisplayName,
	}
}

func tokenNotFoundDetail(err error) handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeTokenNotFound,
		Message:     err.Error(),
		UserMessage: errUserTokenNotFound,
	}
}
