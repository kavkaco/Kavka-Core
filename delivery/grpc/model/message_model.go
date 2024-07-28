package grpc_model

import (
	"github.com/kavkaco/Kavka-Core/internal/model"
	modelv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/model/message/v1"
	"github.com/kavkaco/Kavka-Core/utils"
)

func TransformMessageToGrpcModel(message *model.Message) *modelv1.Message {
	transformedModel := &modelv1.Message{
		MessageId: message.MessageID.Hex(),
		SenderId:  message.SenderID,
		CreatedAt: message.CreatedAt.Unix(),
		Edited:    message.Edited,
		Seen:      message.Seen,
		Type:      message.Type,
	}

	switch transformedModel.Type {
	case "text":
		messageContent, _ := utils.TypeConverter[model.TextMessage](message.Content)

		transformedModel.Payload = &modelv1.Message_TextMessage{
			TextMessage: &modelv1.TextMessage{
				Text: messageContent.Text,
			},
		}
	case "label":
		messageContent, _ := utils.TypeConverter[model.LabelMessage](message.Content)

		transformedModel.Payload = &modelv1.Message_LabelMessage{
			LabelMessage: &modelv1.LabelMessage{
				Text: messageContent.Text,
			},
		}
	}

	return transformedModel
}

func TransformMessagesToGrpcModel(messages []model.Message) []*modelv1.Message {
	var transformedMessages []*modelv1.Message

	for _, v := range messages {
		transformedMessages = append(transformedMessages, TransformMessageToGrpcModel(&v))
	}

	return transformedMessages
}
