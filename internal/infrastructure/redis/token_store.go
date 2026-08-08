package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenData struct {
	UserID string `json:"user_id"`
}

type TokenStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTokenStore(client *redis.Client, ttl time.Duration) *TokenStore {
	return &TokenStore{client: client, ttl: ttl}
}

func (s *TokenStore) Save(ctx context.Context, token string, userID string) error {
	data, err := json.Marshal(TokenData{UserID: userID})
	if err != nil {
		return fmt.Errorf("marshal token data: %w", err)
	}
	key := refreshKey(token)
	return s.client.Set(ctx, key, data, s.ttl).Err()
}

func (s *TokenStore) Get(ctx context.Context, token string) (string, error) {
	key := refreshKey(token)
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found")
	}
	if err != nil {
		return "", err
	}
	var td TokenData
	if err := json.Unmarshal(data, &td); err != nil {
		return "", err
	}
	return td.UserID, nil
}

func (s *TokenStore) Delete(ctx context.Context, token string) error {
	key := refreshKey(token)
	return s.client.Del(ctx, key).Err()
}

func refreshKey(token string) string {
	return "refresh_token:" + token
}
