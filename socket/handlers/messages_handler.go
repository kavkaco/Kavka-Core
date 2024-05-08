package handlers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewMessagesHandler(args HandlerArgs) (bool, error) {
	event := args.Message.Event

	switch event {
	case "insert_text_message":
		return InsertTextMessage(event, args)
	case "delete_message":
		return DeleteMessage(event, args)
	}

	return false, ErrInvalidHandlerEvent
}

func InsertTextMessage(_ string, args HandlerArgs) (bool, error) {
	chatID := args.Message.Data["chat_id"]
	messageContent := args.Message.Data["message_content"]

	chatID, err := primitive.ObjectIDFromHex(chatID.(string))
	if err != nil {
		return false, err
	}

	_, err = args.Services.MsgService.InsertTextMessage(chatID.(primitive.ObjectID), args.UserStaticID, messageContent.(string))
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteMessage(_ string, args HandlerArgs) (bool, error) {
	chatID := args.Message.Data["chat_id"]
	messageID := args.Message.Data["message_id"]

	chatID, err := primitive.ObjectIDFromHex(chatID.(string))
	if err != nil {
		return false, err
	}

	messageID, err = primitive.ObjectIDFromHex(messageID.(string))
	if err != nil {
		return false, err
	}

	err = args.Services.MsgService.DeleteMessage(chatID.(primitive.ObjectID), messageID.(primitive.ObjectID))
	if err != nil {
		return false, err
	}

	return true, nil
}
