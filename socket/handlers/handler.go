package handlers

import (
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/message"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"github.com/kavkaco/Kavka-Core/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var ErrInvalidHandlerEvent = errors.New("invalid handler event")

type HandlerServices struct {
	UserService user.Service
	ChatService chat.Service
	MsgService  message.Service
}

type HandlerArgs struct {
	Logger       *zap.Logger
	Adapter      socket.SocketAdapter
	UserStaticID primitive.ObjectID
	Message      socket.SocketMessage
	Services     *HandlerServices
}

func NewSocketHandler(logger *zap.Logger, adapter socket.SocketAdapter, conn interface{}, services *HandlerServices) error {
	err := adapter.HandleMessages(conn, func(msg socket.SocketMessage) {
		// Define HandlerArgs
		handlerArgs := HandlerArgs{
			Logger:       logger,
			Adapter:      adapter,
			Message:      msg,
			Services:     services,
			UserStaticID: primitive.NilObjectID,
		}

		// Add Handlers
		_, err := NewChatsHandler(handlerArgs)
		if err != nil {
			logger.Error("Unhandled event on (chats handler): " + err.Error())
		}

		_, err = NewChatsHandler(handlerArgs)
		if err != nil {
			logger.Error("Unhandled event on (chats handler): " + err.Error())
		}
	})
	if err != nil {
		logger.Error("Unable to call adapter.HandleMessages() in handler.go: " + err.Error())
		return err
	}

	return nil
}
