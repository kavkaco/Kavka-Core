package repository_mongo

import (
	"context"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type chatRepository struct {
	usersCollection *mongo.Collection
	chatsCollection *mongo.Collection
}

func NewChatMongoRepository(db *mongo.Database) repository.ChatRepository {
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
		return repository.ErrNotModified
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
		return repository.ErrNotModified
	}

	return nil
}

func (repo *chatRepository) FindManyByChatID(ctx context.Context, chatIDs []model.ChatID) ([]model.Chat, error) {
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
func (repo *chatRepository) IsUsernameOccupied(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}

	chatCount, err := repo.chatsCollection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	userCount, err := repo.usersCollection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return (userCount == 1 || chatCount == 1), nil
}
