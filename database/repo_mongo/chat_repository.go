package repository_mongo

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type chatRepository struct {
	usersCollection *mongo.Collection
	chatsCollection *mongo.Collection
}

func NewChatMongoRepository(db *mongo.Database) repository.ChatRepository {
	return &chatRepository{db.Collection(database.UsersCollection), db.Collection(database.ChatsCollection)}
}

func (repo *chatRepository) JoinChat(ctx context.Context, userID string, chatID primitive.ObjectID) error {
	userFilter := bson.M{"user_id": userID}
	userUpdate := bson.M{
		"$addToSet": bson.M{
			"chats_list_ids": chatID,
		},
	}
	_, err := repo.usersCollection.UpdateOne(ctx, userFilter, userUpdate)
	if err != nil {
		return err
	}

	chatFilter := bson.M{"_id": chatID}
	chatUpdate := bson.M{
		"$addToSet": bson.M{
			"chat_detail.members": userID,
		},
	}

	_, err = repo.chatsCollection.UpdateOne(ctx, chatFilter, chatUpdate)
	if err != nil {
		return err
	}

	return nil
}

// FIXME
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

func (repo *chatRepository) GetUserChats(ctx context.Context, chatIDs []model.ChatID) ([]model.ChatDTO, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"_id": bson.M{"$in": chatIDs},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "messages",
				"localField":   "_id",
				"foreignField": "chat_id",
				"as":           "chat_messages",
			},
		},
		bson.M{
			"$addFields": bson.M{
				"last_message": bson.M{
					"$arrayElemAt": bson.A{
						bson.M{
							"$arrayElemAt": bson.A{
								"$chat_messages.messages",
								bson.M{"$subtract": bson.A{
									bson.M{"$size": "$chat_messages.messages"},
									1,
								}},
							},
						},
						0,
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"chat_messages": 0,
			},
		},
	}

	cursor, err := repo.chatsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var chats []model.ChatDTO

	decodeErr := cursor.All(ctx, &chats)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return chats, nil
}

func (repo *chatRepository) findOne(ctx context.Context, filter bson.M) (*model.Chat, error) {
	var model *model.Chat

	result := repo.chatsCollection.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, repository.ErrNotFound
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	err := result.Decode(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (repo *chatRepository) GetChat(ctx context.Context, chatID model.ChatID) (*model.Chat, error) {
	result := repo.chatsCollection.FindOne(ctx, bson.M{
		"_id": chatID,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var chat *model.Chat
	err := result.Decode(&chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (repo *chatRepository) FindBySides(ctx context.Context, userID model.UserID, recipientUserID model.UserID) (*model.Chat, error) {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"$and": bson.A{
				bson.M{"chat_detail.user_id": userID},
				bson.M{"chat_detail.recipient_user_id": recipientUserID},
			}},
			bson.M{"$and": bson.A{
				bson.M{"chat_detail.user_id": recipientUserID},
				bson.M{"chat_detail.recipient_user_id": userID},
			}},
		},
		"chat_detail.chat_type": bson.M{"$ne": "direct"},
	}

	return repo.findOne(ctx, filter)
}
