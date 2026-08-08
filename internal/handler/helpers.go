package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

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

type ctxKey string

const userIDKey ctxKey = "user_id"

func setUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func userIDFromContext(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userIDKey).(uuid.UUID)
	return id
}
