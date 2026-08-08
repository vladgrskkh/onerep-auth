package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

func WriteDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		WriteError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		WriteError(w, http.StatusConflict, "email already exists")
	case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrInvalidPassword):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrTokenNotFound):
		WriteError(w, http.StatusNotFound, err.Error())
	default:
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func DecodeAndValidate(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &domain.ValidationError{Msg: "invalid request body"}
	}
	return nil
}
