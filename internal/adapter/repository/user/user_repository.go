package repository

import (
	"Kavka/database"
	"Kavka/internal/domain/user"
	"Kavka/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPhoneAlreadyTaken = errors.New("phone already taken")
)

type UserRepository struct {
	usersCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db.Collection(database.UsersCollection),
	}
}

func (repo *UserRepository) Create(u *user.CreateUserData) (*user.User, error) {
	result, err := repo.usersCollection.InsertOne(context.TODO(), u)
	if database.IsDuplicateKeyError(err) {
		return nil, ErrPhoneAlreadyTaken
	} else if err != nil {
		return nil, err
	}

	user, convertErr := utils.TypeConverter[user.User](&u)
	if convertErr != nil {
		return nil, convertErr
	}

	user.StaticID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (repo *UserRepository) Where(filter any) ([]*user.User, error) {
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

func (repo *UserRepository) FindByID(staticID primitive.ObjectID) (*user.User, error) {
	filter := bson.D{{Key: "_id", Value: staticID}}
	result, err := repo.Where(filter)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		user := result[len(result)-1]

		return user, nil
	}

	return nil, ErrUserNotFound
}
