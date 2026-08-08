package jwt_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	infrajwt "github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
)

type JWTTestSuite struct {
	suite.Suite

	tm infrajwt.TokenManager
}

func (s *JWTTestSuite) SetupTest() {
	key, err := infrajwt.GenerateKeyPair()
	s.Require().NoError(err)

	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	s.Require().NoError(err)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	tm, err := infrajwt.NewTokenManager(string(privPEM), 15*time.Minute)
	s.Require().NoError(err)
	s.tm = *tm
}

func (s *JWTTestSuite) TestIssueAndValidate() {
	userID := uuid.Must(uuid.NewV7())
	email := "test@example.com"

	pair, err := s.tm.IssueTokenPair(userID, email)
	s.Require().NoError(err)
	s.NotEmpty(pair.AccessToken)
	s.NotEmpty(pair.RefreshToken)
	s.Equal(900, pair.ExpiresIn)

	id, err := infrajwt.GetUserIDFromToken(pair.AccessToken, s.tm.PublicKeyPEM())
	s.Require().NoError(err)
	s.Equal(userID, id)
}

func TestJWTSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(JWTTestSuite))
}
