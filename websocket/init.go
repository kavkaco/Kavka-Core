package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func InitWebSocket(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		// authorized, _ := auth.AuthenticateUser(c)

		// if authorized {

		// } else {

		// }

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("accepted", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", ws)
}
