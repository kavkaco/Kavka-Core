package repository

import (
	"context"
	"fmt"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"

	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ChatTestSuite struct {
	suite.Suite
	userRepo repository.UserRepository
	repo     repository.ChatRepository

	userID               model.UserID
	createdChannelChatID model.ChatID
	createdGroupChatID   model.ChatID
	createdDirectChatID  model.ChatID
	recipientUserID      model.UserID
}

func (s *ChatTestSuite) SetupSuite() {
	s.repo = repository_mongo.NewChatMongoRepository(db)
	s.userRepo = repository_mongo.NewUserMongoRepository(db)

	s.userID = fmt.Sprintf("%d", random.GenerateUserID())
	s.recipientUserID = fmt.Sprintf("%d", random.GenerateUserID())
}

func (s *ChatTestSuite) TestA_CreateChannel() {
	ctx := context.TODO()

	chatDetail := model.ChannelChatDetail{
		Title:       "Test Channel",
		Username:    "TestChannelUsername",
		Members:     []model.UserID{s.userID},
		Admins:      []model.UserID{s.userID},
		Owner:       s.userID,
		Description: "Test Channel Description",
	}
	chatModel := model.NewChat(model.TypeChannel, chatDetail)

	saved, err := s.repo.Create(ctx, *chatModel)
	require.NoError(s.T(), err)

	savedChatDetail, err := utils.TypeConverter[model.ChannelChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeChannel)
	require.Equal(s.T(), chatDetail.Title, savedChatDetail.Title)
	require.Equal(s.T(), chatDetail.Username, savedChatDetail.Username)
	require.Equal(s.T(), chatDetail.Members, savedChatDetail.Members)
	require.Equal(s.T(), chatDetail.Admins, savedChatDetail.Admins)
	require.Equal(s.T(), chatDetail.Owner, savedChatDetail.Owner)
	require.Equal(s.T(), chatDetail.Description, savedChatDetail.Description)

	s.createdChannelChatID = saved.ChatID
}

func (s *ChatTestSuite) TestB_CreateGroup() {
	ctx := context.TODO()

	chatDetail := model.GroupChatDetail{
		Title:       "Test Group",
		Username:    "TestGroupUsername",
		Members:     []model.UserID{s.userID},
		Admins:      []model.UserID{s.userID},
		Owner:       s.userID,
		Description: "Test Group Description",
	}
	chatModel := model.NewChat(model.TypeGroup, chatDetail)

	saved, err := s.repo.Create(ctx, *chatModel)
	require.NoError(s.T(), err)

	savedChatDetail, err := utils.TypeConverter[model.GroupChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeGroup)
	require.Equal(s.T(), chatDetail.Title, savedChatDetail.Title)
	require.Equal(s.T(), chatDetail.Username, savedChatDetail.Username)
	require.Equal(s.T(), chatDetail.Members, savedChatDetail.Members)
	require.Equal(s.T(), chatDetail.Admins, savedChatDetail.Admins)
	require.Equal(s.T(), chatDetail.Owner, savedChatDetail.Owner)
	require.Equal(s.T(), chatDetail.Description, savedChatDetail.Description)

	s.createdGroupChatID = saved.ChatID
}

func (s *ChatTestSuite) TestC_CreateDirect() {
	ctx := context.TODO()

	chatDetail := model.DirectChatDetail{
		UserID:          s.userID,
		RecipientUserID: s.recipientUserID,
	}
	chatModel := model.NewChat(model.TypeDirect, chatDetail)

	saved, err := s.repo.Create(ctx, *chatModel)
	require.NoError(s.T(), err)

	savedChatDetail, err := utils.TypeConverter[model.DirectChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeDirect)
	require.True(s.T(), savedChatDetail.HasSide(s.userID))
	require.True(s.T(), savedChatDetail.HasSide(s.recipientUserID))

	s.createdDirectChatID = saved.ChatID
}

func (s *ChatTestSuite) TestD_FindBySides() {
	ctx := context.TODO()

	chat, err := s.repo.FindBySides(ctx, s.recipientUserID, s.userID)
	require.NoError(s.T(), err)

	chatDetail, err := utils.TypeConverter[model.DirectChatDetail](chat.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), chat.ChatType, model.TypeDirect)
	require.Equal(s.T(), chat.ChatID, s.createdDirectChatID)
	require.True(s.T(), chatDetail.HasSide(s.userID))
	require.True(s.T(), chatDetail.HasSide(s.recipientUserID))
}

func (s *ChatTestSuite) TestE_GetChat() {
	ctx := context.TODO()

	chat, err := s.repo.GetChat(ctx, s.createdDirectChatID)
	require.NoError(s.T(), err)

	chatDetail, err := utils.TypeConverter[model.DirectChatDetail](chat.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), chat.ChatType, model.TypeDirect)
	require.Equal(s.T(), chat.ChatID, s.createdDirectChatID)
	require.True(s.T(), chatDetail.HasSide(s.userID))
	require.True(s.T(), chatDetail.HasSide(s.recipientUserID))
}

func (s *ChatTestSuite) TestGetUserChats() {
	ctx := context.TODO()

	chatIDs := []model.ChatID{s.createdChannelChatID, s.createdGroupChatID, s.createdDirectChatID}
	chats, err := s.repo.GetUserChats(ctx, chatIDs)
	require.NoError(s.T(), err)

	require.Len(s.T(), chats, 3)
}

func TestChatSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ChatTestSuite))
}
