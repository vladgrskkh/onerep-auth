package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	jwtsvc "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/jwt"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/oauth"
)

type OAuthAccountRepository interface {
	Create(
		ctx context.Context,
		account authdomain.OAuthAccount,
	) (authdomain.OAuthAccount, error)
	FindByProviderID(
		ctx context.Context,
		provider authdomain.OAuthProvider,
		providerUserID string,
	) (authdomain.OAuthAccount, error)
}

type OAuthService struct {
	users         UserRepository
	oauthAccounts OAuthAccountRepository
	tokens        jwtsvc.TokenManager
	tokenStore    TokenStorer
	googleAdapter *oauth.GoogleAdapter
}

func NewOAuthService(
	users UserRepository,
	oauthAccounts OAuthAccountRepository,
	tokens jwtsvc.TokenManager,
	tokenStore TokenStorer,
	googleAdapter *oauth.GoogleAdapter,
) *OAuthService {
	return &OAuthService{
		users:         users,
		oauthAccounts: oauthAccounts,
		tokens:        tokens,
		tokenStore:    tokenStore,
		googleAdapter: googleAdapter,
	}
}

func (s *OAuthService) Login(ctx context.Context, provider authdomain.OAuthProvider, code string) (TokenPair, error) {
	switch provider {
	case authdomain.OAuthGoogle:
		return s.googleLogin(ctx, code)
	case authdomain.OAuthApple:
		return TokenPair{}, fmt.Errorf("apple oauth not yet implemented")
	default:
		return TokenPair{}, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *OAuthService) googleLogin(ctx context.Context, code string) (TokenPair, error) {
	gu, err := s.googleAdapter.Exchange(ctx, code)
	if err != nil {
		return TokenPair{}, fmt.Errorf("oauth exchange: %w", err)
	}

	acc, err := s.oauthAccounts.FindByProviderID(ctx, authdomain.OAuthGoogle, gu.ID)
	if err == nil {
		return s.issuePairForUser(ctx, acc.UserID)
	}

	user := authdomain.UserFromOAuth(gu.Email, gu.Name)
	user.AvatarURL = &gu.Picture

	createdUser, err := s.users.Create(ctx, *user)
	if err != nil {
		createdUser, err = s.users.FindByEmail(ctx, gu.Email)
		if err != nil {
			return TokenPair{}, err
		}
	}

	oa := authdomain.NewOAuthAccount(createdUser.ID, authdomain.OAuthGoogle, gu.ID)
	if _, err := s.oauthAccounts.Create(ctx, oa); err != nil {
		return TokenPair{}, err
	}

	return s.issuePairForUser(ctx, createdUser.ID)
}

func (s *OAuthService) issuePairForUser(ctx context.Context, userID uuid.UUID) (TokenPair, error) {
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
