package socket

import (
	"fmt"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewMessagesHandler(args MessageHandlerArgs) bool {
	event := args.message.Event

	switch event {
	case "insert":
		return InsertMessage(args.message, args.conn, args.staticID)
	}

	return false
}

func InsertMessage(message *SocketMessage, conn *websocket.Conn, staticID primitive.ObjectID) bool {
	content := message.Data["content"]

	fmt.Println(content)

	return true
}
