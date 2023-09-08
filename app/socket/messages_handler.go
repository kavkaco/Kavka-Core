package socket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func NewMessagesHandler(args MessageHandlerArgs) bool {
	event := args.message.Event

	log.Println(event)

	switch event {
	case "insert":
		return InsertMessage(args.message, args.conn, args.staticID)
	}

	return false
}

func InsertMessage(message *SocketMessage, conn *websocket.Conn, staticID string) bool {
	content := message.Data["content"]

	fmt.Println(content)

	return true
}
