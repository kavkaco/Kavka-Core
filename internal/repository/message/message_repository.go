package repository

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model/message"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	ErrChatNotFound = errors.New("chat not found")
	ErrMsgNotFound  = errors.New("message not found")
	ErrNoAccess     = errors.New("no access")
)

type MessageRepository interface {
	Insert(ctx context.Context, chatID primitive.ObjectID, msg *message.Message) (*message.Message, error)
	Update(ctx context.Context, chatID primitive.ObjectID, messageID primitive.ObjectID, fieldsToUpdate bson.M) error
	Delete(ctx context.Context, chatID primitive.ObjectID, messageID primitive.ObjectID) error
}

type messageRepository struct {
	logger             *zap.Logger
	messagesCollection *mongo.Collection
}

func NewRepository(logger *zap.Logger, db *mongo.Database) MessageRepository {
	return &messageRepository{logger, db.Collection(database.MessagesCollection)}
}

func (repo *messageRepository) Insert(ctx context.Context, chatID primitive.ObjectID, msg *message.Message) (*message.Message, error) {
	filter := bson.M{"chat_id": chatID}
	update := bson.M{"$push": bson.M{"messages": msg}}

	result := repo.messagesCollection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		if database.IsRowExistsError(result.Err()) {
			return nil, ErrChatNotFound
		}
		return nil, result.Err()
	}

	return msg, nil
}

func (repo *messageRepository) Update(ctx context.Context, chatID primitive.ObjectID, messageID primitive.ObjectID, fieldsToUpdate bson.M) error {
	filter := bson.M{"chat_id": chatID, "messages.message_id": messageID}
	update := bson.M{"$set": fieldsToUpdate}

	_, err := repo.messagesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		if database.IsRowExistsError(err) {
			return ErrChatNotFound
		}

		return err
	}

	return nil
}

func (repo *messageRepository) Delete(ctx context.Context, chatID primitive.ObjectID, messageID primitive.ObjectID) error {
	filter := bson.M{"chat_id": chatID}
	update := bson.M{"$pull": bson.M{"messages": bson.M{"message_id": messageID}}}

	result, err := repo.messagesCollection.UpdateOne(ctx, filter, update)
	if err != nil && result.ModifiedCount != 1 {
		if database.IsRowExistsError(err) {
			return ErrChatNotFound
		}

		return err
	}

	return nil
}
