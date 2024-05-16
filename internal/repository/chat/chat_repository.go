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

type chatRepository struct {
	logger             *zap.Logger
	chatsCollection    *mongo.Collection
	messagesCollection *mongo.Collection
}

func NewRepository(logger *zap.Logger, db *mongo.Database) chat.Repository {
	return &chatRepository{logger, db.Collection(database.ChatsCollection), db.Collection(database.MessagesCollection)}
}

func (repo *chatRepository) GetChatMembers(chatID primitive.ObjectID) []chat.Member {
	// FIXME
	return []chat.Member{}
}

func (repo *chatRepository) Create(newChat chat.Chat) (*chat.Chat, error) {
	_, err := repo.chatsCollection.InsertOne(context.Background(), newChat)
	if err != nil {
		return nil, err
	}

	_, err = repo.messagesCollection.InsertOne(context.Background(), bson.M{
		"chat_id":  newChat.ChatID,
		"messages": bson.A{},
	})
	if err != nil {
		return nil, err
	}

	return &newChat, nil
}

func (repo *chatRepository) Destroy(chatID primitive.ObjectID) error {
	filter := bson.M{"chat_id": chatID}

	_, err := repo.chatsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	_, err = repo.messagesCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo *chatRepository) FindMany(filter bson.M) ([]chat.Chat, error) {
	cursor, err := repo.chatsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var chats []chat.Chat

	decodeErr := cursor.All(context.Background(), &chats)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return chats, nil
}

func (repo *chatRepository) FindOne(filter bson.M) (*chat.Chat, error) {
	var model *chat.Chat
	result := repo.chatsCollection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	err := result.Decode(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (repo *chatRepository) FindByID(staticID primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{"chat_id": staticID}
	return repo.FindOne(filter)
}

func (repo *chatRepository) FindChatOrSidesByStaticID(staticID primitive.ObjectID) (*chat.ChatC, error) {
	ctx := context.TODO()

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

func (repo *chatRepository) FindBySides(sides [2]primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{
		"chat_detail.sides":     sides,
		"chat_detail.chat_type": bson.M{"$ne": "direct"},
	}

	return repo.FindOne(filter)
}
