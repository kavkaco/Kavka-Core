package socket

import (
	"fmt"

	"github.com/fasthttp/websocket"
)

func NewChatsHandler(message *SocketMessage, conn *websocket.Conn, staticID string) {
	event := message.Event

	switch event {
	case "NewChat":
		NewChat(message, conn, staticID)
	}
}

func NewChat(message *SocketMessage, conn *websocket.Conn, staticID string) {
	content := message.Data["Content"]
	fmt.Println(content.(string))
}
