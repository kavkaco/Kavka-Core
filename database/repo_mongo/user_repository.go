package repository_mongo

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/repository"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	usersCollection *mongo.Collection
}

func NewUserMongoRepository(db *mongo.Database) repository.UserRepository {
	return &userRepository{db.Collection(database.UsersCollection)}
}

func (repo *userRepository) Update(ctx context.Context, userID string, name, lastName, username, biography string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"name":      name,
			"last_name": lastName,
			"username":  username,
			"biography": biography,
		},
	}

	result, err := repo.usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 || result.ModifiedCount == 0 {
		return repository.ErrNotModified
	}

	return nil
}

func (repo *userRepository) AddToUserChats(ctx context.Context, userID model.UserID, chatID model.ChatID) error {
	update := bson.M{
		"$addToSet": bson.M{
			"chats_list_ids": chatID,
		},
	}
	result := repo.usersCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (repo *userRepository) GetChats(ctx context.Context, userID model.UserID) ([]model.ChatID, error) {
	foundUser, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return foundUser.ChatsListIDs, nil
}

func (repo *userRepository) Create(ctx context.Context, userModel *model.User) (*model.User, error) {
	_, err := repo.usersCollection.InsertOne(context.TODO(), userModel)
	if database.IsDuplicateKeyError(err) {
		return nil, repository.ErrUniqueConstraint
	} else if err != nil {
		return nil, err
	}

	return userModel, nil
}

func (repo *userRepository) FindOne(ctx context.Context, filter bson.M) (*model.User, error) {
	var model *model.User

	result := repo.usersCollection.FindOne(context.TODO(), filter)
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

func (repo *userRepository) FindByUserID(ctx context.Context, userID model.UserID) (*model.User, error) {
	filter := bson.M{"user_id": userID}
	return repo.FindOne(ctx, filter)
}

func (repo *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	filter := bson.M{"username": username}
	return repo.FindOne(ctx, filter)
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	filter := bson.M{"email": email}
	return repo.FindOne(ctx, filter)
}

func (repo *userRepository) DeleteByID(ctx context.Context, userID model.UserID) error {
	filter := bson.M{"user_id": userID}
	result, err := repo.usersCollection.DeleteOne(ctx, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return repository.ErrNotFound
	} else if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return repository.ErrNotDeleted
	}
	return nil
}
