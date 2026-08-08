package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/vladgrskkh/onerep-auth/internal/domain/auth"
)

type DomainTestSuite struct {
	suite.Suite
}

func (s *DomainTestSuite) TestNewUser() {
	u, err := auth.NewUser("test@example.com", "password123", "Test User")
	s.Require().NoError(err)
	s.Equal("test@example.com", u.Email)
	s.Equal("Test User", u.DisplayName)
	s.Equal(auth.GenderOther, u.Gender)
	s.Equal(1, u.Version)
}

func (s *DomainTestSuite) TestNewUser_InvalidEmail() {
	_, err := auth.NewUser("notanemail", "password123", "Test User")
	s.ErrorIs(err, auth.ErrInvalidEmail)
}

func (s *DomainTestSuite) TestNewUser_ShortPassword() {
	_, err := auth.NewUser("test@example.com", "1234567", "Test User")
	s.ErrorIs(err, auth.ErrInvalidPassword)
}

func (s *DomainTestSuite) TestNewUser_EmptyName() {
	_, err := auth.NewUser("test@example.com", "password123", "")
	s.Error(err)
}

func (s *DomainTestSuite) TestUser_ToProfile_Self() {
	u, _ := auth.NewUser("a@b.com", "password123", "Alice")
	profile := u.ToProfile(u.ID)
	s.Equal("a@b.com", *profile.Email)
}

func (s *DomainTestSuite) TestUser_ToProfile_Other() {
	u, _ := auth.NewUser("a@b.com", "password123", "Alice")
	otherID := uuid.New()
	profile := u.ToProfile(otherID)
	s.Nil(profile.Email)
}

func TestDomainSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DomainTestSuite))
}
