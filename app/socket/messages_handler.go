package socket

import (
	"fmt"

	"github.com/fasthttp/websocket"
)

func NewMessagesHandler(message *SocketMessage, conn *websocket.Conn, staticID string) {
	event := message.Event

	switch event {
	case "insert":
		InsertMessage(message, conn, staticID)
	}
}

func InsertMessage(message *SocketMessage, conn *websocket.Conn, staticID string) {
	content := message.Data["Content"]
	fmt.Println(content.(string))
}
