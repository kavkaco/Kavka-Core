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
	"go.uber.org/zap"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPhoneAlreadyTaken = errors.New("phone already taken")
	ErrInvalidOtpCode    = errors.New("invalid otp Code")
)

type userRepository struct {
	logger          *zap.Logger
	usersCollection *mongo.Collection
}

func NewRepository(logger *zap.Logger, db *mongo.Database) user.UserRepository {
	return &userRepository{logger, db.Collection(database.UsersCollection)}
}

// Using mongodb aggregation pipeline to fetch user-chats.
// This process is a little special because we do not fetch all of the messages because it's really heavy query!
// Then, only last message are going to be fetched by pipeline.
func (repo *userRepository) GetChats(userStaticID primitive.ObjectID) ([]chat.ChatC, error) {
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

func (repo *userRepository) Create(user *user.User) (*user.User, error) {
	_, err := repo.usersCollection.InsertOne(context.TODO(), user)
	if database.IsDuplicateKeyError(err) {
		return nil, ErrPhoneAlreadyTaken
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) FindOne(filter bson.M) (*user.User, error) {
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

func (repo *userRepository) FindMany(filter bson.M) ([]*user.User, error) {
	cursor, err := repo.usersCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []*user.User

	err = cursor.All(context.Background(), &users)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *userRepository) FindByID(staticID primitive.ObjectID) (*user.User, error) {
	filter := bson.M{"id": staticID}
	return repo.FindOne(filter)
}

func (repo *userRepository) FindByUsername(username string) (*user.User, error) {
	filter := bson.M{"username": username}
	return repo.FindOne(filter)
}

func (repo *userRepository) FindByPhone(phone string) (*user.User, error) {
	filter := bson.M{"phone": phone}
	return repo.FindOne(filter)
}
