package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("TEST_VAR", "hello")
	assert.Equal(t, "hello", envOrDefault("TEST_VAR", "default"))
	assert.Equal(t, "default", envOrDefault("UNSET_VAR", "default"))
}
