package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"
)

type PasswordTestSuite struct {
	suite.Suite

	hasher crypto.PasswordHasher
}

func (s *PasswordTestSuite) SetupTest() {
	s.hasher = *crypto.NewPasswordHasher()
}

func (s *PasswordTestSuite) TestHashAndCompare() {
	hash, err := s.hasher.Hash("mypassword")
	s.Require().NoError(err)
	s.NotEmpty(hash)

	err = s.hasher.Compare(hash, "mypassword")
	s.Require().NoError(err)

	err = s.hasher.Compare(hash, "wrongpassword")
	s.Error(err)
}

func TestPasswordSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(PasswordTestSuite))
}
