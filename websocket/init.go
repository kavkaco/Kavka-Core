package websocket

import (
	"Tahagram/logs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func InitWebSocket(app *fiber.App) {
	app.Use("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			msgType, msg, msgErr := c.ReadMessage()
			if msgErr != nil {
				logs.ErrorLogger.Println("Error in reading socket message.")
				break
			}

			WebSocketMessageHandler(string(msg), msgType, c)
		}
	}))
}
