package presenters

import (
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
)

type ChatDto struct {
	Event string     `json:"event"`
	Chat  *chat.Chat `json:"chat"`
}

func ChatAsJSON(_ string, obj *chat.Chat) interface{} {
	if obj.ChatType == chat.TypeDirect {
		obj.ChatDetail = nil
	}

	return obj
}
