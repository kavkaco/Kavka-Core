package message

import (
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/utils"
)

func channelReceiversIDs(chat *model.Chat) ([]model.UserID, error) {
	chatDetail, err := utils.TypeConverter[model.ChannelChatDetail](chat.ChatDetail)
	if err != nil {
		return nil, err
	}

	return chatDetail.Members, nil
}

func groupReceiversIDs(chat *model.Chat) ([]model.UserID, error) {
	chatDetail, err := utils.TypeConverter[model.GroupChatDetail](chat.ChatDetail)
	if err != nil {
		return nil, err
	}

	return chatDetail.Members, nil
}

func directReceiversIDs(chat *model.Chat) ([]model.UserID, error) {
	chatDetail, err := utils.TypeConverter[model.DirectChatDetail](chat.ChatDetail)
	if err != nil {
		return nil, err
	}

	return []model.UserID{chatDetail.UserID, chatDetail.RecipientUserID}, nil
}

func ReceiversIDs(chat *model.Chat) ([]model.UserID, error) {
	switch chat.ChatType {
	case model.TypeChannel:
		return channelReceiversIDs(chat)
	case model.TypeGroup:
		return groupReceiversIDs(chat)
	case model.TypeDirect:
		return directReceiversIDs(chat)
	default:
		return nil, errors.New("invalid chat type detected on method ReceiversIDs")
	}
}
