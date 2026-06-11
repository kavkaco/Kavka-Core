package repository_mongo

import (
	"context"
	"errors"
	"time"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MessagesV2Collection = "messages_v2"

type messageV2Repository struct {
	messagesV2Collection *mongo.Collection
	usersCollection      *mongo.Collection
}

func NewMessageV2MongoRepository(db *mongo.Database) *messageV2Repository {
	return &messageV2Repository{
		messagesV2Collection: db.Collection(MessagesV2Collection),
		usersCollection:      db.Collection(database.UsersCollection),
	}
}

func (repo *messageV2Repository) ensureIndexes(ctx context.Context) error {
	_, err := repo.messagesV2Collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "_id", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "chat_id", Value: 1}},
		},
	})
	return err
}

func (repo *messageV2Repository) InsertV2(ctx context.Context, doc *model.MessageDocumentV2) (*model.MessageDocumentV2, error) {
	doc.ID = model.NewMessageID()
	doc.CreatedAt = time.Now()

	_, err := repo.messagesV2Collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (repo *messageV2Repository) FetchMessagesV2(ctx context.Context, chatID model.ChatID, skip, limit int64) ([]*model.MessageDocumentV2, error) {
	filter := bson.M{"chat_id": chatID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := repo.messagesV2Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var messages []*model.MessageDocumentV2
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (repo *messageV2Repository) FetchMessageV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID) (*model.MessageDocumentV2, error) {
	filter := bson.M{
		"chat_id": chatID,
		"_id":     messageID,
	}

	result := repo.messagesV2Collection.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, repository.ErrNotFound
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	var doc model.MessageDocumentV2
	err := result.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (repo *messageV2Repository) FetchLastMessageV2(ctx context.Context, chatID model.ChatID) (*model.MessageDocumentV2, error) {
	filter := bson.M{"chat_id": chatID}
	opts := options.FindOne().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	result := repo.messagesV2Collection.FindOne(ctx, filter, opts)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, repository.ErrNotFound
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	var doc model.MessageDocumentV2
	err := result.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (repo *messageV2Repository) UpdateMessageContentV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error {
	filter := bson.M{
		"chat_id": chatID,
		"_id":     messageID,
	}
	update := bson.M{
		"$set": bson.M{
			"content.text": newMessageContent,
			"edited":       true,
		},
	}

	result, err := repo.messagesV2Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (repo *messageV2Repository) DeleteV2(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error {
	filter := bson.M{
		"chat_id": chatID,
		"_id":     messageID,
	}

	result, err := repo.messagesV2Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return repository.ErrNotModified
	}

	return nil
}

func (repo *messageV2Repository) CountMessagesV2(ctx context.Context, chatID model.ChatID) (int64, error) {
	filter := bson.M{"chat_id": chatID}
	return repo.messagesV2Collection.CountDocuments(ctx, filter)
}
