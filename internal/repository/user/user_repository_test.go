package repository

import (
	"context"
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const Phone = "user-phone"

type MyTestSuite struct {
	suite.Suite
	db         *mongo.Database
	userRepo   user.UserRepository
	sampleUser user.User
}

func (s *MyTestSuite) SetupSuite() {
	// Connecting to test database!
	cfg := config.Read()
	cfg.Mongo.DBName = "test"
	db, connErr := database.GetMongoDBInstance(cfg.Mongo)
	assert.NoError(s.T(), connErr)
	s.db = db

	// Drop test db
	err := s.db.Drop(context.TODO())
	assert.NoError(s.T(), err)

	s.userRepo = NewUserRepository(db)
}

func (s *MyTestSuite) TestA_Create() {
	user := user.NewUser(Phone)
	user.Name = "John"
	user.LastName = "Doe"

	user, err := s.userRepo.Create(user)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Name, user.Name)
	assert.Equal(s.T(), user.LastName, user.LastName)
	assert.Equal(s.T(), user.Phone, Phone)
	assert.NotEmpty(s.T(), user.StaticID)

	s.sampleUser = *user
}

func (s *MyTestSuite) TestB_Where() {
	users, err := s.userRepo.Where(bson.M{"phone": Phone})

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(users), 1)
}

func (s *MyTestSuite) TestC_Find() {
	cases := []struct {
		name   string
		filter bson.M
		length int
	}{
		{
			name:   "Should not find anything",
			filter: bson.M{"name": "sample"},
			length: 0,
		},
		{
			name:   "Should find the user",
			filter: bson.M{"name": s.sampleUser.Name},
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

func (s *MyTestSuite) TestD_FindByID() {
	cases := []struct {
		name     string
		StaticID primitive.ObjectID
		exist    bool
	}{
		{
			name:     "Should not find anything",
			StaticID: primitive.NewObjectID(),
			exist:    false,
		},
		{
			name:     "Should find the user",
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
			} else {
				assert.Empty(s.T(), user)
			}
		})
	}
}

func (s *MyTestSuite) TestE_FindByID() {
	cases := []struct {
		name      string
		StaticIDs []primitive.ObjectID
		exist     bool
	}{
		{
			name:      "Should not find anything",
			StaticIDs: []primitive.ObjectID{primitive.NewObjectID()},
			exist:     false,
		},
		{
			name:      "Should find the user",
			StaticIDs: []primitive.ObjectID{s.sampleUser.StaticID},
			exist:     true,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			user, err := s.userRepo.FindMany(tt.StaticIDs)

			if tt.exist {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), user)
			} else {
				assert.Empty(s.T(), user)
			}
		})
	}
}

func (s *MyTestSuite) TestF_FindByPhone() {
	cases := []struct {
		name  string
		Phone string
		exist bool
	}{
		{
			name:  "Should not find anything",
			Phone: "sample",
			exist: false,
		},
		{
			name:  "Should find the user",
			Phone: s.sampleUser.Phone,
			exist: true,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			user, err := s.userRepo.FindByPhone(tt.Phone)

			if tt.exist {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), user)
			} else {
				assert.Empty(s.T(), user)
			}
		})
	}
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
