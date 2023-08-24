package repository

import (
	"Kavka/config"
	"Kavka/database"
	"Kavka/internal/domain/user"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PHONE = "sample_phone_number"

type MyTestSuite struct {
	suite.Suite
	userRepo   *UserRepository
	sampleUser user.User
}

func (s *MyTestSuite) SetupSuite() {
	// Load configs
	configs := config.Read()

	configs.Mongo.DBName = "test"

	mongoClient, connErr := database.GetMongoDBInstance(configs.Mongo)
	if connErr != nil {
		panic(connErr)
	}

	mongoClient.Drop(context.TODO())

	s.userRepo = NewUserRepository(mongoClient)
}

func (s *MyTestSuite) TestCreate() {
	var (
		name     = "John"
		lastName = "Doe"
	)
	user, err := s.userRepo.Create(name, lastName, PHONE)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Name, name)
	assert.Equal(s.T(), user.LastName, lastName)
	assert.Equal(s.T(), user.Phone, PHONE)
	assert.NotEmpty(s.T(), user.StaticID)

	s.sampleUser = *user
}

func (s *MyTestSuite) TestFind() {
	cases := []struct {
		name   string
		filter bson.D
		length int
	}{
		{
			name:   "empty",
			filter: bson.D{{Key: "name", Value: "sample"}},
			length: 0,
		},
		{
			name:   "found_just_one",
			filter: bson.D{{Key: "name", Value: s.sampleUser.Name}},
			length: 1,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			users, err := s.userRepo.Where(tt.filter)

			assert.NoError(s.T(), err)
			assert.Len(s.T(), users, tt.length)
		})
	}
}

func (s *MyTestSuite) TestFindByID() {
	cases := []struct {
		name     string
		StaticID primitive.ObjectID
		exist    bool
	}{
		{
			name:     "empty",
			StaticID: primitive.NewObjectID(),
			exist:    false,
		},
		{
			name:     "found_just_one",
			StaticID: s.sampleUser.StaticID,
			exist:    true,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			user, err := s.userRepo.FindByID(tt.StaticID)

			if tt.exist {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), user)

				s.T().Log(user.FullName())
			} else {
				assert.Empty(s.T(), user)
			}
		})
	}
}

func (s *MyTestSuite) TestFindByPhone() {
	cases := []struct {
		name  string
		phone string
		exist bool
	}{
		{
			name:  "empty",
			phone: "sample",
			exist: false,
		},
		{
			name:  "found_just_one",
			phone: s.sampleUser.Phone,
			exist: true,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			user, err := s.userRepo.FindByPhone(tt.phone)

			s.T().Log(err)

			if tt.exist {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), user)

				s.T().Log(user.FullName())
			} else {
				assert.Empty(s.T(), user)
			}
		})
	}
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
