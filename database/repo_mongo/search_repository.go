package repository_mongo

import (
	"context"
	"time"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type searchRepository struct {
	chatRepository  *mongo.Collection
	usersRepository *mongo.Collection
}

func NewSearchRepository(db *mongo.Database) repository.SearchRepository {
	return &searchRepository{db.Collection(database.ChatsCollection), db.Collection(database.UsersCollection)}
}

func (s *searchRepository) Search(ctx context.Context, input string) (*model.SearchResultDTO, error) {
	// Search in chats collection
	cursor, err := s.chatRepository.Find(ctx, bson.M{
		"$text": bson.M{
			"$search": input,
		},
	}, options.Find().SetLimit(10).SetMaxTime(10*time.Second))
	if err != nil {
		return nil, err
	}

	var chats []model.ChatDTO
	err = cursor.All(ctx, &chats)
	if err != nil {
		return nil, err
	}

	// Search in users collection
	cursor, err = s.usersRepository.Find(ctx, bson.M{
		"$text": bson.M{
			"$search": input,
		},
	})
	if err != nil {
		return nil, err
	}

	var users []model.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return &model.SearchResultDTO{
		Chats: chats,
		Users: users,
	}, nil
}

// SearchInChat implements repository.SearchRepository.
func (s *searchRepository) SearchInChat(ctx context.Context, input string) (*model.MessageGetter, error) {
	panic("unimplemented")
}
