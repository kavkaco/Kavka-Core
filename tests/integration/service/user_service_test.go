package service

import (
	"context"

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
	service service.UserService
	repo    repository.UserRepository
	lem     *lorem.Lorem
	userId  model.UserID
}

func (s *UserTestSuite) SetupSuite() {
	s.lem = lorem.New()
	s.repo = repository_mongo.NewUserMongoRepository(db)
	s.service = service.NewUserService(s.repo)
}

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

	saved, err := s.repo.Create(ctx, userModel)

	require.NoError(s.T(), err)
	require.Equal(s.T(), saved.UserID, userModel.UserID)
	require.Equal(s.T(), saved.Email, email)
	require.Equal(s.T(), saved.Username, username)
	require.Equal(s.T(), saved.Name, name)
	require.Equal(s.T(), saved.LastName, lastName)

	s.userId = saved.UserID
}

func (s *UserTestSuite) TestB_UpdateProfile() {
	ctx := context.TODO()
	name := s.lem.FirstName(0)
	lastName := s.lem.LastName()
	username := s.lem.Word(1, 10)
	biography := s.lem.Word(1, 30)
	err := s.service.UpdateProfile(ctx, s.userId, name, lastName, username, biography)
	require.NoError(s.T(), err)
}

func (s *UserTestSuite) TestC_DeleteUser() {
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

	saved, err := s.repo.Create(ctx, userModel)

	require.NoError(s.T(), err)
	require.Equal(s.T(), saved.UserID, userModel.UserID)
	require.Equal(s.T(), saved.Email, email)
	require.Equal(s.T(), saved.Username, username)
	require.Equal(s.T(), saved.Name, name)
	require.Equal(s.T(), saved.LastName, lastName)

	s.userId = saved.UserID
}
