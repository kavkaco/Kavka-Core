package service

import (
	"context"
	"testing"

	lorem "github.com/bozaro/golorem"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	service "github.com/kavkaco/Kavka-Core/internal/service/user"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTestSuite struct {
	suite.Suite
	service  service.UserService
	userRepo repository.UserRepository
	lem      *lorem.Lorem

	userID model.UserID
	email  string
}

func (s *UserTestSuite) SetupSuite() {
	s.lem = lorem.New()
	s.userRepo = repository_mongo.NewUserMongoRepository(db)
	s.service = service.NewUserService(s.userRepo)
}

// Check here
func (s *UserTestSuite) TestA_CreateUser() {
	ctx := context.TODO()

	chatsListIDs := []model.ChatID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}

	name := s.lem.FirstName(2)
	lastName := s.lem.LastName()
	email := s.lem.Email()
	username := s.lem.Word(1, 1)
	userModel := model.NewUser(name, lastName, email, username)
	userModel.ChatsListIDs = chatsListIDs

	saved, err := s.userRepo.Create(ctx, userModel)

	require.NoError(s.T(), err)
	require.Equal(s.T(), saved.UserID, userModel.UserID)
	require.Equal(s.T(), saved.Email, email)
	require.Equal(s.T(), saved.Username, username)
	require.Equal(s.T(), saved.Name, name)
	require.Equal(s.T(), saved.LastName, lastName)

	s.userID = saved.UserID
	s.email = email
}

func (s *UserTestSuite) TestB_UpdateProfile() {
	ctx := context.TODO()
	name := s.lem.FirstName(0)
	lastName := s.lem.LastName()
	username := s.lem.Word(1, 10)
	biography := s.lem.Word(1, 30)

	err := s.service.UpdateProfile(ctx, s.userID, name, lastName, username, biography)
	require.NoError(s.T(), err)

	// Find user again to be sure that his profile changed!
	user, err := s.userRepo.FindByUserID(ctx, s.userID)
	require.NoError(s.T(), err)

	require.NoError(s.T(), err)
	require.Equal(s.T(), user.UserID, s.userID)
	require.Equal(s.T(), user.Email, s.email)
	require.Equal(s.T(), user.Username, username)
	require.Equal(s.T(), user.Name, name)
	require.Equal(s.T(), user.LastName, lastName)
}

func TestUserSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(UserTestSuite))
}
