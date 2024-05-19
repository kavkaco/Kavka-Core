package repository

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/internal/model/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPhoneAlreadyTaken = errors.New("phone already taken")
	ErrInvalidOtpCode    = errors.New("invalid otp Code")
)

type UserRepository interface {
	GetChats(ctx context.Context, userStaticID primitive.ObjectID) ([]chat.ChatC, error)
	Create(ctx context.Context, user *user.User) (*user.User, error)
	FindOne(ctx context.Context, filter bson.M) (*user.User, error)
	FindMany(ctx context.Context, filter bson.M) ([]*user.User, error)
	FindByID(ctx context.Context, staticID primitive.ObjectID) (*user.User, error)
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByPhone(ctx context.Context, phone string) (*user.User, error)
}

type userRepository struct {
	usersCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) UserRepository {
	return &userRepository{db.Collection(database.UsersCollection)}
}

// Using mongodb aggregation pipeline to fetch user-chats.
// This process is a little special because we do not fetch all of the messages because it's really heavy query!
// Then, only last message are going to be fetched by pipeline.
func (repo *userRepository) GetChats(ctx context.Context, userStaticID primitive.ObjectID) ([]chat.ChatC, error) {
	// FIXME
	return []chat.ChatC{}, nil
	// ctx := context.TODO()

	// pipeline := []bson.M{
	// 	{
	// 		"$match": bson.M{
	// 			"$or": []interface{}{
	// 				bson.M{"chat_detail.sides": bson.M{"$in": []interface{}{userStaticID}}},
	// 				bson.M{"chat_detail.members": bson.M{"$in": []interface{}{userStaticID}}},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		"$lookup": bson.M{
	// 			"from":         "users",
	// 			"localField":   "chat_detail.sides",
	// 			"foreignField": "id",
	// 			"as":           "chat_detail.fetchedUsers",
	// 			"pipeline": []interface{}{
	// 				bson.M{
	// 					"$match": bson.M{
	// 						"$expr": bson.M{"$eq": []interface{}{"$chat_type", "direct"}},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		"$lookup": bson.M{
	// 			"from":         "messages",
	// 			"localField":   "chat_id",
	// 			"foreignField": "chat_id",
	// 			"as":           "chatMessages",
	// 		},
	// 	},
	// 	{
	// 		"$unwind": "$chatMessages",
	// 	},
	// 	{
	// 		"$project": bson.M{
	// 			"chat_id":     1,
	// 			"chat_type":   1,
	// 			"chat_detail": 1,
	// 			"messages":    bson.M{"$slice": []interface{}{"$chatMessages.messages", -1}},
	// 		},
	// 	},
	// }

	// cursor, err := repo.chatsCollection.Aggregate(ctx, pipeline)
	// if err != nil {
	// 	return []chat.ChatC{}, err
	// }
	// defer cursor.Close(ctx)

	// var chatsList []chat.ChatC
	// if err := cursor.All(ctx, &chatsList); err != nil {
	// 	return nil, err
	// }

	// return chatsList, nil
}

func (repo *userRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	_, err := repo.usersCollection.InsertOne(context.TODO(), user)
	if database.IsDuplicateKeyError(err) {
		return nil, ErrPhoneAlreadyTaken
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) FindOne(ctx context.Context, filter bson.M) (*user.User, error) {
	var model *user.User

	result := repo.usersCollection.FindOne(context.TODO(), filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	err := result.Decode(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (repo *userRepository) FindMany(ctx context.Context, filter bson.M) ([]*user.User, error) {
	cursor, err := repo.usersCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []*user.User

	err = cursor.All(ctx, &users)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *userRepository) FindByID(ctx context.Context, staticID primitive.ObjectID) (*user.User, error) {
	filter := bson.M{"id": staticID}
	return repo.FindOne(ctx, filter)
}

func (repo *userRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	filter := bson.M{"username": username}
	return repo.FindOne(ctx, filter)
}

func (repo *userRepository) FindByPhone(ctx context.Context, phone string) (*user.User, error) {
	filter := bson.M{"phone": phone}
	return repo.FindOne(ctx, filter)
}
