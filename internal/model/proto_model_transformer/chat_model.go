package proto_model_transformer

import (
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/utils"
	modelv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/model/chat/v1"
)

var ErrTransformation = errors.New("unable to transform the model")

func ChatToProto(chat model.ChatDTO) (*modelv1.Chat, error) {
	var chatType modelv1.ChatType
	switch chat.ChatType {
	case "channel":
		chatType = modelv1.ChatType_CHAT_TYPE_CHANNEL
	case "group":
		chatType = modelv1.ChatType_CHAT_TYPE_GROUP
	case "direct":
		chatType = modelv1.ChatType_CHAT_TYPE_DIRECT
	}

	chatDetailProto, err := ChatDetailToProto(chat.ChatType, chat.ChatDetail)
	if err != nil {
		return nil, err
	}

	lastMessage := &modelv1.LastMessage{
		MessageType:    "",
		MessageCaption: "",
	}

	if chat.LastMessage != nil {
		lastMessage.MessageType = chat.LastMessage.Type

		switch chat.LastMessage.Type {
		case model.TypeTextMessage, model.TypeLabelMessage:
			messageContent, err := utils.TypeConverter[model.TextMessage](chat.LastMessage.Content)
			if err != nil {
				return nil, err
			}

			lastMessage.MessageCaption = messageContent.Text
		}
	}

	return &modelv1.Chat{
		ChatId:      chat.ChatID.Hex(),
		ChatType:    chatType,
		ChatDetail:  chatDetailProto,
		LastMessage: lastMessage,
	}, nil
}

var transformedChats []*modelv1.Chat

func ChatsToProto(chats []model.ChatDTO) ([]*modelv1.Chat, error) {
	transformedChats = []*modelv1.Chat{}

	for _, v := range chats {
		c, err := ChatToProto(v)
		if err != nil {
			return nil, err
		}

		transformedChats = append(transformedChats, c)
	}

	return transformedChats, nil
}

func ChatDetailToProto(chatType string, chatDetail interface{}) (*modelv1.ChatDetail, error) {
	switch chatType {
	case "channel":
		cd, err := utils.TypeConverter[model.ChannelChatDetail](chatDetail)
		if err != nil {
			return nil, err
		}

		return &modelv1.ChatDetail{
			ChatDetailType: &modelv1.ChatDetail_ChannelDetail{
				ChannelDetail: &modelv1.ChannelChatDetail{
					Title:        cd.Title,
					Username:     cd.Username,
					Members:      cd.Members,
					Admins:       cd.Admins,
					Owner:        cd.Owner,
					Description:  cd.Description,
					RemovedUsers: cd.RemovedUsers,
				},
			},
		}, nil
	case "group":
		cd, err := utils.TypeConverter[model.GroupChatDetail](chatDetail)
		if err != nil {
			return nil, err
		}

		return &modelv1.ChatDetail{
			ChatDetailType: &modelv1.ChatDetail_GroupDetail{
				GroupDetail: &modelv1.GroupChatDetail{
					Title:        cd.Title,
					Username:     cd.Username,
					Members:      cd.Members,
					Admins:       cd.Admins,
					Owner:        cd.Owner,
					Description:  cd.Description,
					RemovedUsers: cd.RemovedUsers,
				},
			},
		}, nil
	case "direct":
		cd, err := utils.TypeConverter[model.DirectChatDetailDTO](chatDetail)
		if err != nil {
			return nil, err
		}

		if cd == nil || cd.Recipient == nil {
			return nil, errors.New("unable to transform direct chat detail to DirectChatDetailDTO model")
		}

		return &modelv1.ChatDetail{
			ChatDetailType: &modelv1.ChatDetail_DirectDetail{
				DirectDetail: &modelv1.DirectChatDetail{
					Recipient: UserToProto(*cd.Recipient),
				},
			},
		}, nil
	}

	return nil, ErrTransformation
}
