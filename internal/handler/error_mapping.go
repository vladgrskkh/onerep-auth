package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
)

func WriteDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authdomain.ErrInvalidCredentials):
		WriteError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, authdomain.ErrEmailAlreadyExists):
		WriteError(w, http.StatusConflict, "email already exists")
	case errors.Is(err, authdomain.ErrInvalidEmail), errors.Is(err, authdomain.ErrInvalidPassword):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, authdomain.ErrUserNotFound), errors.Is(err, authdomain.ErrTokenNotFound):
		WriteError(w, http.StatusNotFound, err.Error())
	default:
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func DecodeAndValidate(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &authdomain.ValidationError{Msg: "invalid request body"}
	}
	return nil
}
