package service

import (
	"context"
	"fmt"
	"testing"

	lorem "github.com/bozaro/golorem"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	repository_mongo "github.com/kavkaco/Kavka-Core/internal/repository/mongo"
	service "github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ChatTestSuite struct {
	suite.Suite
	userRepo repository.UserRepository
	service  service.ChatService
	lem      *lorem.Lorem

	userID               model.UserID
	createdChannelChatID model.ChatID
	createdGroupChatID   model.ChatID
	createdDirectChatID  model.ChatID
	recipientUserID      model.UserID
}

func (s *ChatTestSuite) SetupSuite() {
	s.lem = lorem.New()

	chatRepo := repository_mongo.NewChatMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)

	s.userRepo = userRepo
	s.service = service.NewChatService(chatRepo, userRepo)

	s.userID = fmt.Sprintf("%d", random.GenerateUserID())
	s.recipientUserID = fmt.Sprintf("%d", random.GenerateUserID())
}

func (s *ChatTestSuite) TestCreateChannel() {
	ctx := context.TODO()

	title := s.lem.Word(3, 6)
	username := s.lem.LastName()
	description := s.lem.Paragraph(1, 4)
	members := []model.UserID{s.userID}
	admins := []model.UserID{s.userID}
	owner := s.userID

	saved, err := s.service.CreateChannel(ctx, s.userID, title, username, description)
	require.NoError(s.T(), err)

	chatDetail, err := utils.TypeConverter[model.ChannelChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeChannel)
	require.Equal(s.T(), chatDetail.Title, title)
	require.Equal(s.T(), chatDetail.Username, username)
	require.Equal(s.T(), chatDetail.Members, members)
	require.Equal(s.T(), chatDetail.Admins, admins)
	require.Equal(s.T(), chatDetail.Owner, owner)
	require.Equal(s.T(), chatDetail.Description, description)

	s.createdChannelChatID = saved.ChatID
}

func (s *ChatTestSuite) TestCreateGroup() {
	ctx := context.TODO()

	title := s.lem.Word(3, 6)
	username := s.lem.LastName()
	description := s.lem.Paragraph(1, 4)
	members := []model.UserID{s.userID}
	admins := []model.UserID{s.userID}
	owner := s.userID

	saved, err := s.service.CreateGroup(ctx, s.userID, title, username, description)
	require.NoError(s.T(), err)

	chatDetail, err := utils.TypeConverter[model.GroupChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeGroup)
	require.Equal(s.T(), chatDetail.Title, title)
	require.Equal(s.T(), chatDetail.Username, username)
	require.Equal(s.T(), chatDetail.Members, members)
	require.Equal(s.T(), chatDetail.Admins, admins)
	require.Equal(s.T(), chatDetail.Owner, owner)
	require.Equal(s.T(), chatDetail.Description, description)

	s.createdGroupChatID = saved.ChatID
}

func (s *ChatTestSuite) TestCreateDirect() {
	ctx := context.TODO()

	saved, err := s.service.CreateDirect(ctx, s.userID, s.recipientUserID)
	require.NoError(s.T(), err)

	chatDetail, err := utils.TypeConverter[model.DirectChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.True(s.T(), chatDetail.HasSide(s.userID))
	require.True(s.T(), chatDetail.HasSide(s.recipientUserID))

	s.createdDirectChatID = saved.ChatID
}

func (s *ChatTestSuite) TestGetChat() {
	ctx := context.TODO()

	chat, err := s.service.GetChat(ctx, s.createdChannelChatID)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), chat)
	require.Equal(s.T(), chat.ChatID, s.createdChannelChatID)
}

func (s *ChatTestSuite) TestGetUserChats() {
	ctx := context.TODO()

	// Create a real user to test GetUserChats
	userModel := model.NewUser(s.lem.FirstName(1), s.lem.LastName(), s.lem.Email(), s.lem.Word(1, 10))
	userModel.ChatsListIDs = []model.ChatID{
		s.createdChannelChatID,
		s.createdGroupChatID,
		s.createdDirectChatID,
	}
	user, err := s.userRepo.Create(ctx, userModel)
	require.NoError(s.T(), err)

	userChatsList, err := s.service.GetUserChats(ctx, user.UserID)
	require.NoError(s.T(), err)

	for _, v := range userChatsList {
		switch v.ChatType {
		case model.TypeChannel:
			require.Equal(s.T(), v.ChatID, s.createdChannelChatID)
		case model.TypeGroup:
			require.Equal(s.T(), v.ChatID, s.createdGroupChatID)
		case model.TypeDirect:
			require.Equal(s.T(), v.ChatID, s.createdDirectChatID)
		}
	}
}

func TestChatSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ChatTestSuite))
}
