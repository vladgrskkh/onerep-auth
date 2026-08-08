package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vladgrskkh/onerep-auth/internal/auth"
	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

type OAuthAccountRepo struct {
	pool *pgxpool.Pool
}

func NewOAuthAccountRepo(pool *pgxpool.Pool) *OAuthAccountRepo {
	return &OAuthAccountRepo{pool: pool}
}

func (r *OAuthAccountRepo) Create(ctx context.Context, account auth.OAuthAccount) (auth.OAuthAccount, error) {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth.oauth_accounts (id, user_id, provider, provider_user_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, account.ID, account.UserID, account.Provider, account.ProviderUserID, account.CreatedAt)
	return account, err
}

func (r *OAuthAccountRepo) FindByProviderID(
	ctx context.Context,
	provider auth.OAuthProvider,
	providerUserID string,
) (auth.OAuthAccount, error) {
	var a auth.OAuthAccount
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, provider, provider_user_id, created_at
		FROM auth.oauth_accounts WHERE provider = $1 AND provider_user_id = $2
	`, provider, providerUserID).Scan(&a.ID, &a.UserID, &a.Provider, &a.ProviderUserID, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.OAuthAccount{}, domain.ErrUserNotFound
	}
	return a, err
}

func (r *OAuthAccountRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]auth.OAuthAccount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, provider, provider_user_id, created_at
		FROM auth.oauth_accounts WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []auth.OAuthAccount
	for rows.Next() {
		var a auth.OAuthAccount
		if err := rows.Scan(&a.ID, &a.UserID, &a.Provider, &a.ProviderUserID, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}
