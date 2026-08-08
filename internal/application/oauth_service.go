package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/oauth"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/jwt"
)

type oauthAccountCreator interface {
	Create(ctx context.Context, account domain.OAuthAccount) (domain.OAuthAccount, error)
	FindByProviderID(ctx context.Context, provider domain.OAuthProvider, providerUserID string) (domain.OAuthAccount, error)
}

type OAuthService struct {
	userCreator    userCreator
	userFinder     userFinder
	userByIDFinder userByIDFinder
	oauthAccRepo   oauthAccountCreator
	jwtSvc         jwtsvc.Service
	tokenStore     tokenStorer
	googleAdapter  *oauth.GoogleAdapter
}

func NewOAuthService(
	uc userCreator,
	uf userFinder,
	ubi userByIDFinder,
	oa oauthAccountCreator,
	jwt jwtsvc.Service,
	ts tokenStorer,
	ga *oauth.GoogleAdapter,
) *OAuthService {
	return &OAuthService{
		userCreator:    uc,
		userFinder:     uf,
		userByIDFinder: ubi,
		oauthAccRepo:   oa,
		jwtSvc:         jwt,
		tokenStore:     ts,
		googleAdapter:  ga,
	}
}

func (s *OAuthService) Login(ctx context.Context, provider domain.OAuthProvider, code string) (TokenPair, error) {
	switch provider {
	case domain.OAuthGoogle:
		return s.googleLogin(ctx, code)
	default:
		return TokenPair{}, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *OAuthService) googleLogin(ctx context.Context, code string) (TokenPair, error) {
	gu, err := s.googleAdapter.Exchange(ctx, code)
	if err != nil {
		return TokenPair{}, fmt.Errorf("oauth exchange: %w", err)
	}

	acc, err := s.oauthAccRepo.FindByProviderID(ctx, domain.OAuthGoogle, gu.ID)
	if err == nil {
		return s.issuePairForUser(ctx, acc.UserID)
	}

	user := domain.UserFromOAuth(gu.Email, gu.Name)
	user.AvatarURL = &gu.Picture

	createdUser, err := s.userCreator.Create(ctx, *user)
	if err != nil {
		createdUser, err = s.userFinder.FindByEmail(ctx, gu.Email)
		if err != nil {
			return TokenPair{}, err
		}
	}

	oa := domain.NewOAuthAccount(createdUser.ID, domain.OAuthGoogle, gu.ID)
	if _, err := s.oauthAccRepo.Create(ctx, oa); err != nil {
		return TokenPair{}, err
	}

	return s.issuePairForUser(ctx, createdUser.ID)
}

func (s *OAuthService) issuePairForUser(ctx context.Context, userID uuid.UUID) (TokenPair, error) {
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
