package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type SearchRepository interface {
	Search(ctx context.Context, input string, limit int) (*model.SearchResultDTO, error)
	SearchInChat(ctx context.Context, chatID model.ChatID, input string) ([]*model.MessageGetter, error)
}
