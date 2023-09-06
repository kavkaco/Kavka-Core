package socket

import (
	"fmt"

	"github.com/fasthttp/websocket"
)

func NewMessagesHandler(args MessageHandlerArgs) bool {
	event := args.message.Event

	switch event {
	case "insert":
		return InsertMessage(args.message, args.conn.Conn, args.staticID)
	}

	return false
}

func InsertMessage(message *SocketMessage, conn *websocket.Conn, staticID string) bool {
	content := message.Data["Content"]
	fmt.Println(content.(string))

	return true
}
