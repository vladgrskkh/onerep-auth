package user

import (
	"errors"
	"net/http"

	"github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

func mapError(err error) (int, handler.ErrorDetail) {
	switch {
	case errors.Is(err, auth.ErrUserNotFound):
		return http.StatusNotFound, handler.ErrorDetail{
			Code:        "USER_NOT_FOUND",
			Message:     "user not found",
			UserMessage: "User not found",
			Type:        handler.ErrorTypeUser,
		}
	default:
		return http.StatusInternalServerError, handler.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
			Type:    handler.ErrorTypeSystem,
		}
	}
}
