package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	redislib "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	authredis "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/redis"
)

const defaultRedisURL = "redis://localhost:6379/0"

type TokenStoreTestSuite struct {
	suite.Suite

	client *redislib.Client
	store  *authredis.TokenStore
}

func (s *TokenStoreTestSuite) SetupTest() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = defaultRedisURL
	}

	opts, err := redislib.ParseURL(redisURL)
	s.Require().NoError(err)

	s.client = redislib.NewClient(opts)
	s.T().Cleanup(func() {
		if err := s.client.Close(); err != nil {
			s.T().Errorf("close Redis client: %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.client.Ping(ctx).Err(); err != nil {
		s.T().Skipf("Redis unavailable at %s: %v", redisURL, err)
	}

	s.store = authredis.NewTokenStore(s.client, time.Hour)
}

func (s *TokenStoreTestSuite) TestConsumeReturnsStoredTokenOnce() {
	ctx := context.Background()
	token := "auth-review-" + uuid.NewString()
	key := "refresh_token:" + token
	userID := "user-123"

	s.T().Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.client.Del(cleanupCtx, key).Err(); err != nil {
			s.T().Errorf("clean up Redis key %q: %v", key, err)
		}
	})

	s.Require().NoError(s.store.Save(ctx, token, userID))

	consumedUserID, err := s.store.Consume(ctx, token)
	s.Require().NoError(err)
	s.Equal(userID, consumedUserID)

	consumedUserID, err = s.store.Consume(ctx, token)
	s.Empty(consumedUserID)
	s.Require().EqualError(err, "token not found")
}

func TestTokenStoreSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TokenStoreTestSuite))
}
