package auth

import (
	"context"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/crypto"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"
)

type UserRepository interface {
	Create(ctx context.Context, user authdomain.User) (authdomain.User, error)
	FindByEmail(ctx context.Context, email string) (authdomain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (authdomain.User, error)
}

type TokenStorer interface {
	Save(ctx context.Context, token string, userID string) error
	Get(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}

type TokenPair = jwtsvc.TokenPair

type AuthService struct {
	users      UserRepository
	hasher     crypto.PasswordHasher
	tokens     jwtsvc.TokenManager
	tokenStore TokenStorer
}

func NewAuthService(
	users UserRepository,
	hasher crypto.PasswordHasher,
	tokens jwtsvc.TokenManager,
	tokenStore TokenStorer,
) *AuthService {
	return &AuthService{
		users:      users,
		hasher:     hasher,
		tokens:     tokens,
		tokenStore: tokenStore,
	}
}

func (s *AuthService) Register(ctx context.Context, cmd RegisterCommand) (TokenPair, error) {
	_, err := s.users.FindByEmail(ctx, cmd.Email)
	if err == nil {
		return TokenPair{}, authdomain.ErrEmailAlreadyExists
	}

	user, err := authdomain.NewUser(cmd.Email, cmd.Password, cmd.DisplayName)
	if err != nil {
		return TokenPair{}, err
	}

	hash, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return TokenPair{}, err
	}
	user.PasswordHash = &hash

	user, err = s.users.Create(ctx, user)
	if err != nil {
		return TokenPair{}, err
	}

	pair, err := s.tokens.IssueTokenPair(user.ID, user.Email)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.tokenStore.Save(ctx, pair.RefreshToken, user.ID.String()); err != nil {
		return TokenPair{}, err
	}

	return pair, nil
}

func (s *AuthService) Login(ctx context.Context, cmd LoginCommand) (TokenPair, error) {
	user, err := s.users.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return TokenPair{}, authdomain.ErrInvalidCredentials
	}

	if user.PasswordHash == nil {
		return TokenPair{}, authdomain.ErrInvalidCredentials
	}

	if cmpErr := s.hasher.Compare(*user.PasswordHash, cmd.Password); cmpErr != nil {
		return TokenPair{}, authdomain.ErrInvalidCredentials
	}

	pair, err := s.tokens.IssueTokenPair(user.ID, user.Email)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.tokenStore.Save(ctx, pair.RefreshToken, user.ID.String()); err != nil {
		return TokenPair{}, err
	}

	return pair, nil
}

func (s *AuthService) Logout(ctx context.Context, cmd LogoutCommand) error {
	return s.tokenStore.Delete(ctx, cmd.RefreshToken)
}

func (s *AuthService) Refresh(ctx context.Context, cmd RefreshCommand) (TokenPair, error) {
	userIDStr, err := s.tokenStore.Get(ctx, cmd.RefreshToken)
	if err != nil {
		return TokenPair{}, authdomain.ErrTokenNotFound
	}

	if delErr := s.tokenStore.Delete(ctx, cmd.RefreshToken); delErr != nil {
		return TokenPair{}, delErr
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return TokenPair{}, err
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return TokenPair{}, err
	}

	pair, err := s.tokens.IssueTokenPair(user.ID, user.Email)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.tokenStore.Save(ctx, pair.RefreshToken, user.ID.String()); err != nil {
		return TokenPair{}, err
	}

	return pair, nil
}
