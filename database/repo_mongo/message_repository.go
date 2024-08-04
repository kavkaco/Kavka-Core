package repository_mongo

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/repository"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageRepository struct {
	messagesCollection *mongo.Collection
}

func NewMessageMongoRepository(db *mongo.Database) repository.MessageRepository {
	return &messageRepository{db.Collection(database.MessagesCollection)}
}

func (repo *messageRepository) FindMessage(ctx context.Context, chatID model.ChatID, messageID model.MessageID) (*model.MessageGetter, error) {
	filter := bson.M{"chat_id": chatID}

	result := repo.messagesCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var message model.MessageGetter
	var chatMessages *model.ChatMessages

	err := result.Decode(&chatMessages)
	if err != nil {
		return nil, err
	}

	for i, m := range chatMessages.Messages {
		if m.Message.MessageID == messageID {
			message = *m
			break
		}

		if i == len(chatMessages.Messages)-1 {
			return nil, repository.ErrNotFound
		}
	}

	return &message, nil
}

func (repo *messageRepository) Create(ctx context.Context, chatID model.ChatID) error {
	messageStoreModel := model.ChatMessages{
		ChatID:   chatID,
		Messages: []*model.MessageGetter{},
	}
	_, err := repo.messagesCollection.InsertOne(ctx, messageStoreModel)
	if err != nil {
		return nil
	}

	return nil
}

func (repo *messageRepository) FetchMessages(ctx context.Context, chatID model.ChatID) (*model.ChatMessages, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"chat_id": chatID,
			},
		},
		bson.M{"$unwind": bson.M{"path": "$messages"}},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "messages.sender_id",
				"foreignField": "user_id",
				"as":           "sender",
			},
		},
		bson.M{"$unwind": bson.M{"path": "$sender"}},
		bson.M{
			"$group": bson.M{
				"_id": "$_id",
				"chat_id": bson.M{
					"$first": "$chat_id",
				},
				"messages": bson.M{
					"$push": bson.M{
						"message": "$messages",
						"sender":  "$sender",
					},
				},
			},
		},
	}

	cursor, err := repo.messagesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &model.ChatMessages{}, nil
		}

		return nil, err
	}

	var chatMessages []model.ChatMessages
	err = cursor.All(ctx, &chatMessages)
	if err != nil {
		return nil, err
	}

	return &chatMessages[0], nil
}

func (repo *messageRepository) Insert(ctx context.Context, chatID model.ChatID, message *model.Message) (*model.Message, error) {
	filter := bson.M{"chat_id": chatID}
	update := bson.M{"$push": bson.M{"messages": message}}

	result, err := repo.messagesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		if database.IsRowExistsError(err) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	if result.MatchedCount == 0 || result.ModifiedCount == 0 {
		return nil, repository.ErrNotModified
	}

	return message, nil
}

func (repo *messageRepository) UpdateMessageContent(ctx context.Context, chatID model.ChatID, messageID model.MessageID, newMessageContent string) error {
	return repo.updateMessageFields(ctx, chatID, messageID, bson.M{"$set": bson.M{
		"messages.$.content.data": newMessageContent,
	}})
}

func (repo *messageRepository) updateMessageFields(ctx context.Context, chatID model.ChatID, messageID model.MessageID, updateQuery bson.M) error {
	filter := bson.M{"chat_id": chatID, "messages": bson.M{"$elemMatch": bson.M{"message_id": messageID}}}

	result, err := repo.messagesCollection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		if database.IsRowExistsError(err) {
			return repository.ErrNotFound
		}

		return err
	}

	if result.MatchedCount == 0 || result.ModifiedCount == 0 {
		return repository.ErrNotModified
	}

	return nil
}

func (repo *messageRepository) Delete(ctx context.Context, chatID model.ChatID, messageID model.MessageID) error {
	filter := bson.M{"chat_id": chatID}
	update := bson.M{"$pull": bson.M{"messages": bson.M{"message_id": messageID}}}

	result, err := repo.messagesCollection.UpdateOne(ctx, filter, update)
	if err != nil && result.ModifiedCount != 1 {
		if database.IsRowExistsError(err) {
			return repository.ErrNotFound
		}

		return err
	}

	return nil
}
