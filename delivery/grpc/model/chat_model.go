package grpc_model

import (
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/model"
	modelv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/model/chat/v1"
	"github.com/kavkaco/Kavka-Core/utils"
)

var ErrTransformation = errors.New("unable to transform the model")

func TransformChatToGrpcModel(chat model.Chat) (*modelv1.Chat, error) {
	var chatType modelv1.ChatType
	switch chat.ChatType {
	case "channel":
		chatType = modelv1.ChatType_CHAT_TYPE_CHANNEL
	case "group":
		chatType = modelv1.ChatType_CHAT_TYPE_GROUP
	case "direct":
		chatType = modelv1.ChatType_CHAT_TYPE_DIRECT
	}

	chatDetailGrpcModel, err := TransformChatDetailToGrpcModel(chat.ChatType, chat.ChatDetail)
	if err != nil {
		return nil, err
	}

	return &modelv1.Chat{
		ChatId:     chat.ChatID.Hex(),
		ChatType:   chatType,
		ChatDetail: chatDetailGrpcModel,
	}, nil
}

func TransformChatDetailToGrpcModel(chatType string, chatDetail interface{}) (*modelv1.ChatDetail, error) {
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
					Admins:       cd.Members,
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
					Admins:       cd.Members,
					Owner:        cd.Owner,
					Description:  cd.Description,
					RemovedUsers: cd.RemovedUsers,
				},
			},
		}, nil
	case "direct":
		cd, err := utils.TypeConverter[model.DirectChatDetail](chatDetail)
		if err != nil {
			return nil, err
		}

		return &modelv1.ChatDetail{
			ChatDetailType: &modelv1.ChatDetail_DirectDetail{
				DirectDetail: &modelv1.DirectChatDetail{
					Sides: []string{cd.Sides[0], cd.Sides[1]},
				},
			},
		}, nil
	}

	return nil, ErrTransformation
}
