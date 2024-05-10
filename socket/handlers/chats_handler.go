package handlers

import (
	"github.com/kavkaco/Kavka-Core/app/presenters"
	"github.com/kavkaco/Kavka-Core/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewChatsHandler(args HandlerArgs) (ok bool, err error) {
	event := args.Message.Event

	switch event {
	case "get_chat":
		return GetChat(event, args)
	case "create_direct":
		return CreateDirect(event, args)
	case "create_group":
		return CreateGroup(event, args)
	case "create_channel":
		return CreateChannel(event, args)
	}

	return false, ErrInvalidHandlerEvent
}

func CreateDirect(event string, args HandlerArgs) (bool, error) {
	staticID := args.Message.Data["static_id"]

	staticID, err := primitive.ObjectIDFromHex(staticID.(string))
	if err != nil {
		return false, err
	}

	_, err = args.Services.ChatService.CreateDirect(args.UserStaticID, staticID.(primitive.ObjectID))
	if err != nil {
		return false, err
	}

	// FIXME
	// err = args.conn.WriteJSON(presenters.ChatAsJSON(event, chat))

	return true, nil
}

func GetChat(event string, args HandlerArgs) (bool, error) {
	staticID := args.Message.Data["static_id"]

	staticID, err := primitive.ObjectIDFromHex(staticID.(string))
	if err != nil {
		return false, err
	}

	foundChat, err := args.Services.ChatService.GetChat(staticID.(primitive.ObjectID))
	if err != nil {
		return false, err
	}

	chatJson, err := presenters.ChatAsJSON(*foundChat, args.UserStaticID)
	if err != nil {
		return false, err
	}

	err = args.Adapter.WriteMessage(args.Conn, socket.OutgoingSocketMessage{
		Status: 200,
		Event:  "chat_found",
		Data:   chatJson,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateGroup(event string, args HandlerArgs) (bool, error) {
	title := args.Message.Data["title"]
	username := args.Message.Data["username"]
	description := args.Message.Data["description"]

	if title != nil && username != nil && description != nil {
		_, err := args.Services.ChatService.CreateGroup(args.UserStaticID, title.(string), username.(string), description.(string))
		if err != nil {
			return false, err
		}

		// FIXME
		// err := args.conn.WriteJSON(presenters.ChatAsJSON(event, chat))

		return true, nil
	}

	return false, nil
}

func CreateChannel(event string, args HandlerArgs) (bool, error) {
	title := args.Message.Data["title"]
	username := args.Message.Data["username"]
	description := args.Message.Data["description"]

	if title != nil && username != nil && description != nil {
		_, err := args.Services.ChatService.CreateChannel(args.UserStaticID, title.(string), username.(string), description.(string))
		if err != nil {
			return false, err
		}

		// FIXME
		// err := args.conn.WriteJSON(presenters.ChatAsJSON(event, chat))

		return true, nil
	}

	return false, nil
}
