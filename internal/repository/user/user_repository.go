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

func (repo *userRepository) Create(name string, lastName string, phone string) (*user.User, error) {
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

func (repo *userRepository) Where(filter any) ([]*user.User, error) {
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

func (repo *userRepository) findBy(filter bson.D) (*user.User, error) {
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

func (repo *userRepository) FindByID(staticID primitive.ObjectID) (*user.User, error) {
	filter := bson.D{{Key: "_id", Value: staticID}}
	return repo.findBy(filter)
}

func (repo *userRepository) FindByUsername(username string) (*user.User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	return repo.findBy(filter)
}

func (repo *userRepository) FindByPhone(phone string) (*user.User, error) {
	filter := bson.D{{Key: "phone", Value: phone}}
	return repo.findBy(filter)
}

func (repo *userRepository) FindOrCreateGuestUser(phone string) (*user.User, error) {
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
