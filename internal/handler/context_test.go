package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/vladgrskkh/onerep-auth/internal/handler"
)

type testValidationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ContextTestSuite struct {
	suite.Suite
}

func (s *ContextTestSuite) TestDecodeAndValidate_Success() {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"a@b.com"}`))
	var req testValidationRequest
	s.Require().NoError(handler.DecodeAndValidate(r, &req))
	s.Equal("a@b.com", req.Email)
}

func (s *ContextTestSuite) TestDecodeAndValidate_InvalidJSON() {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{not json`))
	var req testValidationRequest
	err := handler.DecodeAndValidate(r, &req)
	s.Require().Error(err)
	s.NotErrorIs(err, handler.ErrValidationFailed)
}

func (s *ContextTestSuite) TestDecodeAndValidate_ValidationFailed() {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"notanemail"}`))
	var req testValidationRequest
	err := handler.DecodeAndValidate(r, &req)
	s.ErrorIs(err, handler.ErrValidationFailed)
}

func (s *ContextTestSuite) TestDecodeAndValidate_MissingRequiredField() {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	var req testValidationRequest
	err := handler.DecodeAndValidate(r, &req)
	s.ErrorIs(err, handler.ErrValidationFailed)
}

func TestContextSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ContextTestSuite))
}
