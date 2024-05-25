package handlers

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	chat "github.com/kavkaco/Kavka-Core/internal/service/chat"
	message "github.com/kavkaco/Kavka-Core/internal/service/message"
	user "github.com/kavkaco/Kavka-Core/internal/service/user"
	"github.com/kavkaco/Kavka-Core/socket"
	"go.uber.org/zap"
)

var HandlersList = []func(args HandlerArgs) (ok bool, err error){
	NewChatsHandler,
	NewMessagesHandler,
}

type HandlerServices struct {
	UserService    user.UserService
	ChatService    chat.ChatService
	MessageService message.MessageService
}

type HandlerArgs struct {
	Ctx      context.Context
	Logger   *zap.Logger
	Adapter  socket.SocketAdapter
	UserID   model.UserID
	Message  socket.IncomingSocketMessage
	Services *HandlerServices
	Conn     interface{}
}

func NewSocketHandler(ctx context.Context, logger *zap.Logger, adapter socket.SocketAdapter, conn interface{}, services *HandlerServices, userID model.UserID) error {
	err := adapter.HandleMessages(conn, func(msg socket.IncomingSocketMessage) {
		// Define HandlerArgs
		handlerArgs := HandlerArgs{
			Ctx:      ctx,
			Conn:     conn,
			Logger:   logger,
			Adapter:  adapter,
			Message:  msg,
			Services: services,
			UserID:   userID,
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
