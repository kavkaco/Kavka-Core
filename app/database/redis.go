package database

import (
	"Nexus/internal/configs"
	"fmt"
	"log"

	"github.com/gofiber/storage/redis"
)

var RedisStore *redis.Storage

func EstablishRedisDBConnection(redisConfigs configs.RedisConfigs) {
	client := redis.New(redis.Config{
		Host:     redisConfigs.Host,
		Username: redisConfigs.Username,
		Password: redisConfigs.Password,
		Database: redisConfigs.Database,
		Port:     redisConfigs.Port,
	})

	if client != nil {
		fmt.Println("Successfully connected to redis database")
		RedisStore = client
	} else {
		log.Fatal("Error in connecting to redis database")
	}
}
