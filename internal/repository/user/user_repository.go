package repository

import (
	"context"
	"errors"

	"Kavka/database"
	"Kavka/internal/domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPhoneAlreadyTaken = errors.New("phone already taken")
	ErrInvalidOtpCode    = errors.New("invalid otp Code")
)

type UserRepository struct {
	usersCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db.Collection(database.UsersCollection),
	}
}

func (repo *UserRepository) Create(name string, lastName string, phone string) (*user.User, error) {
	user := user.NewUser(phone)

	user.Name = name
	user.LastName = lastName

	_, err := repo.usersCollection.InsertOne(context.TODO(), user)
	if database.IsDuplicateKeyError(err) {
		return nil, ErrPhoneAlreadyTaken
	} else if err != nil {
		return nil, err
	}

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

func (repo *UserRepository) findBy(filter bson.D) (*user.User, error) {
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

func (repo *UserRepository) FindByID(staticID primitive.ObjectID) (*user.User, error) {
	filter := bson.D{{Key: "_id", Value: staticID}}
	return repo.findBy(filter)
}

func (repo *UserRepository) FindByUsername(username string) (*user.User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	return repo.findBy(filter)
}

func (repo *UserRepository) FindByPhone(phone string) (*user.User, error) {
	filter := bson.D{{Key: "phone", Value: phone}}
	return repo.findBy(filter)
}

func (repo *UserRepository) FindOrCreateGuestUser(phone string) (*user.User, error) {
	foundUser, findErr := repo.FindByPhone(phone)

	if findErr != nil {
		// Creating a new user with entered Phone
		newUser, createErr := repo.Create("guest", "guest", phone)

		if createErr != nil {
			return nil, createErr
		}

		return newUser, nil
	}

	return foundUser, nil
}
