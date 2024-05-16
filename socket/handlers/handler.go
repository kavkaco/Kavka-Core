package handlers

import (
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/internal/model/message"
	"github.com/kavkaco/Kavka-Core/internal/model/user"
	"github.com/kavkaco/Kavka-Core/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var HandlersList = []func(args HandlerArgs) (ok bool, err error){
	NewChatsHandler,
	NewMessagesHandler,
}

type HandlerServices struct {
	UserService user.Service
	ChatService chat.Service
	MsgService  message.Service
}

type HandlerArgs struct {
	Logger       *zap.Logger
	Adapter      socket.SocketAdapter
	UserStaticID primitive.ObjectID
	Message      socket.IncomingSocketMessage
	Services     *HandlerServices
	Conn         interface{}
}

func NewSocketHandler(logger *zap.Logger, adapter socket.SocketAdapter, conn interface{}, services *HandlerServices, userStaticID primitive.ObjectID) error {
	err := adapter.HandleMessages(conn, func(msg socket.IncomingSocketMessage) {
		// Define HandlerArgs
		handlerArgs := HandlerArgs{
			Conn:         conn,
			Logger:       logger,
			Adapter:      adapter,
			Message:      msg,
			Services:     services,
			UserStaticID: userStaticID,
		}

		// Add Handlers Of HandlersList
		for _, handler := range HandlersList {
			_, err := handler(handlerArgs)
			if err != nil {
				logger.Error("Unhandled event on handlers: " + err.Error())
			}
		}
	})
	if err != nil {
		logger.Error("Unable to call adapter.HandleMessages() in handler.go: " + err.Error())
		return err
	}

	return nil
}
