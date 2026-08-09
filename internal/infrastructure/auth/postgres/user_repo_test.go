//go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/postgres"
)

type UserRepoTestSuite struct {
	suite.Suite

	pool *pgxpool.Pool
	repo *postgres.UserRepo
	ctx  context.Context
}

func (s *UserRepoTestSuite) SetupTest() {
	pool, err := pgxpool.New(context.Background(), "postgres://gym:gym_pass@localhost:5432/gym?sslmode=disable")
	if err != nil {
		panic(err)
	}
	s.pool = pool
	s.repo = postgres.NewUserRepo(pool)
	s.ctx = context.Background()
}

func (s *UserRepoTestSuite) TearDownTest() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *UserRepoTestSuite) TestCreateAndFindByID() {
	u, err := auth.NewUser("findbyid@test.com", "password123", "Find Test")
	s.Require().NoError(err)

	created, err := s.repo.Create(s.ctx, u)
	s.Require().NoError(err)
	s.Equal(u.ID, created.ID)

	found, err := s.repo.FindByID(s.ctx, u.ID)
	s.Require().NoError(err)
	s.Equal("findbyid@test.com", found.Email)
}

func (s *UserRepoTestSuite) TestFindByEmail_NotFound() {
	_, err := s.repo.FindByEmail(s.ctx, "noone@test.com")
	s.ErrorIs(err, auth.ErrUserNotFound)
}

func (s *UserRepoTestSuite) TestUpdate_VersionConflict() {
	u, _ := auth.NewUser("version@test.com", "password123", "Version Test")
	created, _ := s.repo.Create(s.ctx, u)

	created.DisplayName = "Updated"
	updated, err := s.repo.Update(s.ctx, created)
	s.Require().NoError(err)
	s.Equal(2, updated.Version)

	updated.Version = 1
	_, err = s.repo.Update(s.ctx, updated)
	s.ErrorIs(err, auth.ErrUserNotFound)
}

func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
