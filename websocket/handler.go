package websocket

import (
	"Tahagram/logs"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

type SSendMessage struct {
	Text string `json:"text"`
}

var ws = websocket.New(func(c *websocket.Conn) {
	fmt.Println("user is upgraded")
	fmt.Println(c.Params("id"))
	fmt.Println(c.Locals("accepted"))

	for {
		msgType, msg, msgErr := c.ReadMessage()
		if msgErr != nil {
			logs.ErrorLogger.Println("Error in reading socket message.")
			break
		}

		WebSocketMessageHandler(string(msg), msgType, c)
	}
})

func WebSocketMessageHandler(msg string, msgType int, c *websocket.Conn) {
	var s *SSendMessage = &SSendMessage{}
	err := c.ReadJSON(s)
	if err != nil {
		logs.ErrorLogger.Println("Error in parsing websocket message to JSON: " + err.Error())
	}

	fmt.Println(s.Text)
}
