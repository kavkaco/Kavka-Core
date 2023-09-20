package database

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/kavkaco/Kavka-Core/config"
)

var (
	redisLock     = &sync.Mutex{}
	redisInstance *redis.Client
)

func GetRedisDBInstance(redisConfigs config.Redis) *redis.Client {
	if redisInstance == nil {
		redisLock.Lock()
		defer redisLock.Unlock()

		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConfigs.Host, redisConfigs.Port),
			Password: redisConfigs.Password,
			DB:       redisConfigs.DB,
		})

		redisInstance = client
	}

	return redisInstance
}
