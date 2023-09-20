package repository

import (
	"github.com/fatih/structs"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepository struct {
	users []*user.User
}

func NewMockUserRepository() user.UserRepository {
	return &MockUserRepository{}
}

func (repo *MockUserRepository) Create(user *user.User) (*user.User, error) {
	repo.users = append(repo.users, user)

	return user, nil
}

func (repo *MockUserRepository) Where(filter bson.M) ([]*user.User, error) {
	var filterKey string
	var filterValue interface{}

	for k, v := range filter {
		filterKey = k
		filterValue = v
	}

	var foundUsers []*user.User

	if len(repo.users) == 0 {
		return foundUsers, nil
	}

	for _, user := range repo.users {
		// Check filter for user
		fields := structs.Fields(user)

		for _, field := range fields {
			tagValue := field.Tag("bson")
			fieldValue := field.Value()

			if filterKey == tagValue {
				switch filterValue := filterValue.(type) {
				case primitive.ObjectID:
					if fieldValue.(primitive.ObjectID).Hex() == filterValue.Hex() {
						foundUsers = append(foundUsers, user)
					}
				case any:
					if fieldValue == filterValue {
						foundUsers = append(foundUsers, user)
					}
				}
			}
		}
	}

	return foundUsers, nil
}

func (repo *MockUserRepository) findBy(filter bson.M) (*user.User, error) {
	users, err := repo.Where(filter)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		user := users[len(users)-1]

		return user, nil
	}

	return nil, ErrUserNotFound
}

func (repo *MockUserRepository) FindByID(staticID primitive.ObjectID) (*user.User, error) {
	filter := bson.M{"id": staticID}
	return repo.findBy(filter)
}

func (repo *MockUserRepository) FindByUsername(username string) (*user.User, error) {
	filter := bson.M{"username": username}
	return repo.findBy(filter)
}

func (repo *MockUserRepository) FindByPhone(phone string) (*user.User, error) {
	filter := bson.M{"phone": phone}
	return repo.findBy(filter)
}

func (repo *MockUserRepository) FindOrCreateGuestUser(phone string) (*user.User, error) {
	return nil, nil
}
