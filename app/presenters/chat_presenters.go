package presenters

import (
	"Kavka/internal/domain/chat"
)

type ChatDto struct {
	Event string     `json:"event"`
	Chat  *chat.Chat `json:"chat"`
}
