package websocket

import (
	"Kavka/app/httpstatus"
	"Kavka/pkg/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func InitWebSocket(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("accepted", true)

			authorized, _ := auth.AuthenticateUser(c)

			if authorized {
				return c.Next()
			} else {
				httpstatus.Unauthorized(c)
			}
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", ws)
}
