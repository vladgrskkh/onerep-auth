package auth_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"sync"
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

	svc          *auth.AuthService
	userRepo     *authmocks.MockUserRepository
	tokenStore   *authmocks.MockTokenStorer
	tokenManager jwtsvc.TokenManager
}

func (s *ServiceTestSuite) SetupTest() {
	key, err := jwtsvc.GenerateKeyPair()
	s.Require().NoError(err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	s.Require().NoError(err)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	tm, err := jwtsvc.NewTokenManager(string(privPEM), 15*time.Minute)
	s.Require().NoError(err)
	s.tokenManager = *tm

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

	pair, err := s.svc.Register(context.Background(), auth.RegisterCommand{
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	})
	s.Require().NoError(err)
	s.NotEmpty(pair.AccessToken)
	s.NotEmpty(pair.RefreshToken)
	s.Equal(900, pair.ExpiresIn)
}

func (s *ServiceTestSuite) TestRegister_DuplicateEmail() {
	s.userRepo.EXPECT().FindByEmail(mock.Anything, "dupe@example.com").Return(authdomain.User{}, nil)

	_, err := s.svc.Register(context.Background(), auth.RegisterCommand{
		Email:       "dupe@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	})
	s.ErrorIs(err, authdomain.ErrEmailAlreadyExists)
}

func (s *ServiceTestSuite) TestRegister_InvalidEmail() {
	s.userRepo.EXPECT().FindByEmail(mock.Anything, "").Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := s.svc.Register(context.Background(), auth.RegisterCommand{
		Email:       "",
		Password:    "password123",
		DisplayName: "Test User",
	})
	s.ErrorIs(err, authdomain.ErrInvalidEmail)
}

func (s *ServiceTestSuite) TestRegister_EmptyPassword() {
	s.userRepo.EXPECT().
		FindByEmail(mock.Anything, "test@example.com").
		Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := s.svc.Register(context.Background(), auth.RegisterCommand{
		Email:       "test@example.com",
		Password:    "",
		DisplayName: "Test User",
	})
	s.ErrorIs(err, authdomain.ErrInvalidPassword)
}

func (s *ServiceTestSuite) TestRegister_TooLongDisplayName() {
	s.userRepo.EXPECT().
		FindByEmail(mock.Anything, "test@example.com").
		Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := s.svc.Register(context.Background(), auth.RegisterCommand{
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: strings.Repeat("a", authdomain.MaxDisplayNameLength+1),
	})
	s.ErrorIs(err, authdomain.ErrInvalidDisplayName)
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

	pair, err := s.svc.Login(context.Background(), auth.LoginCommand{
		Email:    "login@example.com",
		Password: "password123",
	})
	s.Require().NoError(err)
	s.NotEmpty(pair.AccessToken)
}

func (s *ServiceTestSuite) TestLogin_InvalidCredentials() {
	s.userRepo.EXPECT().
		FindByEmail(mock.Anything, "bad@example.com").
		Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := s.svc.Login(context.Background(), auth.LoginCommand{
		Email:    "bad@example.com",
		Password: "wrong",
	})
	s.ErrorIs(err, authdomain.ErrInvalidCredentials)
}

func (s *ServiceTestSuite) TestLogout_Success() {
	s.tokenStore.EXPECT().Delete(mock.Anything, "some-refresh-token").Return(nil)

	err := s.svc.Logout(context.Background(), auth.LogoutCommand{RefreshToken: "some-refresh-token"})
	s.NoError(err)
}

func (s *ServiceTestSuite) TestRefresh_Success() {
	userID := uuid.Must(uuid.NewV7())

	s.tokenStore.EXPECT().Consume(mock.Anything, "valid-refresh").Return(userID.String(), nil)
	s.userRepo.EXPECT().
		FindByID(mock.Anything, userID).
		Return(authdomain.User{ID: userID, Email: "refresh@example.com"}, nil)
	s.tokenStore.EXPECT().Save(mock.Anything, mock.Anything, userID.String()).Return(nil)

	pair, err := s.svc.Refresh(context.Background(), auth.RefreshCommand{RefreshToken: "valid-refresh"})
	s.Require().NoError(err)
	s.NotEmpty(pair.AccessToken)
}

func (s *ServiceTestSuite) TestRefresh_TokenNotFound() {
	s.tokenStore.EXPECT().Consume(mock.Anything, "expired-refresh").Return("", authdomain.ErrTokenNotFound)

	_, err := s.svc.Refresh(context.Background(), auth.RefreshCommand{RefreshToken: "expired-refresh"})
	s.ErrorIs(err, authdomain.ErrTokenNotFound)
}

func (s *ServiceTestSuite) TestRefresh_ConcurrentRequestsConsumeTokenOnlyOnce() {
	userID := uuid.Must(uuid.NewV7())
	var mu sync.Mutex
	consumed := false
	s.tokenStore.EXPECT().Consume(mock.Anything, "shared-refresh").RunAndReturn(
		func(context.Context, string) (string, error) {
			mu.Lock()
			defer mu.Unlock()

			if consumed {
				return "", authdomain.ErrTokenNotFound
			}
			consumed = true
			return userID.String(), nil
		},
	)
	s.userRepo.EXPECT().FindByID(mock.Anything, userID).Return(
		authdomain.User{ID: userID, Email: "refresh@example.com"}, nil,
	).Once()
	s.tokenStore.EXPECT().Save(mock.Anything, mock.Anything, userID.String()).Return(nil).Once()

	svc := auth.NewAuthService(
		s.userRepo,
		*crypto.NewPasswordHasher(),
		s.tokenManager,
		s.tokenStore,
	)

	errs := make(chan error, 2)
	for range 2 {
		go func() {
			_, err := svc.Refresh(context.Background(), auth.RefreshCommand{RefreshToken: "shared-refresh"})
			errs <- err
		}()
	}

	var notFound int
	for range 2 {
		if err := <-errs; errors.Is(err, authdomain.ErrTokenNotFound) {
			notFound++
		} else {
			s.Require().NoError(err)
		}
	}

	s.Equal(1, notFound)
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServiceTestSuite))
}
