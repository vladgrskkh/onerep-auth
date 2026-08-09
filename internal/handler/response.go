package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeUser   ErrorType = "user_error"
	ErrorTypeSystem ErrorType = "system_error"
)

type ErrorCode string

type ErrorDetail struct {
	Code        ErrorCode `json:"code"`
	Message     string    `json:"message"`
	UserMessage string    `json:"user_message,omitempty"`
	Type        ErrorType `json:"type"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// WriteJSON writes a JSON response. Once WriteHeader is called the response is
// committed; an encoding failure at this point is not recoverable server-side.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, detail ErrorDetail) {
	WriteJSON(w, status, ErrorResponse{Error: detail})
}

func WriteSystemError(w http.ResponseWriter, status int, code ErrorCode) {
	WriteError(w, status, ErrorDetail{
		Code:    code,
		Message: "internal server error",
		Type:    ErrorTypeSystem,
	})
}
