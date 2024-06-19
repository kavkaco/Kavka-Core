package repository

import (
	"context"
	"testing"

	lorem "github.com/bozaro/golorem"
	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTestSuite struct {
	suite.Suite
	repo      repository.UserRepository
	lem       *lorem.Lorem
	savedUser *model.User
}

func (s *UserTestSuite) SetupSuite() {
	s.lem = lorem.New()
	s.repo = repository_mongo.NewUserMongoRepository(db)
}

func (s *UserTestSuite) TestA_Create() {
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

func (s *UserTestSuite) TestD_AddToUserChats() {
	ctx := context.TODO()

	// Add channel chat to user's chats
	var chatID model.ChatID = primitive.NewObjectID()
	err := s.repo.AddToUserChats(ctx, s.savedUser.UserID, chatID)
	require.NoError(s.T(), err)

	// Let's find the user and check that chat added to list or not
	foundUser, err := s.repo.FindByUserID(ctx, s.savedUser.UserID)
	require.NoError(s.T(), err)

	require.Len(s.T(), foundUser.ChatsListIDs, len(s.savedUser.ChatsListIDs)+1)
}

func (s *UserTestSuite) TestE_Update() {
	ctx := context.TODO()

	name := s.lem.FirstName(2)
	lastName := s.lem.LastName()
	username := s.lem.Word(1, 1)
	biography := s.lem.Sentence(1, 3)

	err := s.repo.Update(ctx, s.savedUser.UserID, name, lastName, username, biography)
	require.NoError(s.T(), err)

	// Get the user again to be sure it's updated
	user, err := s.repo.FindByUserID(ctx, s.savedUser.UserID)
	require.NoError(s.T(), err)

	require.Equal(s.T(), user.Name, name)
	require.Equal(s.T(), user.LastName, lastName)
	require.Equal(s.T(), user.Username, username)
	require.Equal(s.T(), user.Biography, biography)
}
func (s *UserTestSuite) TestF_Delete() {
	ctx := context.TODO()
	err := s.repo.DeleteByID(ctx, s.savedUser.UserID)
	require.NoError(s.T(), err)
}

func TestUserSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(UserTestSuite))
}
