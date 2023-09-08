package presenters

import (
	"Kavka/internal/domain/chat"
)

type ChatDto struct {
	Event string     `json:"event"`
	Chat  *chat.Chat `json:"chat"`
}

func ChatAsJSON(event string, obj *chat.Chat) interface{} {
	if obj.ChatType == chat.ChatTypeDirect {
		obj.ChatDetail = nil
	}

	return obj
}
