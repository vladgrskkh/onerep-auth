package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth.users (id, email, password_hash, display_name, avatar_url, gender, birth_date, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, user.ID, user.Email, user.PasswordHash, user.DisplayName, user.AvatarURL, user.Gender, user.BirthDate, user.CreatedAt, user.UpdatedAt, user.Version)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, avatar_url, gender, birth_date, created_at, updated_at, version
		FROM auth.users WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Gender, &u.BirthDate, &u.CreatedAt, &u.UpdatedAt, &u.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	return u, err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, avatar_url, gender, birth_date, created_at, updated_at, version
		FROM auth.users WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Gender, &u.BirthDate, &u.CreatedAt, &u.UpdatedAt, &u.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	return u, err
}

func (r *UserRepo) Update(ctx context.Context, user domain.User) (domain.User, error) {
	user.Version++
	tag, err := r.pool.Exec(ctx, `
		UPDATE auth.users SET
			display_name = $1, avatar_url = $2, gender = $3, birth_date = $4,
			updated_at = $5, version = $6
		WHERE id = $7 AND version = $8
	`, user.DisplayName, user.AvatarURL, user.Gender, user.BirthDate, user.UpdatedAt, user.Version, user.ID, user.Version-1)
	if err != nil {
		return domain.User{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}
