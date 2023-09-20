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

func InsertTextMessage(_ string, args MessageHandlerArgs) bool {
	chatID := args.message.Data["chat_id"]
	messageContent := args.message.Data["message_content"]

	chatID, parseErr := primitive.ObjectIDFromHex(chatID.(string))
	if parseErr != nil {
		return false
	}

	_, err := args.socketService.msgService.InsertTextMessage(chatID.(primitive.ObjectID), args.staticID, messageContent.(string))

	return err == nil
}

func DeleteMessage(_ string, args MessageHandlerArgs) bool {
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

	err := args.socketService.msgService.DeleteMessage(chatID.(primitive.ObjectID), messageID.(primitive.ObjectID))

	return err == nil
}
