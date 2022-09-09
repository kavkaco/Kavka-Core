package cache

import (
	"Kavka/config"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type UserCache_TokenData struct {
	Username string `json:"username"`
}

type UserCache struct {
	redisClient *redis.Client
	authConfigs config.Auth
}

func NewUserCacheRepository(redisClient *redis.Client, authConfigs config.Auth) *UserCache {
	return &UserCache{redisClient, authConfigs}
}

func (userCache *UserCache) SetRefreshToken(token string, username string) error {
	tokenData, marshalErr := json.Marshal(UserCache_TokenData{
		Username: username,
	})
	if marshalErr != nil {
		return marshalErr
	}

	setErr := userCache.redisClient.Set(context.Background(), token, tokenData, userCache.authConfigs.RF_TTL_MINUTE).Err()
	if setErr != nil {
		return setErr
	}

	return nil
}

func (userCache *UserCache) GetRefreshToken(token string) (*UserCache_TokenData, error) {
	tokenDataString, getErr := userCache.redisClient.Get(context.Background(), token).Result()
	if getErr != nil {
		return nil, getErr
	}

	var tokenData UserCache_TokenData
	unmarshalErr := json.Unmarshal([]byte(tokenDataString), &tokenData)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &tokenData, nil
}
