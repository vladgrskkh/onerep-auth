package dto_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"

	"github.com/vladgrskkh/onerep-auth/internal/handler/user/dto"
)

func timeNow() time.Time {
	return time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC)
}

type DTOTestSuite struct {
	suite.Suite
}

func (s *DTOTestSuite) TestUserProfileResponse_OmitZero() {
	resp := dto.UserProfileResponse{
		ID:          [16]byte{1, 2, 3},
		DisplayName: "Alice",
		Gender:      "other",
		CreatedAt:   timeNow(),
	}

	raw, err := json.Marshal(resp)
	s.Require().NoError(err)

	var m map[string]any
	s.Require().NoError(json.Unmarshal(raw, &m))

	s.NotContains(m, "email")
	s.NotContains(m, "birth_date")
	s.NotContains(m, "avatar_url")
	s.Equal("Alice", m["display_name"])
}

func (s *DTOTestSuite) TestUserProfileResponse_IncludesEmailForOwner() {
	resp := dto.UserProfileResponse{
		ID:          [16]byte{1, 2, 3},
		DisplayName: "Alice",
		Gender:      "female",
		Email:       "a@b.com",
		CreatedAt:   timeNow(),
	}

	raw, err := json.Marshal(resp)
	s.Require().NoError(err)

	var m map[string]any
	s.Require().NoError(json.Unmarshal(raw, &m))

	s.Equal("a@b.com", m["email"])
}

func (s *DTOTestSuite) TestUpdateProfileRequest_Validation() {
	cases := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{name: "empty body", body: `{}`, wantErr: false},
		{
			name: "all fields valid",
			body: `{"display_name":"Alice","gender":"female","birth_date":"2000-01-02"}`,
		},
		{name: "invalid gender", body: `{"gender":"attack"}`, wantErr: true},
		{name: "invalid birth date format", body: `{"birth_date":"01/02/2000"}`, wantErr: true},
		{
			name:    "display name too long",
			body:    fmt.Sprintf(`{"display_name":%q}`, strings.Repeat("a", 101)),
			wantErr: true,
		},
	}

	validate := validator.New()
	for _, tc := range cases {
		s.Run(tc.name, func() {
			var req dto.UpdateProfileRequest
			s.Require().NoError(json.Unmarshal([]byte(tc.body), &req))
			err := validate.Struct(req)
			if tc.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestDTOSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DTOTestSuite))
}
