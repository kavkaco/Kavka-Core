package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

type UserRepository interface {
	GetChats(ctx context.Context, userID model.UserID) ([]model.ChatID, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	AddToUserChats(ctx context.Context, userID model.UserID, chatID model.ChatID) error
	Update(ctx context.Context, userID string, name, lastName, username, biography string) error
	FindByUserID(ctx context.Context, userID model.UserID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	DeleteByID(ctx context.Context, userID model.UserID) error
	IsIndexesUnique(ctx context.Context, email string, username string) (isUnique bool, unUniqueFields []string)
	IsUsernameOccupied(ctx context.Context,username string)(bool,error)
}
