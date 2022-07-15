package session

import (
	"Nexus/app/database"
	"log"

	"github.com/gofiber/fiber/v2/middleware/session"
)

func InitSession() *session.Store {
	if database.RedisStore == nil {
		log.Fatal("Error in initializing Session. RedisStore is empty!")
		return nil
	}

	session := session.New(session.Config{
		KeyLookup: "cookie:session_id",
		Storage:   database.RedisStore,
	})

	return session
}

var SessionStore *session.Store
