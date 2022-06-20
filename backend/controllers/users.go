package controllers

import (
	"fmt"

	"Tahagram/session"

	"github.com/gofiber/fiber/v2"
)

type User struct {
}

func LoginAction(c *fiber.Ctx) error {
	sess, sessErr := session.SessionStore.Get(c)
	if sessErr != nil {
		// FIXME - fix auto error handlers
		c.Status(500).SendString("Internal Server Error")
	}

	fmt.Println(sess)

	c.Status(200).SendString("Hello World")

	return nil
}
