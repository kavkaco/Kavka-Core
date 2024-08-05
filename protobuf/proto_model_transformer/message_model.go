package proto_model_transformer

import (
	"github.com/kavkaco/Kavka-Core/internal/model"
	messagev1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/model/message/v1"
	"github.com/kavkaco/Kavka-Core/utils"
)

func MessageToProto(messageGetter *model.MessageGetter) *messagev1.Message {
	message := messageGetter.Message

	transformedModel := &messagev1.Message{
		MessageId: message.MessageID.Hex(),
		Sender:    MessageSenderToProto(messageGetter.Sender),
		CreatedAt: message.CreatedAt.Unix(),
		Edited:    message.Edited,
		Seen:      message.Seen,
		Type:      message.Type,
	}

	switch transformedModel.Type {
	case "text":
		messageContent, _ := utils.TypeConverter[model.TextMessage](message.Content)

		transformedModel.Payload = &messagev1.Message_TextMessage{
			TextMessage: &messagev1.TextMessage{
				Text: messageContent.Text,
			},
		}
	case "label":
		messageContent, _ := utils.TypeConverter[model.LabelMessage](message.Content)

		transformedModel.Payload = &messagev1.Message_LabelMessage{
			LabelMessage: &messagev1.LabelMessage{
				Text: messageContent.Text,
			},
		}
	}

	return transformedModel
}

func MessageSenderToProto(messageSender *model.MessageSender) *messagev1.MessageSender {
	return &messagev1.MessageSender{
		UserId:   messageSender.UserID,
		Name:     messageSender.Name,
		LastName: messageSender.LastName,
		Username: messageSender.Username,
	}
}

func MessagesToProto(messageGetters []*model.MessageGetter) []*messagev1.Message {
	var transformedMessages []*messagev1.Message

	for _, v := range messageGetters {
		transformedMessages = append(transformedMessages, MessageToProto(v))
	}

	return transformedMessages
}
