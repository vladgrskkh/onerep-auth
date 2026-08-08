package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/crypto"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
)

type userCreator interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

type userFinder interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type userByIDFinder interface {
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type tokenStorer interface {
	Save(ctx context.Context, token string, userID string) error
	Get(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}

type TokenPair = jwtsvc.TokenPair

type AuthService struct {
	userCreator    userCreator
	userFinder     userFinder
	userByIDFinder userByIDFinder
	hasher         crypto.PasswordHasher
	jwtSvc         jwtsvc.Service
	tokenStore     tokenStorer
}

func NewAuthService(
	userCreator userCreator,
	userFinder userFinder,
	userByIDFinder userByIDFinder,
	hasher crypto.PasswordHasher,
	jwtSvc jwtsvc.Service,
	tokenStore tokenStorer,
) *AuthService {
	return &AuthService{
		userCreator:    userCreator,
		userFinder:     userFinder,
		userByIDFinder: userByIDFinder,
		hasher:         hasher,
		jwtSvc:         jwtSvc,
		tokenStore:     tokenStore,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (TokenPair, error) {
	_, err := s.userFinder.FindByEmail(ctx, email)
	if err == nil {
		return TokenPair{}, domain.ErrEmailAlreadyExists
	}

	user, err := domain.NewUser(email, password, displayName)
	if err != nil {
		return TokenPair{}, err
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return TokenPair{}, err
	}
	user.PasswordHash = &hash

	user, err = s.userCreator.Create(ctx, user)
	if err != nil {
		return TokenPair{}, err
	}

	pair, err := s.jwtSvc.IssueTokenPair(user)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.tokenStore.Save(ctx, pair.RefreshToken, user.ID.String()); err != nil {
		return TokenPair{}, err
	}

	return pair, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (TokenPair, error) {
	user, err := s.userFinder.FindByEmail(ctx, email)
	if err != nil {
		return TokenPair{}, domain.ErrInvalidCredentials
	}

	if user.PasswordHash == nil {
		return TokenPair{}, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(*user.PasswordHash, password); err != nil {
		return TokenPair{}, domain.ErrInvalidCredentials
	}

	pair, err := s.jwtSvc.IssueTokenPair(user)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.tokenStore.Save(ctx, pair.RefreshToken, user.ID.String()); err != nil {
		return TokenPair{}, err
	}

	return pair, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokenStore.Delete(ctx, refreshToken)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	userIDStr, err := s.tokenStore.Get(ctx, refreshToken)
	if err != nil {
		return TokenPair{}, domain.ErrTokenNotFound
	}

	if err := s.tokenStore.Delete(ctx, refreshToken); err != nil {
		return TokenPair{}, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return TokenPair{}, err
	}

	user, err := s.userByIDFinder.FindByID(ctx, userID)
	if err != nil {
		return TokenPair{}, err
	}

	pair, err := s.jwtSvc.IssueTokenPair(user)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.tokenStore.Save(ctx, pair.RefreshToken, user.ID.String()); err != nil {
		return TokenPair{}, err
	}

	return pair, nil
}
