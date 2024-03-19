package presenters

import (
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/utils"
)

type ChatDto struct {
	Event string     `json:"event"`
	Chat  *chat.Chat `json:"chat"`
}

func ChatAsJSON(obj chat.Chat) (interface{}, error) {
	if obj.ChatType == chat.TypeDirect {
		obj.ChatDetail = nil
	} else {
		// Determine the specific ChatDetail type based on chatType
		var chatDetail interface{}
		var convertErr error

		switch obj.ChatType {
		case chat.TypeChannel:
			chatDetail, convertErr = utils.TypeConverter[chat.ChannelChatDetail](obj.ChatDetail)
		case chat.TypeGroup:
			chatDetail, convertErr = utils.TypeConverter[chat.GroupChatDetail](obj.ChatDetail)
		}

		if convertErr != nil {
			return nil, convertErr
		}

		obj.ChatDetail = chatDetail
	}

	return obj, nil
}
