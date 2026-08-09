package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vladgrskkh/onerep-auth/internal/handler/middleware"

	"github.com/google/uuid"
)

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, middleware.UserIDKey, id)
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}

func DecodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
