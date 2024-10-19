package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type SearchRepository interface {
	Search(ctx context.Context, input string) (*model.SearchResultDTO, error)
	SearchInChat(ctx context.Context, input string) (*model.MessageGetter, error)
}
