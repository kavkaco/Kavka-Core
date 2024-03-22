package repository

import (
	"context"
	"errors"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPhoneAlreadyTaken = errors.New("phone already taken")
	ErrInvalidOtpCode    = errors.New("invalid otp Code")
)

type userRepository struct {
	usersCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) user.UserRepository {
	return &userRepository{db.Collection(database.UsersCollection)}
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
	if result.Err() != nil {
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
	if err != nil {
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
