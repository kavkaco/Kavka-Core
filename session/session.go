package session

import (
	"Tahagram/configs"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

func InitSession() *session.Store {
	wd, _ := os.Getwd()

	redisConfigs, redisConfigsErr := configs.ParseRedisConfigs(wd + "/configs/redis.yml")
	if redisConfigsErr != nil {
		log.Fatal("Error in connecting to redis database")
	}

	var redisStore = redis.New(redis.Config{
		Host:     redisConfigs.Host,
		Username: redisConfigs.Username,
		Password: redisConfigs.Password,
		Database: redisConfigs.Database,
		Port:     redisConfigs.Port,
	})

	var session = session.New(session.Config{
		KeyLookup:  "cookie:session_id",
		Storage:    redisStore,
		Expiration: 10 * 24 * time.Hour, // 14 Days
	})

	fmt.Println("Connected To Redis Database")

	return session
}

var SessionStore session.Store = *InitSession()
