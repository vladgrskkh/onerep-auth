package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		writeError(w, http.StatusConflict, "email already exists")
	case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrInvalidPassword):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrTokenNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func decodeAndValidate(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &domain.ValidationError{Msg: "invalid request body"}
	}
	// validator.Validate(v) call for structs with validate tags
	return nil
}
