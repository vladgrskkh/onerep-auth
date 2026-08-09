package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/domain/user"
)

type DomainTestSuite struct {
	suite.Suite
}

func (s *DomainTestSuite) TestToProfile_Self() {
	u, err := authdomain.NewUser("a@b.com", "password123", "Alice")
	s.Require().NoError(err)

	profile := user.ToProfile(u, u.ID)
	s.Equal("a@b.com", *profile.Email)
	s.Nil(profile.BirthDate)
}

func (s *DomainTestSuite) TestToProfile_Other() {
	u, err := authdomain.NewUser("a@b.com", "password123", "Alice")
	s.Require().NoError(err)

	profile := user.ToProfile(u, uuid.New())
	s.Nil(profile.Email)
	s.Nil(profile.BirthDate)
	s.Equal("Alice", profile.DisplayName)
}

func TestDomainSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DomainTestSuite))
}
