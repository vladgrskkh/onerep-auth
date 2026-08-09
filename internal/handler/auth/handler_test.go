package auth_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/auth"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"
	authsvc "github.com/vladgrskkh/onerep-auth/internal/service/auth"
)

type fakeAuthService struct {
	registered bool
	lastCmd    authsvc.RegisterCommand
	pair       jwtsvc.TokenPair
}

func (f *fakeAuthService) Register(_ context.Context, cmd authsvc.RegisterCommand) (jwtsvc.TokenPair, error) {
	f.registered = true
	f.lastCmd = cmd
	return f.pair, nil
}

func (f *fakeAuthService) Login(_ context.Context, _ authsvc.LoginCommand) (jwtsvc.TokenPair, error) {
	return f.pair, nil
}

func (f *fakeAuthService) Logout(_ context.Context, _ authsvc.LogoutCommand) error {
	return nil
}

func (f *fakeAuthService) Refresh(_ context.Context, _ authsvc.RefreshCommand) (jwtsvc.TokenPair, error) {
	return f.pair, nil
}

type HandlerTestSuite struct {
	suite.Suite

	handler *auth.AuthHandler
	fake    *fakeAuthService
}

func (s *HandlerTestSuite) SetupTest() {
	s.fake = &fakeAuthService{}
	logger := slog.New(slog.DiscardHandler)
	s.handler = auth.NewAuthHandler(s.fake, logger)
}

func (s *HandlerTestSuite) post(body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(body))
	w := httptest.NewRecorder()
	s.handler.Register(w, req)
	return w
}

func (s *HandlerTestSuite) errorCode(w *httptest.ResponseRecorder) string {
	var resp handler.ErrorResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	return string(resp.Error.Code)
}

func (s *HandlerTestSuite) TestRegister_InvalidEmail() {
	w := s.post(`{"email":"notanemail","password":"password123","display_name":"Alice"}`)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("VALIDATION_ERROR", s.errorCode(w))
	s.False(s.fake.registered)
}

func (s *HandlerTestSuite) TestRegister_MissingPassword() {
	w := s.post(`{"email":"a@b.com"}`)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("VALIDATION_ERROR", s.errorCode(w))
	s.False(s.fake.registered)
}

func (s *HandlerTestSuite) TestRegister_InvalidJSON() {
	w := s.post(`{not json`)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("INVALID_REQUEST_BODY", s.errorCode(w))
	s.False(s.fake.registered)
}

func (s *HandlerTestSuite) TestRegister_Success() {
	s.fake.pair = jwtsvc.TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 900}
	w := s.post(`{"email":"a@b.com","password":"password123","display_name":"Alice"}`)

	s.Equal(http.StatusCreated, w.Code)
	s.True(s.fake.registered)
	s.Equal("a@b.com", s.fake.lastCmd.Email)
	s.Equal("Alice", s.fake.lastCmd.DisplayName)
}

func TestHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HandlerTestSuite))
}
