package application_test

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

	"github.com/vladgrskkh/onerep-auth/internal/application"
	"github.com/vladgrskkh/onerep-auth/internal/domain"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

type mockTokenStore struct {
	mock.Mock
}

func (m *mockTokenStore) Save(ctx context.Context, token string, userID string) error {
	args := m.Called(ctx, token, userID)
	return args.Error(0)
}

func (m *mockTokenStore) Get(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *mockTokenStore) Delete(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func setupAuthService(t *testing.T) (*application.AuthService, *mockUserRepository, *mockTokenStore) {
	t.Helper()

	key, err := jwtsvc.GenerateKeyPair()
	require.NoError(t, err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	tm, err := jwtsvc.NewTokenManager(string(privPEM), 15*time.Minute)
	require.NoError(t, err)

	userRepo := new(mockUserRepository)
	tokenStore := new(mockTokenStore)
	hasher := *crypto.NewPasswordHasher()

	svc := application.NewAuthService(userRepo, hasher, *tm, tokenStore)
	return svc, userRepo, tokenStore
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, userRepo, tokenStore := setupAuthService(t)

	userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(domain.User{}, domain.ErrUserNotFound)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("domain.User")).Return(domain.User{}, nil)
	tokenStore.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	pair, err := svc.Register(context.Background(), "test@example.com", "password123", "Test User")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, 900, pair.ExpiresIn)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc, userRepo, _ := setupAuthService(t)

	existingUser := domain.User{}
	userRepo.On("FindByEmail", mock.Anything, "dupe@example.com").Return(existingUser, nil)

	_, err := svc.Register(context.Background(), "dupe@example.com", "password123", "Test User")
	assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
}

func TestAuthService_Register_InvalidEmail(t *testing.T) {
	svc, userRepo, _ := setupAuthService(t)

	userRepo.On("FindByEmail", mock.Anything, "bad").Return(domain.User{}, domain.ErrUserNotFound)

	_, err := svc.Register(context.Background(), "bad", "password123", "Test User")
	assert.ErrorIs(t, err, domain.ErrInvalidEmail)
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, userRepo, tokenStore := setupAuthService(t)

	hasher := crypto.NewPasswordHasher()
	hash, _ := hasher.Hash("password123")

	user := domain.User{}
	user.ID = uuid.Must(uuid.NewV7())
	user.Email = "login@example.com"
	user.PasswordHash = &hash

	userRepo.On("FindByEmail", mock.Anything, "login@example.com").Return(user, nil)
	tokenStore.On("Save", mock.Anything, mock.Anything, user.ID.String()).Return(nil)

	pair, err := svc.Login(context.Background(), "login@example.com", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	svc, userRepo, _ := setupAuthService(t)

	userRepo.On("FindByEmail", mock.Anything, "bad@example.com").Return(domain.User{}, domain.ErrUserNotFound)

	_, err := svc.Login(context.Background(), "bad@example.com", "wrong")
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestAuthService_Logout_Success(t *testing.T) {
	svc, _, tokenStore := setupAuthService(t)

	tokenStore.On("Delete", mock.Anything, "some-refresh-token").Return(nil)

	err := svc.Logout(context.Background(), "some-refresh-token")
	assert.NoError(t, err)
}

func TestAuthService_Refresh_Success(t *testing.T) {
	svc, userRepo, tokenStore := setupAuthService(t)

	userID := uuid.Must(uuid.NewV7())
	user := domain.User{ID: userID, Email: "refresh@example.com"}

	tokenStore.On("Get", mock.Anything, "valid-refresh").Return(userID.String(), nil)
	tokenStore.On("Delete", mock.Anything, "valid-refresh").Return(nil)
	userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
	tokenStore.On("Save", mock.Anything, mock.Anything, userID.String()).Return(nil)

	pair, err := svc.Refresh(context.Background(), "valid-refresh")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
}

func TestAuthService_Refresh_TokenNotFound(t *testing.T) {
	svc, _, tokenStore := setupAuthService(t)

	tokenStore.On("Get", mock.Anything, "expired-refresh").Return("", domain.ErrTokenNotFound)

	_, err := svc.Refresh(context.Background(), "expired-refresh")
	assert.ErrorIs(t, err, domain.ErrTokenNotFound)
}
