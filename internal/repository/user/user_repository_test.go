package repository

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const Phone = "user-phone"

type MyTestSuite struct {
	suite.Suite
	userRepo   user.UserRepository
	sampleUser user.User
}

func (s *MyTestSuite) SetupSuite() {
	s.userRepo = NewMockUserRepository()
}

func (s *MyTestSuite) TestA_Create() {
	var (
		name     = "John"
		lastName = "Doe"
	)
	user, err := s.userRepo.Create(name, lastName, Phone)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Name, name)
	assert.Equal(s.T(), user.LastName, lastName)
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
			name:   "empty",
			filter: bson.M{"name": "sample"},
			length: 0,
		},
		{
			name:   "find_one",
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
			name:     "empty",
			StaticID: primitive.NewObjectID(),
			exist:    false,
		},
		{
			name:     "find_one",
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

func (s *MyTestSuite) TestE_FindByPhone() {
	cases := []struct {
		name  string
		Phone string
		exist bool
	}{
		{
			name:  "empty",
			Phone: "sample",
			exist: false,
		},
		{
			name:  "found_just_one",
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
