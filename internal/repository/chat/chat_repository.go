package repository

import (
	"Kavka/database"
	"Kavka/internal/domain/chat"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
)

type ChatRepository struct {
	chatsCollection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) *ChatRepository {
	return &ChatRepository{
		db.Collection(database.ChatsCollection),
	}
}

func (repo *ChatRepository) Create(chatType string, chatDetail interface{}) (*chat.Chat, error) {
	chat := chat.NewChat(chatType, chatDetail)
	result, err := repo.chatsCollection.InsertOne(context.Background(), chat)
	if err != nil {
		return nil, err
	}

	chat.ChatID = result.InsertedID.(primitive.ObjectID)

	return chat, nil
}

func (repo *ChatRepository) Destroy(chatID primitive.ObjectID) error {
	filter := bson.M{"_id": chatID}

	_, err := repo.chatsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ChatRepository) Where(filter any) ([]*chat.Chat, error) {
	cursor, err := repo.chatsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var chats []*chat.Chat

	err = cursor.All(context.Background(), &chats)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (repo *ChatRepository) findBy(filter any) (*chat.Chat, error) {
	result, err := repo.Where(filter)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		user := result[len(result)-1]

		return user, nil
	}

	return nil, ErrChatNotFound
}

func (repo *ChatRepository) FindByID(staticID primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.D{{Key: "_id", Value: staticID}}
	return repo.findBy(filter)
}

func (repo *ChatRepository) FindChatOrSidesByStaticID(staticID *primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"chat_detail.sides": bson.M{"$in": []*primitive.ObjectID{staticID}}},
			bson.M{"_id": staticID},
		},
	}

	return repo.findBy(filter)
}

func (repo *ChatRepository) FindBySides(sides [2]*primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{
		"chat_detail.sides":     sides,
		"chat_detail.chat_type": bson.M{"$ne": "direct"},
	}

	return repo.findBy(filter)
}
