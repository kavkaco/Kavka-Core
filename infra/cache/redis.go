package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss for key: %s", key)
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	return json.Unmarshal(data, dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *RedisCache) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func UserCacheKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

func ChatCacheKey(chatID string) string {
	return fmt.Sprintf("chat:%s", chatID)
}

func ChatListCacheKey(userID string) string {
	return fmt.Sprintf("chats:%s", userID)
}

const (
	UserCacheTTL      = 5 * time.Minute
	ChatCacheTTL      = 5 * time.Minute
	ChatListCacheTTL  = 2 * time.Minute
)
