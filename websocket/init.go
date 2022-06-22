package websocket

import (
	"Tahagram/logs"
	"encoding/json"

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

func SendJson(inter interface{}, c *websocket.Conn, msgType int) bool {
	data, marshalErr := json.Marshal(inter)
	if marshalErr != nil {
		return false
	}

	writeErr := c.WriteMessage(msgType, data)

	if writeErr != nil {
		return false
	}

	return true
}

func ReadJson(data string) map[string]interface{} {
	// FIXME
	var parsedData interface{}

	parseErr := json.Unmarshal([]byte(data), &parsedData)

	if parseErr != nil {
		return nil
	}

	return parsedData.(map[string]interface{})
}
