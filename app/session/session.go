package session

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

func InitSession(redisClient *redis.Storage) *session.Store {
	session := session.New(session.Config{
		KeyLookup: "cookie:session_id",
		Storage:   redisClient,
	})

	return session
}

var SessionStore *session.Store
