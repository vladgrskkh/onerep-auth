package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

func withUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func userIDFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
