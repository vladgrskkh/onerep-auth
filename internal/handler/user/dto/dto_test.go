package dto_test

import (
	"encoding/json"
	"testing"
	"time"

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

func TestDTOSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DTOTestSuite))
}
