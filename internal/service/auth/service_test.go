package auth_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/domain/auth/mocks"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
	"github.com/vladgrskkh/onerep-auth/internal/service/auth"
)

func setupAuthService(t *testing.T) (*auth.AuthService, *mocks.UserRepository, *mocks.TokenStorer) {
	t.Helper()

	key, err := jwtsvc.GenerateKeyPair()
	require.NoError(t, err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	tm, err := jwtsvc.NewTokenManager(string(privPEM), 15*time.Minute)
	require.NoError(t, err)

	userRepo := mocks.NewUserRepository(t)
	tokenStore := mocks.NewTokenStorer(t)
	hasher := *crypto.NewPasswordHasher()

	svc := auth.NewAuthService(userRepo, hasher, *tm, tokenStore)
	return svc, userRepo, tokenStore
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, userRepo, tokenStore := setupAuthService(t)

	userRepo.EXPECT().
		FindByEmail(mock.Anything, "test@example.com").
		Return(authdomain.User{}, authdomain.ErrUserNotFound)
	userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("auth.User")).Return(authdomain.User{}, nil)
	tokenStore.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	pair, err := svc.Register(context.Background(), "test@example.com", "password123", "Test User")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, 900, pair.ExpiresIn)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc, userRepo, _ := setupAuthService(t)

	userRepo.EXPECT().FindByEmail(mock.Anything, "dupe@example.com").Return(authdomain.User{}, nil)

	_, err := svc.Register(context.Background(), "dupe@example.com", "password123", "Test User")
	assert.ErrorIs(t, err, authdomain.ErrEmailAlreadyExists)
}

func TestAuthService_Register_InvalidEmail(t *testing.T) {
	svc, userRepo, _ := setupAuthService(t)

	userRepo.EXPECT().FindByEmail(mock.Anything, "bad").Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := svc.Register(context.Background(), "bad", "password123", "Test User")
	assert.ErrorIs(t, err, authdomain.ErrInvalidEmail)
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, userRepo, tokenStore := setupAuthService(t)

	hasher := crypto.NewPasswordHasher()
	hash, _ := hasher.Hash("password123")
	userID := uuid.Must(uuid.NewV7())

	userRepo.EXPECT().FindByEmail(mock.Anything, "login@example.com").Return(authdomain.User{
		ID:           userID,
		Email:        "login@example.com",
		PasswordHash: &hash,
	}, nil)
	tokenStore.EXPECT().Save(mock.Anything, mock.Anything, userID.String()).Return(nil)

	pair, err := svc.Login(context.Background(), "login@example.com", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	svc, userRepo, _ := setupAuthService(t)

	userRepo.EXPECT().
		FindByEmail(mock.Anything, "bad@example.com").
		Return(authdomain.User{}, authdomain.ErrUserNotFound)

	_, err := svc.Login(context.Background(), "bad@example.com", "wrong")
	assert.ErrorIs(t, err, authdomain.ErrInvalidCredentials)
}

func TestAuthService_Logout_Success(t *testing.T) {
	svc, _, tokenStore := setupAuthService(t)

	tokenStore.EXPECT().Delete(mock.Anything, "some-refresh-token").Return(nil)

	err := svc.Logout(context.Background(), "some-refresh-token")
	assert.NoError(t, err)
}

func TestAuthService_Refresh_Success(t *testing.T) {
	svc, userRepo, tokenStore := setupAuthService(t)

	userID := uuid.Must(uuid.NewV7())

	tokenStore.EXPECT().Get(mock.Anything, "valid-refresh").Return(userID.String(), nil)
	tokenStore.EXPECT().Delete(mock.Anything, "valid-refresh").Return(nil)
	userRepo.EXPECT().
		FindByID(mock.Anything, userID).
		Return(authdomain.User{ID: userID, Email: "refresh@example.com"}, nil)
	tokenStore.EXPECT().Save(mock.Anything, mock.Anything, userID.String()).Return(nil)

	pair, err := svc.Refresh(context.Background(), "valid-refresh")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
}

func TestAuthService_Refresh_TokenNotFound(t *testing.T) {
	svc, _, tokenStore := setupAuthService(t)

	tokenStore.EXPECT().Get(mock.Anything, "expired-refresh").Return("", authdomain.ErrTokenNotFound)

	_, err := svc.Refresh(context.Background(), "expired-refresh")
	assert.ErrorIs(t, err, authdomain.ErrTokenNotFound)
}
