package repository_mongo

import (
	"context"
	"testing"

	lorem "github.com/bozaro/golorem"
	"github.com/kavkaco/Kavka-Core/internal/model"
	userRepo "github.com/kavkaco/Kavka-Core/internal/repository/user"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTestSuite struct {
	suite.Suite
	repo      userRepo.UserRepository
	lem       *lorem.Lorem
	savedUser *model.User
}

func (s *UserTestSuite) SetupTest() {
	s.lem = lorem.New()
	s.repo = userRepo.NewRepository(db)
}

func (s *UserTestSuite) TestA_Create() {
	ctx := context.TODO()

	var chatsListIDs = []model.ChatID{
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

	s.savedUser = saved
}

func (s *UserTestSuite) TestB_FindOne() {
	ctx := context.TODO()

	found, err := s.repo.FindByEmail(ctx, s.savedUser.Email)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), found)
	require.Equal(s.T(), found.Name, s.savedUser.Name)
	require.Equal(s.T(), found.LastName, s.savedUser.LastName)
	require.Equal(s.T(), found.Email, s.savedUser.Email)
	require.Equal(s.T(), found.Username, s.savedUser.Username)
}

func (s *UserTestSuite) TestC_GetChats() {
	ctx := context.TODO()

	chatsListIDs, err := s.repo.GetChats(ctx, s.savedUser.UserID)
	require.NoError(s.T(), err)

	require.Equal(s.T(), len(chatsListIDs), len(s.savedUser.ChatsListIDs))
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
