package repository

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
)

type repository struct {
	chatsCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) chat.Repository {
	return &repository{db.Collection(database.ChatsCollection)}
}

func (repo *repository) Create(newChat chat.Chat) (*chat.Chat, error) {
	_, err := repo.chatsCollection.InsertOne(context.Background(), newChat)
	if err != nil {
		return nil, err
	}

	return &newChat, nil
}

func (repo *repository) Destroy(chatID primitive.ObjectID) error {
	filter := bson.M{"chat_id": chatID}

	_, err := repo.chatsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) FindMany(filter bson.M) ([]chat.Chat, error) {
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

func (repo *repository) FindOne(filter bson.M) (*chat.Chat, error) {
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

func (repo *repository) GetUserChats(userStaticID primitive.ObjectID) ([]chat.Chat, error) {
	ctx := context.TODO()

	memberMatchStage := bson.M{"chat_detail.members": bson.M{"$in": bson.A{userStaticID}}}
	sidesMatchStage := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key: "chat_detail.sides",
			Value: bson.D{{
				Key:   "$in",
				Value: bson.A{userStaticID},
			}},
		}},
	}}

	// Define the pipeline for direct chats (with aggregation for sides)
	directPipeline := mongo.Pipeline{
		sidesMatchStage,
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "users",
				"localField":   "chat_detail.sides",
				"foreignField": "id",
				"as":           "chat_detail.fetchedUsers",
			},
		}},
		{{
			Key: "$project",
			Value: bson.M{
				"chat_id":     1,
				"chat_type":   1,
				"chat_detail": 1,
			},
		}},
	}

	normalCursor, err := repo.chatsCollection.Find(context.TODO(), memberMatchStage)
	if err != nil {
		return []chat.Chat{}, err
	}
	defer normalCursor.Close(ctx)

	// Aggregate for direct chats
	directCursor, err := repo.chatsCollection.Aggregate(ctx, directPipeline)
	if err != nil {
		return nil, err
	}
	defer directCursor.Close(ctx)

	// Decode and merge results
	var allChats []chat.Chat
	if err := normalCursor.All(ctx, &allChats); err != nil {
		return nil, err
	}

	var directChats []chat.Chat
	if err := directCursor.All(ctx, &directChats); err != nil {
		return nil, err
	}

	allChats = append(allChats, directChats...)

	return allChats, nil
}

func (repo *repository) FindByID(staticID primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{"chat_id": staticID}
	return repo.FindOne(filter)
}

func (repo *repository) FindChatOrSidesByStaticID(staticID primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"chat_detail.sides": bson.M{"$in": []primitive.ObjectID{staticID}}},
			bson.M{"chat_id": staticID},
		},
	}

	return repo.FindOne(filter)
}

func (repo *repository) FindBySides(sides [2]primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{
		"chat_detail.sides":     sides,
		"chat_detail.chat_type": bson.M{"$ne": "direct"},
	}

	return repo.FindOne(filter)
}
