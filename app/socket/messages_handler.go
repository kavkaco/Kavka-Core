package socket

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewMessagesHandler(args MessageHandlerArgs) bool {
	event := args.message.Event

	switch event {
	case "insert_text_message":
		return InsertTextMessage(event, args)
	case "delete_message":
		return DeleteMessage(event, args)
	}

	return false
}

func InsertTextMessage(event string, args MessageHandlerArgs) bool {
	chatID := args.message.Data["chat_id"]
	messageContent := args.message.Data["message_content"]

	chatID, parseErr := primitive.ObjectIDFromHex(chatID.(string))
	if parseErr != nil {
		return false
	}

	args.socketService.msgService.InsertTextMessage(chatID.(primitive.ObjectID), args.staticID, messageContent.(string))

	return true
}

func DeleteMessage(event string, args MessageHandlerArgs) bool {
	chatID := args.message.Data["chat_id"]
	messageID := args.message.Data["message_id"]

	chatID, parseErr := primitive.ObjectIDFromHex(chatID.(string))
	if parseErr != nil {
		return false
	}

	messageID, parseErr = primitive.ObjectIDFromHex(messageID.(string))
	if parseErr != nil {
		return false
	}

	args.socketService.msgService.DeleteMessage(chatID.(primitive.ObjectID), messageID.(primitive.ObjectID))

	return true
}
