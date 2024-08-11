package service

import (
	"context"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	service "github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ChatTestSuite struct {
	suite.Suite
	userRepo repository.UserRepository
	service  service.ChatService

	// Created chats
	createdChannelChatID model.ChatID
	createdGroupChatID   model.ChatID
	createdDirectChatID  model.ChatID
}

func (s *ChatTestSuite) SetupSuite() {
	chatRepo := repository_mongo.NewChatMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	messageRepo := repository_mongo.NewMessageMongoRepository(db)

	s.userRepo = userRepo
	s.service = service.NewChatService(nil, chatRepo, userRepo, messageRepo, nil)
}

func (s *ChatTestSuite) TestCreateChannel() {
	ctx := context.TODO()

	m := channelChatDetailTestModel

	saved, varror := s.service.CreateChannel(ctx, m.Owner, m.Title, m.Username, m.Description)
	require.Nil(s.T(), varror)

	chatDetail, err := utils.TypeConverter[model.ChannelChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeChannel)
	require.Equal(s.T(), chatDetail.Title, m.Title)
	require.Equal(s.T(), chatDetail.Username, m.Username)
	require.Equal(s.T(), chatDetail.Members, m.Members)
	require.Equal(s.T(), chatDetail.Admins, m.Admins)
	require.Equal(s.T(), chatDetail.Owner, m.Owner)
	require.Equal(s.T(), chatDetail.Description, m.Description)

	s.createdChannelChatID = saved.ChatID
}

func (s *ChatTestSuite) TestCreateGroup() {
	ctx := context.TODO()

	m := groupChatDetailTestModel

	saved, varror := s.service.CreateGroup(ctx, m.Owner, m.Title, m.Username, m.Description)
	require.Nil(s.T(), varror)

	chatDetail, err := utils.TypeConverter[model.GroupChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeGroup)
	require.Equal(s.T(), chatDetail.Title, m.Title)
	require.Equal(s.T(), chatDetail.Username, m.Username)
	require.Equal(s.T(), chatDetail.Members, m.Members)
	require.Equal(s.T(), chatDetail.Admins, m.Admins)
	require.Equal(s.T(), chatDetail.Owner, m.Owner)
	require.Equal(s.T(), chatDetail.Description, m.Description)

	s.createdGroupChatID = saved.ChatID
}

func (s *ChatTestSuite) TestCreateDirect() {
	ctx := context.TODO()

	m := directChatDetailTestModel

	saved, varror := s.service.CreateDirect(ctx, m.Sides[0], m.Sides[1])
	require.Nil(s.T(), varror)

	chatDetail, err := utils.TypeConverter[model.DirectChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.True(s.T(), chatDetail.HasSide(m.Sides[0]))
	require.True(s.T(), chatDetail.HasSide(m.Sides[1]))
	require.False(s.T(), chatDetail.HasSide("invalid-user-id"))

	s.createdDirectChatID = saved.ChatID
}

func (s *ChatTestSuite) TestGetChat_Channel() {
	ctx := context.TODO()

	chat, varror := s.service.GetChat(ctx, s.createdChannelChatID)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), chat)
	require.Equal(s.T(), chat.ChatID, s.createdChannelChatID)
}

func (s *ChatTestSuite) TestGetChat_Group() {
	ctx := context.TODO()

	chat, varror := s.service.GetChat(ctx, s.createdGroupChatID)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), chat)
	require.Equal(s.T(), chat.ChatID, s.createdGroupChatID)
}

// func (s *ChatTestSuite) TestGetUserChats() {
// 	ctx := context.TODO()

// 	// Create a real user to test GetUserChats
// 	userModel := model.NewUser(
// 		"Margaret", "Vega", "margaret_vega@kavka.org", "margaret_vega",
// 	)
// 	userModel.ChatsListIDs = []model.ChatID{
// 		s.createdChannelChatID,
// 		s.createdGroupChatID,
// 		s.createdDirectChatID,
// 	}
// 	user, err := s.userRepo.Create(ctx, userModel)
// 	require.NoError(s.T(), err)

// 	userChatsList, varror := s.service.GetUserChats(ctx, user.UserID)
// 	require.Nil(s.T(), varror)

// 	for _, v := range userChatsList {
// 		switch v.ChatType {
// 		case model.TypeChannel:
// 			require.Equal(s.T(), v.ChatID, s.createdChannelChatID)
// 		case model.TypeGroup:
// 			require.Equal(s.T(), v.ChatID, s.createdGroupChatID)
// 		case model.TypeDirect:
// 			require.Equal(s.T(), v.ChatID, s.createdDirectChatID)
// 		}
// 	}
// }

func TestChatSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ChatTestSuite))
}
