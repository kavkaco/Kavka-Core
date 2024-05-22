package repository

import (
	"context"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository interface {
	SearchInChats(ctx context.Context, key string) ([]model.Chat, error)
	UpdateChatLastMessage(ctx context.Context, chatID model.ChatID, lastMessage model.LastMessage) error
	Create(ctx context.Context, chatModel model.Chat) (*model.Chat, error)
	Destroy(ctx context.Context, chatID model.ChatID) error
	FindMany(ctx context.Context, chatIDs []model.ChatID) ([]model.Chat, error)
	FindByID(ctx context.Context, chatID model.ChatID) (*model.Chat, error)
	FindBySides(ctx context.Context, sides [2]model.UserID) (*model.Chat, error)
	GetChatMembers(chatID model.ChatID) []model.Member
}

type chatRepository struct {
	usersCollection *mongo.Collection
	chatsCollection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) ChatRepository {
	return &chatRepository{db.Collection(database.UsersCollection), db.Collection(database.ChatsCollection)}
}

func (repo *chatRepository) UpdateChatLastMessage(ctx context.Context, chatID model.ChatID, lastMessage model.LastMessage) error {
	filter := bson.M{"_id": chatID}
	update := bson.M{"$set": bson.M{"last_message": lastMessage}}

	result, err := repo.chatsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 || result.ModifiedCount == 0 {
		return ErrNotModified
	}

	return nil
}

// TODO
func (repo *chatRepository) SearchInChats(ctx context.Context, key string) ([]model.Chat, error) {
	panic("unimplemented")
}

// TODO
func (repo *chatRepository) GetChatMembers(chatID model.ChatID) []model.Member {
	return []model.Member{}
}

func (repo *chatRepository) Create(ctx context.Context, chatModel model.Chat) (*model.Chat, error) {
	result, err := repo.chatsCollection.InsertOne(ctx, chatModel)
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(model.ChatID)
	chatModel.ChatID = insertedID

	return &chatModel, nil
}

func (repo *chatRepository) Destroy(ctx context.Context, chatID model.ChatID) error {
	filter := bson.M{"_id": chatID}

	result, err := repo.chatsCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrNotModified
	}

	return nil
}

func (repo *chatRepository) FindMany(ctx context.Context, chatIDs []model.ChatID) ([]model.Chat, error) {
	filter := bson.M{"_id": bson.M{"$in": chatIDs}}

	cursor, err := repo.chatsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var chats []model.Chat

	decodeErr := cursor.All(ctx, &chats)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return chats, nil
}

func (repo *chatRepository) findOne(ctx context.Context, filter bson.M) (*model.Chat, error) {
	var model *model.Chat
	result := repo.chatsCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	err := result.Decode(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (repo *chatRepository) FindByID(ctx context.Context, chatID model.ChatID) (*model.Chat, error) {
	filter := bson.M{"_id": chatID}
	return repo.findOne(ctx, filter)
}

func (repo *chatRepository) FindBySides(ctx context.Context, sides [2]model.UserID) (*model.Chat, error) {
	filter := bson.M{
		"chat_detail.sides":     sides,
		"chat_detail.chat_type": bson.M{"$ne": "direct"},
	}

	return repo.findOne(ctx, filter)
}
