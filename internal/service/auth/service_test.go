package auth_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/crypto"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"
	"github.com/vladgrskkh/onerep-auth/internal/service/auth"
	authmocks "github.com/vladgrskkh/onerep-auth/internal/service/auth/mocks"
)

type ServiceTestSuite struct {
	suite.Suite

	svc        *auth.AuthService
	userRepo   *authmocks.MockUserRepository
	tokenStore *authmocks.MockTokenStorer
}

func (s *ServiceTestSuite) SetupTest() {
	key, err := jwtsvc.GenerateKeyPair()
	s.Require().NoError(err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	s.Require().NoError(err)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	tm, err := jwtsvc.NewTokenManager(string(privPEM), 15*time.Minute)
	s.Require().NoError(err)

	s.userRepo = authmocks.NewMockUserRepository(s.T())
	s.tokenStore = authmocks.NewMockTokenStorer(s.T())
	hasher := *crypto.NewPasswordHasher()

	s.svc = auth.NewAuthService(s.userRepo, hasher, *tm, s.tokenStore)
}

func (s *ServiceTestSuite) TestRegister_Success() {
	s.userRepo.EXPECT().
		FindByEmail(mock.Anything, "test@example.com").
		Return(authdomain.User{}, authdomain.ErrUserNotFound)
	s.userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("auth.User")).Return(authdomain.User{}, nil)
	s.tokenStore.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	pair, err := s.svc.Register(context.Background(), "test@example.com", "password123", "Test User")
	s.Require().NoError(err)
	s.NotEmpty(pair.AccessToken)
	s.NotEmpty(pair.RefreshToken)
	s.Equal(900, pair.ExpiresIn)
}

func (s *ServiceTestSuite) TestRegister_DuplicateEmail() {
	s.userRepo.EXPECT().FindByEmail(mock.Anything, "dupe@example.com").Return(authdomain.User{}, nil)

	_, err := s.svc.Register(context.Background(), "dupe@example.com", "password123", "Test User")
	s.ErrorIs(err, authdomain.ErrEmailAlreadyExists)
}

func (s *ServiceTestSuite) TestRegister_InvalidEmail() {
	s.userRepo.EXPECT().FindByEmail(mock.Anything, "bad").Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := s.svc.Register(context.Background(), "bad", "password123", "Test User")
	s.ErrorIs(err, authdomain.ErrInvalidEmail)
}

func (s *ServiceTestSuite) TestLogin_Success() {
	hasher := crypto.NewPasswordHasher()
	hash, _ := hasher.Hash("password123")
	userID := uuid.Must(uuid.NewV7())

	s.userRepo.EXPECT().FindByEmail(mock.Anything, "login@example.com").Return(authdomain.User{
		ID:           userID,
		Email:        "login@example.com",
		PasswordHash: &hash,
	}, nil)
	s.tokenStore.EXPECT().Save(mock.Anything, mock.Anything, userID.String()).Return(nil)

	pair, err := s.svc.Login(context.Background(), "login@example.com", "password123")
	s.Require().NoError(err)
	s.NotEmpty(pair.AccessToken)
}

func (s *ServiceTestSuite) TestLogin_InvalidCredentials() {
	s.userRepo.EXPECT().
		FindByEmail(mock.Anything, "bad@example.com").
		Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := s.svc.Login(context.Background(), "bad@example.com", "wrong")
	s.ErrorIs(err, authdomain.ErrInvalidCredentials)
}

func (s *ServiceTestSuite) TestLogout_Success() {
	s.tokenStore.EXPECT().Delete(mock.Anything, "some-refresh-token").Return(nil)

	err := s.svc.Logout(context.Background(), "some-refresh-token")
	s.NoError(err)
}

func (s *ServiceTestSuite) TestRefresh_Success() {
	userID := uuid.Must(uuid.NewV7())

	s.tokenStore.EXPECT().Get(mock.Anything, "valid-refresh").Return(userID.String(), nil)
	s.tokenStore.EXPECT().Delete(mock.Anything, "valid-refresh").Return(nil)
	s.userRepo.EXPECT().
		FindByID(mock.Anything, userID).
		Return(authdomain.User{ID: userID, Email: "refresh@example.com"}, nil)
	s.tokenStore.EXPECT().Save(mock.Anything, mock.Anything, userID.String()).Return(nil)

	pair, err := s.svc.Refresh(context.Background(), "valid-refresh")
	s.Require().NoError(err)
	s.NotEmpty(pair.AccessToken)
}

func (s *ServiceTestSuite) TestRefresh_TokenNotFound() {
	s.tokenStore.EXPECT().Get(mock.Anything, "expired-refresh").Return("", authdomain.ErrTokenNotFound)

	_, err := s.svc.Refresh(context.Background(), "expired-refresh")
	s.ErrorIs(err, authdomain.ErrTokenNotFound)
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServiceTestSuite))
}
