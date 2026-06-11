package repository_mongo

import (
	"context"
	"time"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type searchRepository struct {
	chatRepository  *mongo.Collection
	usersRepository *mongo.Collection
}

func NewSearchRepository(db *mongo.Database) repository.SearchRepository {
	return &searchRepository{db.Collection(database.ChatsCollection), db.Collection(database.UsersCollection)}
}

func (s *searchRepository) Search(ctx context.Context, input string, limit int) (*model.SearchResultDTO, error) {
	// Search in chats collection
	cursor, err := s.chatRepository.Find(ctx, bson.M{
		"$text": bson.M{
			"$search": input,
		},
	}, options.Find().SetLimit(int64(limit)).SetMaxTime(10*time.Second))
	if err != nil {
		return nil, err
	}

	var chats []model.ChatDTO
	err = cursor.All(ctx, &chats)
	if err != nil {
		return nil, err
	}

	// Search in users collection
	cursor, err = s.usersRepository.Find(ctx, bson.M{
		"$text": bson.M{
			"$search": input,
		},
	}, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}

	var users []model.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return &model.SearchResultDTO{
		Chats: chats,
		Users: users,
	}, nil
}

// SearchInChat implements repository.SearchRepository.
func (s *searchRepository) SearchInChat(ctx context.Context, chatID model.ChatID, input string) ([]*model.MessageGetter, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"chat_id": chatID,
			},
		},
		bson.M{
			"$addFields": bson.M{
				"messages": bson.M{
					"$filter": bson.M{
						"input": "$messages",
						"as":    "msg",
						"cond": bson.M{
							"$regexMatch": bson.M{
								"input": bson.M{"$ifNull": bson.A{"$$msg.content.text", ""}},
								"regex": input,
								"options": "i",
							},
						},
					},
				},
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

	cursor, err := s.chatRepository.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	type doc struct {
		FetchedMessages []*model.MessageGetter `bson:"fetched_messages"`
	}

	var docs []doc
	err = cursor.All(ctx, &docs)
	if err != nil {
		return nil, err
	}

	if len(docs) > 0 {
		return docs[0].FetchedMessages, nil
	}

	return []*model.MessageGetter{}, nil
}
