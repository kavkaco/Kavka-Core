package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Username string `json:"username"`
}

func LoginAction(c *fiber.Ctx) error {
	// sess, sessErr := session.SessionStore.Get(c)
	// if sessErr != nil {
	// 	// FIXME - fix auto error handlers
	// 	c.Status(500).SendString("Internal Server Error")
	// }

	// fmt.Println(sess)

	c.SendString("Hello world")
	return nil
}
