package cache

import "github.com/go-redis/redis/v8"

type UserCache struct {
	redisClient *redis.Client
}

func NewUserCacheRepository(redisClient *redis.Client) *UserCache {
	return &UserCache{redisClient}
}

// TODO
func (rc *UserCache) SetToken(token string, username string) error { return nil }
func (rc *UserCache) GetToken(token string) (string, error)        { return "", nil }
