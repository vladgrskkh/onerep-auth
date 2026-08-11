package auth

import (
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

func invalidCredentialsDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        "INVALID_CREDENTIALS",
		Message:     "invalid credentials",
		UserMessage: "The email or password you entered is incorrect",
	}
}

const (
	errCodeEmailAlreadyExists = "EMAIL_ALREADY_EXISTS"
	errMsgEmailAlreadyExists  = "email already exists"
	errUserEmailAlreadyExists = "An account with this email already exists"
)

func emailAlreadyExistsDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeEmailAlreadyExists,
		Message:     errMsgEmailAlreadyExists,
		UserMessage: errUserEmailAlreadyExists,
	}
}

const (
	errCodeInvalidEmail = "INVALID_EMAIL"
	errMsgInvalidEmail  = "email is required"
	errUserInvalidEmail = "Please enter a valid email address"
)

func invalidEmailDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidEmail,
		Message:     errMsgInvalidEmail,
		UserMessage: errUserInvalidEmail,
	}
}

const (
	errCodeInvalidPassword = "INVALID_PASSWORD"
	errMsgInvalidPassword  = "password must be at least 8 characters"
	errUserInvalidPassword = "Password must be at least 8 characters long"
)

func invalidPasswordDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeInvalidPassword,
		Message:     errMsgInvalidPassword,
		UserMessage: errUserInvalidPassword,
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

const (
	errCodeTokenNotFound = "TOKEN_NOT_FOUND"
	errMsgTokenNotFound  = "refresh token not found or expired"
	errUserTokenNotFound = "Your session has expired, please log in again"
)

func tokenNotFoundDetail() handler.ErrorDetail {
	return handler.ErrorDetail{
		Code:        errCodeTokenNotFound,
		Message:     errMsgTokenNotFound,
		UserMessage: errUserTokenNotFound,
	}
}
