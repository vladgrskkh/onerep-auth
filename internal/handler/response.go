package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorCode string

type ErrorDetail struct {
	Code        ErrorCode `json:"code"`
	Message     string    `json:"message"`
	UserMessage string    `json:"user_message,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// WriteJSON writes a JSON response. Once WriteHeader is called the response is
// committed; an encoding failure at this point is not recoverable server-side,
// so we only log it.
func WriteJSON(w http.ResponseWriter, logger *slog.Logger, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

func WriteError(w http.ResponseWriter, logger *slog.Logger, status int, detail ErrorDetail) {
	WriteJSON(w, logger, status, ErrorResponse{Error: detail})
}

func WriteSystemError(w http.ResponseWriter, logger *slog.Logger, status int, err error) {
	WriteError(w, logger, status, ErrorDetail{
		Code:    "INTERNAL_ERROR",
		Message: err.Error(),
	})
}

func ValidationErrorDetail() ErrorDetail {
	return ErrorDetail{
		Code:        "VALIDATION_ERROR",
		Message:     "request validation failed",
		UserMessage: "Please check your input",
	}
}
