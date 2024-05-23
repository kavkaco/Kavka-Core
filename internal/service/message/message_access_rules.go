package message

import (
	"slices"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/utils"
)

func HasAccessToSendMessage(chatType string, detail interface{}, userID model.UserID) bool {
	if chatType == model.TypeDirect {
		detail, err := utils.TypeConverter[model.DirectChatDetail](detail)
		if err != nil {
			return false
		}
		return ha_SendMessage_Direct(detail, userID)
	} else if chatType == model.TypeChannel {
		detail, err := utils.TypeConverter[model.ChannelChatDetail](detail)
		if err != nil {
			return false
		}
		return ha_SendMessage_Channel(detail, userID)
	} else if chatType == model.TypeGroup {
		detail, err := utils.TypeConverter[model.GroupChatDetail](detail)
		if err != nil {
			return false
		}
		return ha_SendMessage_Group(detail, userID)
	}

	return false
}

func HasAccessToDeleteMessage(chatType string, detail interface{}, userID model.UserID, message model.Message) bool {
	if chatType == model.TypeDirect {
		detail, err := utils.TypeConverter[model.DirectChatDetail](detail)
		if err != nil {
			return false
		}
		return ha_DeleteMessage_Direct(detail, userID, &message)
	} else if chatType == model.TypeChannel {
		detail, err := utils.TypeConverter[model.ChannelChatDetail](detail)
		if err != nil {
			return false
		}
		return ha_DeleteMessage_Channel(detail, userID, &message)
	} else if chatType == model.TypeGroup {
		detail, err := utils.TypeConverter[model.GroupChatDetail](detail)
		if err != nil {
			return false
		}
		return ha_DeleteMessage_Group(detail, userID, &message)
	}

	return false
}

// Being a member of a group is enough to have access to send the message.
func ha_SendMessage_Group(detail *model.GroupChatDetail, userID model.UserID) bool {
	return slices.Contains(detail.Members, userID)
}

// Only admins of the channel chat can send messages.
func ha_SendMessage_Channel(detail *model.ChannelChatDetail, userID model.UserID) bool {
	return slices.Contains(detail.Admins, userID)
}

// Both users can send message to their direct chat
func ha_SendMessage_Direct(_ *model.DirectChatDetail, _ model.UserID) bool {
	return true
}

// The user is only allowed to delete messages his own messages.
//
//	Admins can delete any messages.
func ha_DeleteMessage_Group(detail *model.GroupChatDetail, userID model.UserID, message *model.Message) bool {
	if message.SenderID == userID {
		return true
	}

	// If is admin
	return slices.Contains(detail.Admins, userID)
}

// Only admins of the channel chat can delete messages.
func ha_DeleteMessage_Channel(detail *model.ChannelChatDetail, userID model.UserID, _ *model.Message) bool {
	return slices.Contains(detail.Admins, userID)
}

// Both users can delete their message in direct chat
func ha_DeleteMessage_Direct(detail *model.DirectChatDetail, userID model.UserID, _ *model.Message) bool {
	return detail.HasSide(userID)
}
