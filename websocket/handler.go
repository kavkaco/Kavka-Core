package websocket

import (
	"Tahagram/logs"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

type SSendMessage struct {
	Text string `json:"text"`
}

func WebSocketMessageHandler(msg string, msgType int, c *websocket.Conn) {
	var s *SSendMessage = &SSendMessage{}
	err := c.ReadJSON(s)
	if err != nil {
		logs.ErrorLogger.Println("Error in parsing websocket message to JSON: " + err.Error())
	}

	fmt.Println(s.Text)
}
