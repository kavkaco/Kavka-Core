package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/ory/dockertest/v3"
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

func GetRedisTestInstance(callback func(redisClient *redis.Client)) {
	var dockerContainerEnvVariables = []string{}

	err := os.Setenv("ENV", "test")
	if err != nil {
		log.Fatalf("Could not set the environment variable to test: %s", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	var client *redis.Client

	resource, err := pool.Run("redis", "alpine", dockerContainerEnvVariables)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Kill the container
	defer func() {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	err = pool.Retry(func() error {
		ipAddr := resource.Container.NetworkSettings.IPAddress + ":6379"

		fmt.Printf("Docker redis container network ip address: %s\n", ipAddr)

		client = redis.NewClient(&redis.Options{
			Addr: ipAddr,
			DB:   0,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Could not connect to Redis: %s", err)
	}

	callback(client)
}
