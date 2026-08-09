package auth_test

import (
	"strings"
	"testing"

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

func (s *DomainTestSuite) TestNewUser_EmptyEmail() {
	_, err := auth.NewUser("", "password123", "Test User")
	s.ErrorIs(err, auth.ErrInvalidEmail)
}

func (s *DomainTestSuite) TestNewUser_WhitespaceEmail() {
	_, err := auth.NewUser("   ", "password123", "Test User")
	s.ErrorIs(err, auth.ErrInvalidEmail)
}

// Email format is not a domain concern: the request validator enforces it.
func (s *DomainTestSuite) TestNewUser_NoEmailFormatCheck() {
	_, err := auth.NewUser("notanemail", "password123", "Test User")
	s.NoError(err)
}

func (s *DomainTestSuite) TestNewUser_ShortPassword() {
	_, err := auth.NewUser("test@example.com", "1234567", "Test User")
	s.ErrorIs(err, auth.ErrInvalidPassword)
}

func (s *DomainTestSuite) TestNewUser_EmptyName() {
	_, err := auth.NewUser("test@example.com", "password123", "")
	s.ErrorIs(err, auth.ErrInvalidDisplayName)
}

func (s *DomainTestSuite) TestNewUser_WhitespaceName() {
	_, err := auth.NewUser("test@example.com", "password123", "   ")
	s.ErrorIs(err, auth.ErrInvalidDisplayName)
}

func (s *DomainTestSuite) TestNewUser_LongDisplayName() {
	_, err := auth.NewUser("test@example.com", "password123", strings.Repeat("a", auth.MaxDisplayNameLength+1))
	s.ErrorIs(err, auth.ErrInvalidDisplayName)
}

func (s *DomainTestSuite) TestNewUser_MaxDisplayNameLength() {
	u, err := auth.NewUser("test@example.com", "password123", strings.Repeat("a", auth.MaxDisplayNameLength))
	s.Require().NoError(err)
	s.Len(u.DisplayName, auth.MaxDisplayNameLength)
}

func (s *DomainTestSuite) TestNewUser_TrimsDisplayName() {
	u, err := auth.NewUser("test@example.com", "password123", "  Test User  ")
	s.Require().NoError(err)
	s.Equal("Test User", u.DisplayName)
}

func (s *DomainTestSuite) TestNewUser_LowercasesEmail() {
	u, err := auth.NewUser("Test@Example.COM", "password123", "Test User")
	s.Require().NoError(err)
	s.Equal("test@example.com", u.Email)
}

func (s *DomainTestSuite) TestGenderIsValid() {
	s.True(auth.GenderMale.IsValid())
	s.True(auth.GenderFemale.IsValid())
	s.True(auth.GenderOther.IsValid())
	s.False(auth.Gender("").IsValid())
	s.False(auth.Gender("attack").IsValid())
}

func TestDomainSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DomainTestSuite))
}
