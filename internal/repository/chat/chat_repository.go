package repository

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
)

type ChatRepository interface {
	GetChatMembers(chatID primitive.ObjectID) []chat.Member
	Create(ctx context.Context, newChat chat.Chat) (*chat.Chat, error)
	Destroy(ctx context.Context, chatID primitive.ObjectID) error
	FindMany(ctx context.Context, filter bson.M) ([]chat.Chat, error)
	FindOne(ctx context.Context, filter bson.M) (*chat.Chat, error)
	FindByID(ctx context.Context, staticID primitive.ObjectID) (*chat.Chat, error)
	FindChatOrSidesByStaticID(ctx context.Context, staticID primitive.ObjectID) (*chat.ChatC, error)
	FindBySides(ctx context.Context, sides [2]primitive.ObjectID) (*chat.Chat, error)
	AddToUserChats(ctx context.Context, userStaticID primitive.ObjectID, chatID primitive.ObjectID) (ok bool, err error)
}

type chatRepository struct {
	logger             *zap.Logger
	usersCollection    *mongo.Collection
	chatsCollection    *mongo.Collection
	messagesCollection *mongo.Collection
}

func NewRepository(logger *zap.Logger, db *mongo.Database) ChatRepository {
	return &chatRepository{logger, db.Collection(database.ChatsCollection), db.Collection(database.UsersCollection), db.Collection(database.MessagesCollection)}
}

func (repo *chatRepository) AddToUserChats(ctx context.Context, userStaticID primitive.ObjectID, chatID primitive.ObjectID) (ok bool, err error) {
	result := repo.usersCollection.FindOneAndUpdate(ctx, bson.M{"static_id": userStaticID}, bson.M{
		"$push": bson.M{
			"chats": primitive.A{chatID},
		},
	})
	if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}

func (repo *chatRepository) GetChatMembers(chatID primitive.ObjectID) []chat.Member {
	// FIXME
	return []chat.Member{}
}

func (repo *chatRepository) Create(ctx context.Context, newChat chat.Chat) (*chat.Chat, error) {
	// Insert chat
	_, err := repo.chatsCollection.InsertOne(ctx, newChat)
	if err != nil {
		return nil, err
	}

	// Create messages
	_, err = repo.messagesCollection.InsertOne(ctx, bson.M{
		"chat_id":  newChat.ChatID,
		"messages": bson.A{},
	})
	if err != nil {
		return nil, err
	}

	return &newChat, nil
}

func (repo *chatRepository) Destroy(ctx context.Context, chatID primitive.ObjectID) error {
	filter := bson.M{"chat_id": chatID}

	_, err := repo.chatsCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	_, err = repo.messagesCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo *chatRepository) FindMany(ctx context.Context, filter bson.M) ([]chat.Chat, error) {
	cursor, err := repo.chatsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var chats []chat.Chat

	decodeErr := cursor.All(ctx, &chats)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return chats, nil
}

func (repo *chatRepository) FindOne(ctx context.Context, filter bson.M) (*chat.Chat, error) {
	var model *chat.Chat
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

func (repo *chatRepository) FindByID(ctx context.Context, staticID primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{"chat_id": staticID}
	return repo.FindOne(ctx, filter)
}

func (repo *chatRepository) FindChatOrSidesByStaticID(ctx context.Context, staticID primitive.ObjectID) (*chat.ChatC, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$or": []interface{}{
					bson.M{"chat_id": staticID},
					bson.M{"chat_detail.sides": bson.M{"$in": []interface{}{staticID}}},
					bson.M{"chat_detail.members": bson.M{"$in": []interface{}{staticID}}},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "chat_detail.sides",
				"foreignField": "id",
				"as":           "chat_detail.fetchedUsers",
				"pipeline": []interface{}{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{"$eq": []interface{}{"$chat_type", "direct"}},
						},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "messages",
				"localField":   "chat_id",
				"foreignField": "chat_id",
				"as":           "chatMessages",
			},
		},
		{
			"$unwind": "$chatMessages",
		},
		{
			"$project": bson.M{
				"chat_id":     1,
				"chat_type":   1,
				"chat_detail": 1,
				"messages":    bson.M{"$slice": []interface{}{"$chatMessages.messages", -1}},
			},
		},
	}

	cursor, err := repo.chatsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var foundChats []chat.ChatC
	if err := cursor.All(ctx, &foundChats); err != nil {
		return nil, err
	}

	if len(foundChats) == 0 {
		return nil, ErrChatNotFound
	}

	return &foundChats[0], nil
}

func (repo *chatRepository) FindBySides(ctx context.Context, sides [2]primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{
		"chat_detail.sides":     sides,
		"chat_detail.chat_type": bson.M{"$ne": "direct"},
	}

	return repo.FindOne(ctx, filter)
}
