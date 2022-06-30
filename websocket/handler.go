package websocket

import (
	"Tahagram/logs"
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

type MessageData struct {
	Text string `json:"message"`
}

type SendMessageData struct {
	Text string `json:"text"`
}

var ws = websocket.New(func(c *websocket.Conn) {
	for {
		_, msg, msgErr := c.ReadMessage()
		if msgErr != nil {
			logs.ErrorLogger.Println("Error in reading socket message.")
			break
		}

		var data *MessageData = &MessageData{}
		parseErr := json.Unmarshal([]byte(msg), &data)
		if parseErr != nil {
			logs.ErrorLogger.Println("Error in parsing websocket message to JSON: " + parseErr.Error())
		} else {
			WebSocketMessageHandler(data, c)
		}
	}
})

func WebSocketMessageHandler(data *MessageData, c *websocket.Conn) {
	c.WriteJSON(data)
	fmt.Println(data)
}
