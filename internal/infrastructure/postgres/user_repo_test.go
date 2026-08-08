//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testPool() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), "postgres://gym:gym_pass@localhost:5432/gym?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return pool
}

func TestUserRepo_CreateAndFindByID(t *testing.T) {
	pool := testPool()
	defer pool.Close()
	repo := postgres.NewUserRepo(pool)
	ctx := context.Background()

	u, err := auth.NewUser("findbyid@test.com", "password123", "Find Test")
	require.NoError(t, err)

	created, err := repo.Create(ctx, u)
	require.NoError(t, err)
	assert.Equal(t, u.ID, created.ID)

	found, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "findbyid@test.com", found.Email)
}

func TestUserRepo_FindByEmail_NotFound(t *testing.T) {
	pool := testPool()
	defer pool.Close()
	repo := postgres.NewUserRepo(pool)
	ctx := context.Background()

	_, err := repo.FindByEmail(ctx, "noone@test.com")
	assert.ErrorIs(t, err, auth.ErrUserNotFound)
}

func TestUserRepo_Update_VersionConflict(t *testing.T) {
	pool := testPool()
	defer pool.Close()
	repo := postgres.NewUserRepo(pool)
	ctx := context.Background()

	u, _ := auth.NewUser("version@test.com", "password123", "Version Test")
	created, _ := repo.Create(ctx, u)

	created.DisplayName = "Updated"
	updated, err := repo.Update(ctx, created)
	require.NoError(t, err)
	assert.Equal(t, 2, updated.Version)

	updated.Version = 1
	_, err = repo.Update(ctx, updated)
	assert.ErrorIs(t, err, auth.ErrUserNotFound)
}
