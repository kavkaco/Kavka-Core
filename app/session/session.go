package session

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

var SessionStore *session.Store

func InitSession(redisClient *redis.Storage) {
	session := session.New(session.Config{
		KeyLookup: "cookie:session_id",
		Storage:   redisClient,
	})

	SessionStore = session
}
