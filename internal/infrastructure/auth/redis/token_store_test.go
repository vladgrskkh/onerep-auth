package redis_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	redislib "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	authredis "github.com/vladgrskkh/onerep-auth/internal/infrastructure/auth/redis"
)

type TokenStoreTestSuite struct {
	suite.Suite

	store  *authredis.TokenStore
	client *observingRedisClient
}

func (s *TokenStoreTestSuite) SetupTest() {
	s.client = newObservingRedisClient()
	s.store = authredis.NewTokenStore(s.client, time.Hour)
}

func (s *TokenStoreTestSuite) TestConsumeReturnsStoredTokenOnce() {
	ctx := context.Background()
	token := "refresh-token"
	userID := "user-123"
	s.Require().NoError(s.store.Save(ctx, token, userID))

	consumedUserID, err := s.store.Consume(ctx, token)
	s.Require().NoError(err)
	s.Equal(userID, consumedUserID)

	consumedUserID, err = s.store.Consume(ctx, token)
	s.Empty(consumedUserID)
	s.Require().EqualError(err, "token not found")
	s.Equal([]string{"refresh_token:" + token, "refresh_token:" + token}, s.client.getDelCalls)
}

type observingRedisClient struct {
	mu          sync.Mutex
	values      map[string]string
	getDelCalls []string
}

func newObservingRedisClient() *observingRedisClient {
	return &observingRedisClient{values: make(map[string]string)}
}

func (c *observingRedisClient) Set(
	ctx context.Context,
	key string,
	value any,
	_ time.Duration,
) *redislib.StatusCmd {
	cmd := redislib.NewStatusCmd(ctx)

	var storedValue string
	switch value := value.(type) {
	case []byte:
		storedValue = string(value)
	case string:
		storedValue = value
	default:
		cmd.SetErr(fmt.Errorf("unsupported SET value type %T", value))
		return cmd
	}

	c.mu.Lock()
	c.values[key] = storedValue
	c.mu.Unlock()
	cmd.SetVal("OK")
	return cmd
}

func (c *observingRedisClient) GetDel(ctx context.Context, key string) *redislib.StringCmd {
	cmd := redislib.NewStringCmd(ctx)

	c.mu.Lock()
	c.getDelCalls = append(c.getDelCalls, key)
	value, ok := c.values[key]
	delete(c.values, key)
	c.mu.Unlock()

	if !ok {
		cmd.SetErr(redislib.Nil)
		return cmd
	}
	cmd.SetVal(value)
	return cmd
}

func (c *observingRedisClient) Del(ctx context.Context, keys ...string) *redislib.IntCmd {
	cmd := redislib.NewIntCmd(ctx)

	c.mu.Lock()
	var deleted int64
	for _, key := range keys {
		if _, ok := c.values[key]; ok {
			delete(c.values, key)
			deleted++
		}
	}
	c.mu.Unlock()
	cmd.SetVal(deleted)
	return cmd
}

func TestTokenStoreSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TokenStoreTestSuite))
}
