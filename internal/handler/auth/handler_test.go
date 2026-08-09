package auth_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/auth"
	authmocks "github.com/vladgrskkh/onerep-auth/internal/handler/auth/mocks"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"
	authsvc "github.com/vladgrskkh/onerep-auth/internal/service/auth"
)

type HandlerTestSuite struct {
	suite.Suite

	handler *auth.AuthHandler
	svc     *authmocks.MockAuthService
}

func (s *HandlerTestSuite) SetupTest() {
	s.svc = authmocks.NewMockAuthService(s.T())
	logger := slog.New(slog.DiscardHandler)
	s.handler = auth.NewAuthHandler(s.svc, logger)
}

func (s *HandlerTestSuite) post(body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(body))
	w := httptest.NewRecorder()
	s.handler.Register(w, req)
	return w
}

func (s *HandlerTestSuite) decodeError(w *httptest.ResponseRecorder) handler.ErrorResponse {
	var resp handler.ErrorResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

func (s *HandlerTestSuite) TestRegister_InvalidEmail() {
	w := s.post(`{"email":"notanemail","password":"password123","display_name":"Alice"}`)

	s.Equal(http.StatusBadRequest, w.Code)
	resp := s.decodeError(w)
	s.Equal("VALIDATION_ERROR", string(resp.Error.Code))
	s.Equal("Please check your input", resp.Error.UserMessage)
	s.Require().Len(resp.Error.Details, 1)
	s.Equal("email", resp.Error.Details[0].Field)
	s.Equal("email", resp.Error.Details[0].Tag)
	s.svc.AssertNotCalled(s.T(), "Register", mock.Anything, mock.Anything)
}

func (s *HandlerTestSuite) TestRegister_MissingPassword() {
	w := s.post(`{"email":"a@b.com"}`)

	s.Equal(http.StatusBadRequest, w.Code)
	resp := s.decodeError(w)
	s.Equal("VALIDATION_ERROR", string(resp.Error.Code))
	s.Require().Len(resp.Error.Details, 2)
	s.Equal("password", resp.Error.Details[0].Field)
	s.Equal("required", resp.Error.Details[0].Tag)
	s.Equal("display_name", resp.Error.Details[1].Field)
	s.Equal("required", resp.Error.Details[1].Tag)
	s.svc.AssertNotCalled(s.T(), "Register", mock.Anything, mock.Anything)
}

func (s *HandlerTestSuite) TestRegister_InvalidJSON() {
	w := s.post(`{not json`)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("INVALID_REQUEST_BODY", string(s.decodeError(w).Error.Code))
	s.svc.AssertNotCalled(s.T(), "Register", mock.Anything, mock.Anything)
}

func (s *HandlerTestSuite) TestRegister_Success() {
	cmd := authsvc.RegisterCommand{
		Email:       "a@b.com",
		Password:    "password123",
		DisplayName: "Alice",
	}
	s.svc.EXPECT().
		Register(mock.Anything, cmd).
		Return(jwtsvc.TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 900}, nil)

	w := s.post(`{"email":"a@b.com","password":"password123","display_name":"Alice"}`)

	s.Equal(http.StatusCreated, w.Code)
	s.svc.AssertCalled(s.T(), "Register", mock.Anything, cmd)
}

func TestHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HandlerTestSuite))
}
