package database

import (
	"Kavka/internal/configs"
	"fmt"
	"log"

	"github.com/gofiber/storage/redis"
)

func InitRedisDB(redisConfigs configs.Redis) *redis.Storage {
	client := redis.New(redis.Config{
		Host:     redisConfigs.Host,
		Username: redisConfigs.Username,
		Password: redisConfigs.Password,
		Database: redisConfigs.Database,
		Port:     redisConfigs.Port,
	})

	if client != nil {
		fmt.Println("Successfully connected to redis database")
		return client
	} else {
		log.Fatal("Error in connecting to redis database")
		return nil
	}
}
