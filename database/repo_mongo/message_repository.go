package repository_mongo

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/repository"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type messagesDoc struct {
	ChatID   model.ChatID           `bson:"chat_id"`
	Messages []*model.MessageGetter `bson:"messages"`
}

type messageRepository struct {
	messagesCollection *mongo.Collection
	usersCollection    *mongo.Collection
}

func NewMessageMongoRepository(db *mongo.Database) repository.MessageRepository {
	return &messageRepository{db.Collection(database.MessagesCollection), db.Collection(database.UsersCollection)}
}

func (repo *messageRepository) FetchMessage(ctx context.Context, chatID primitive.ObjectID, messageID primitive.ObjectID) (*model.Message, error) {
	cursor, err := repo.messagesCollection.Aggregate(ctx, bson.A{
		bson.M{
			"$match": bson.M{
				"chat_id": chatID,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":     1,
				"chat_id": 1,
				"message": bson.M{
					"$arrayElemAt": bson.A{
						bson.M{
							"$filter": bson.M{
								"input": "$messages",
								"as":    "message",
								"cond": bson.M{"$eq": bson.A{
									"$$message.message_id", messageID,
								}},
							},
						},
						0,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	type docModel struct {
		ChatID  model.ChatID  `bson:"chat_id"`
		Message model.Message `bson:"message"`
	}

	var docs []docModel
	err = cursor.All(ctx, &docs)
	if err != nil || len(docs) == 0 {
		return nil, err
	}

	message := &docs[0].Message

	return message, nil
}

func (repo *messageRepository) FetchLastMessage(ctx context.Context, chatID model.ChatID) (*model.Message, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"chat_id": chatID,
			},
		},
		bson.M{
			"$addFields": bson.M{
				"last_message": bson.M{
					"$last": "$messages",
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"messages": 0,
			},
		},
	}

	cursor, err := repo.messagesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	type doc struct {
		LastMessage *model.Message `bson:"last_message"`
	}

	var docs []*doc

	decodeErr := cursor.All(ctx, &docs)
	if decodeErr != nil {
		return nil, decodeErr
	}

	if len(docs) > 0 && docs[0] != nil {
		return docs[0].LastMessage, nil
	}

	return &model.Message{}, nil
}

func (repo *messageRepository) Create(ctx context.Context, chatID model.ChatID) error {
	messageStoreModel := messagesDoc{
		ChatID:   chatID,
		Messages: []*model.MessageGetter{},
	}
	_, err := repo.messagesCollection.InsertOne(ctx, messageStoreModel)
	if err != nil {
		return nil
	}

	return nil
}

func (repo *messageRepository) FetchMessages(ctx context.Context, chatID model.ChatID) ([]*model.MessageGetter, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"chat_id": chatID,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "messages.sender_id",
				"foreignField": "user_id",
				"as":           "senders",
			},
		},
		bson.M{
			"$addFields": bson.M{
				"fetched_messages": bson.M{
					"$map": bson.M{
						"input": "$messages",
						"as":    "message",
						"in": bson.M{
							"sender": bson.M{
								"$arrayElemAt": bson.A{
									bson.M{
										"$filter": bson.M{
											"input": "$senders",
											"as":    "sender",
											"cond": bson.M{
												"$eq": bson.A{"$$sender.user_id", "$$message.sender_id"},
											},
										},
									},
									0,
								},
							},
							"message": "$$message",
						},
					},
				},
			},
		},
	}

	cursor, err := repo.messagesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*model.MessageGetter{}, nil
		}

		return nil, err
	}

	type doc struct {
		ChatID   model.ChatID           `bson:"chat_id"`
		Messages []*model.MessageGetter `bson:"fetched_messages"`
	}

	var docs []*doc

	err = cursor.All(ctx, &docs)
	if err != nil {
		return nil, err
	}

	if len(docs) > 0 {
		return docs[0].Messages, nil
	}

	return []*model.MessageGetter{}, nil
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
