package user_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
	"github.com/vladgrskkh/onerep-auth/internal/handler"
	"github.com/vladgrskkh/onerep-auth/internal/handler/user"
	usermocks "github.com/vladgrskkh/onerep-auth/internal/handler/user/mocks"
	serviceuser "github.com/vladgrskkh/onerep-auth/internal/service/user"
)

type HandlerTestSuite struct {
	suite.Suite

	handler *user.UserHandler
	svc     *usermocks.MockUserProfileService
}

func (s *HandlerTestSuite) SetupTest() {
	s.svc = usermocks.NewMockUserProfileService(s.T())
	logger := slog.New(slog.DiscardHandler)
	s.handler = user.NewUserHandler(s.svc, logger)
}

func (s *HandlerTestSuite) patch(userID uuid.UUID, body string) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Patch("/v1/users/{id}", s.handler.UpdateProfile)

	req := httptest.NewRequest(http.MethodPatch, "/v1/users/"+userID.String(), strings.NewReader(body))
	req = req.WithContext(handler.WithUserID(req.Context(), userID))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func (s *HandlerTestSuite) decodeError(w *httptest.ResponseRecorder) handler.ErrorResponse {
	var resp handler.ErrorResponse
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

func (s *HandlerTestSuite) TestUpdateProfile_InvalidGender() {
	userID := uuid.Must(uuid.NewV7())
	w := s.patch(userID, `{"gender":"attack"}`)

	s.Equal(http.StatusBadRequest, w.Code)
	resp := s.decodeError(w)
	s.Equal("VALIDATION_ERROR", string(resp.Error.Code))
	s.Require().Len(resp.Error.Details, 1)
	s.Equal("gender", resp.Error.Details[0].Field)
	s.Equal("oneof", resp.Error.Details[0].Tag)
	s.svc.AssertNotCalled(s.T(), "UpdateProfile", mock.Anything, mock.Anything)
}

func (s *HandlerTestSuite) TestUpdateProfile_InvalidBirthDate() {
	userID := uuid.Must(uuid.NewV7())
	w := s.patch(userID, `{"birth_date":"01/02/2000"}`)

	s.Equal(http.StatusBadRequest, w.Code)
	resp := s.decodeError(w)
	s.Equal("VALIDATION_ERROR", string(resp.Error.Code))
	s.Require().Len(resp.Error.Details, 1)
	s.Equal("birth_date", resp.Error.Details[0].Field)
	s.Equal("datetime", resp.Error.Details[0].Tag)
	s.svc.AssertNotCalled(s.T(), "UpdateProfile", mock.Anything, mock.Anything)
}

func (s *HandlerTestSuite) TestUpdateProfile_DisplayNameTooLong() {
	userID := uuid.Must(uuid.NewV7())
	w := s.patch(userID, `{"display_name":"`+strings.Repeat("a", 101)+`"}`)

	s.Equal(http.StatusBadRequest, w.Code)
	resp := s.decodeError(w)
	s.Equal("VALIDATION_ERROR", string(resp.Error.Code))
	s.Require().Len(resp.Error.Details, 1)
	s.Equal("display_name", resp.Error.Details[0].Field)
	s.Equal("max", resp.Error.Details[0].Tag)
	s.svc.AssertNotCalled(s.T(), "UpdateProfile", mock.Anything, mock.Anything)
}

func (s *HandlerTestSuite) TestUpdateProfile_Success() {
	userID := uuid.Must(uuid.NewV7())
	cmd := serviceuser.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: new("Alice"),
		Gender:      new("female"),
		BirthDate:   new("2000-01-02"),
	}
	profile := userdomain.UserProfile{
		ID:          userID,
		DisplayName: new("Alice"),
		AvatarURL:   new("avatar.png"),
		Gender:      new(authdomain.GenderFemale),
		BirthDate:   new(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)),
		CreatedAt:   time.Now(),
	}
	s.svc.EXPECT().
		UpdateProfile(mock.Anything, cmd).
		Return(profile, nil)

	w := s.patch(userID, `{"display_name":"Alice","gender":"female","birth_date":"2000-01-02"}`)

	s.Equal(http.StatusOK, w.Code)
	var resp struct {
		DisplayName string `json:"display_name"`
		Gender      string `json:"gender"`
	}
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	s.Equal("Alice", resp.DisplayName)
	s.Equal("female", resp.Gender)
	s.svc.AssertCalled(s.T(), "UpdateProfile", mock.Anything, cmd)
}

func TestHandlerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HandlerTestSuite))
}
