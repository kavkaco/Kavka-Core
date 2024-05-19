package repository

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already taken")
)

type UserRepository interface {
	GetChats(ctx context.Context, userID model.UserID) ([]model.ChatID, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindOne(ctx context.Context, filter bson.M) (*model.User, error)
	FindMany(ctx context.Context, filter bson.M) ([]*model.User, error)
	FindByUserID(ctx context.Context, userID model.UserID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	usersCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) UserRepository {
	return &userRepository{db.Collection(database.UsersCollection)}
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
		return nil, ErrEmailAlreadyTaken
	} else if err != nil {
		return nil, err
	}

	return userModel, nil
}

func (repo *userRepository) FindOne(ctx context.Context, filter bson.M) (*model.User, error) {
	var model *model.User

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

func (repo *userRepository) FindMany(ctx context.Context, filter bson.M) ([]*model.User, error) {
	cursor, err := repo.usersCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []*model.User

	err = cursor.All(ctx, &users)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return users, nil
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
