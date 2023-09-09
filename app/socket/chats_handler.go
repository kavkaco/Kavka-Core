package socket

import (
	"Kavka/app/presenters"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewChatsHandler(args MessageHandlerArgs) bool {
	event := args.message.Event

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

	return false
}
func CreateDirect(event string, args MessageHandlerArgs) bool {
	staticID := args.message.Data["static_id"]

	staticID, parseErr := primitive.ObjectIDFromHex(staticID.(string))
	if parseErr != nil {
		return false
	}

	chat, err := args.socketService.chatService.CreateDirect(args.staticID, staticID.(primitive.ObjectID))
	if err != nil {
		return false
	}

	args.conn.WriteJSON(presenters.ChatAsJSON(event, chat))

	return true
}

func GetChat(event string, args MessageHandlerArgs) bool {
	staticID := args.message.Data["static_id"]

	staticID, parseErr := primitive.ObjectIDFromHex(staticID.(string))
	if parseErr != nil {
		return false
	}

	chat, err := args.socketService.chatService.GetChat(staticID.(primitive.ObjectID))
	if err != nil {
		log.Println("find chat error in socket:", err)
		return false
	}

	args.conn.WriteJSON(presenters.ChatAsJSON(event, chat))

	return true
}

func CreateGroup(event string, args MessageHandlerArgs) bool {
	title := args.message.Data["title"]
	username := args.message.Data["username"]
	description := args.message.Data["description"]

	if title != nil && username != nil && description != nil {
		chat, err := args.socketService.chatService.CreateGroup(args.staticID, title.(string), username.(string), description.(string))
		if err != nil {
			return false
		}

		args.conn.WriteJSON(presenters.ChatAsJSON(event, chat))

		return true
	}

	return false
}

func CreateChannel(event string, args MessageHandlerArgs) bool {
	title := args.message.Data["title"]
	username := args.message.Data["username"]
	description := args.message.Data["description"]

	if title != nil && username != nil && description != nil {
		chat, err := args.socketService.chatService.CreateChannel(args.staticID, title.(string), username.(string), description.(string))
		if err != nil {
			return false
		}

		args.conn.WriteJSON(presenters.ChatAsJSON(event, chat))

		return true
	}

	return false
}
