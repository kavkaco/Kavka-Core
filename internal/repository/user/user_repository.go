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

type repository struct {
	usersCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) user.UserRepository {
	return &repository{db.Collection(database.UsersCollection)}
}

func (repo *repository) Create(user *user.User) (*user.User, error) {
	_, err := repo.usersCollection.InsertOne(context.TODO(), user)
	if database.IsDuplicateKeyError(err) {
		return nil, ErrPhoneAlreadyTaken
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *repository) FindOne(filter bson.M) (*user.User, error) {
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

func (repo *repository) FindMany(filter bson.M) ([]*user.User, error) {
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

func (repo *repository) FindByID(staticID primitive.ObjectID) (*user.User, error) {
	filter := bson.M{"id": staticID}
	return repo.FindOne(filter)
}

func (repo *repository) FindByUsername(username string) (*user.User, error) {
	filter := bson.M{"username": username}
	return repo.FindOne(filter)
}

func (repo *repository) FindByPhone(phone string) (*user.User, error) {
	filter := bson.M{"phone": phone}
	return repo.FindOne(filter)
}
