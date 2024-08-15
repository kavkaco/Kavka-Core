package service

import (
	"context"
	"fmt"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	service "github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ChatTestSuite struct {
	suite.Suite
	userRepo repository.UserRepository
	chatRepo repository.ChatRepository
	service  service.ChatService

	// Created chats
	createdChannelChatID model.ChatID
	createdGroupChatID   model.ChatID
	createdDirectChatID  model.ChatID

	users [2]model.User
}

func (s *ChatTestSuite) SetupSuite() {
	ctx := context.TODO()

	chatRepo := repository_mongo.NewChatMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	messageRepo := repository_mongo.NewMessageMongoRepository(db)

	s.userRepo = userRepo
	s.chatRepo = chatRepo
	s.service = service.NewChatService(nil, chatRepo, userRepo, messageRepo, nil)

	s.users = [2]model.User{
		{
			UserID:       fmt.Sprintf("%d", random.GenerateUserID()),
			Name:         "User2:Name",
			LastName:     "User2:LastName",
			Email:        "user2@kavka.org",
			Username:     "user2",
			Biography:    "User2:biography",
			ChatsListIDs: []model.ChatID{},
		},
		{
			UserID:       fmt.Sprintf("%d", random.GenerateUserID()),
			Name:         "User3:Name",
			LastName:     "User3:LastName",
			Email:        "user3@kavka.org",
			Username:     "user3",
			Biography:    "User3:biography",
			ChatsListIDs: []model.ChatID{},
		},
	}

	_, err := userRepo.Create(ctx, &s.users[0])
	require.NoError(s.T(), err)

	_, err = userRepo.Create(ctx, &s.users[1])
	require.NoError(s.T(), err)
}

func (s *ChatTestSuite) TestCreateChannel() {
	ctx := context.TODO()

	detailModel := model.ChannelChatDetail{
		Title:       "Channel1",
		Username:    "channel1",
		Owner:       s.users[0].UserID,
		Members:     []model.UserID{s.users[0].UserID},
		Admins:      []model.UserID{s.users[0].UserID},
		Description: "Channel1:Description",
	}

	saved, varror := s.service.CreateChannel(ctx, detailModel.Owner, detailModel.Title, detailModel.Username, detailModel.Description)
	require.Nil(s.T(), varror)

	chatDetail, err := utils.TypeConverter[model.ChannelChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeChannel)
	require.Equal(s.T(), chatDetail.Title, detailModel.Title)
	require.Equal(s.T(), chatDetail.Username, detailModel.Username)
	require.Equal(s.T(), chatDetail.Members, detailModel.Members)
	require.Equal(s.T(), chatDetail.Admins, detailModel.Admins)
	require.Equal(s.T(), chatDetail.Owner, detailModel.Owner)
	require.Equal(s.T(), chatDetail.Description, detailModel.Description)

	s.createdChannelChatID = saved.ChatID
}

func (s *ChatTestSuite) TestCreateGroup() {
	ctx := context.TODO()

	detailModel := model.GroupChatDetail{
		Title:       "Group1",
		Username:    "Group1",
		Owner:       s.users[0].UserID,
		Members:     []model.UserID{s.users[0].UserID},
		Admins:      []model.UserID{s.users[0].UserID},
		Description: "Group1:Description",
	}

	saved, varror := s.service.CreateGroup(ctx, detailModel.Owner, detailModel.Title, detailModel.Username, detailModel.Description)
	require.Nil(s.T(), varror)

	chatDetail, err := utils.TypeConverter[model.GroupChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.ChatType, model.TypeGroup)
	require.Equal(s.T(), chatDetail.Title, detailModel.Title)
	require.Equal(s.T(), chatDetail.Username, detailModel.Username)
	require.Equal(s.T(), chatDetail.Members, detailModel.Members)
	require.Equal(s.T(), chatDetail.Admins, detailModel.Admins)
	require.Equal(s.T(), chatDetail.Owner, detailModel.Owner)
	require.Equal(s.T(), chatDetail.Description, detailModel.Description)

	s.createdGroupChatID = saved.ChatID
}

func (s *ChatTestSuite) TestCreateDirect() {
	ctx := context.TODO()

	detailModel := &model.DirectChatDetail{
		Sides: [2]model.UserID{s.users[0].UserID, s.users[1].UserID},
	}

	saved, varror := s.service.CreateDirect(ctx, detailModel.Sides[0], detailModel.Sides[1])
	require.Nil(s.T(), varror)

	chatDetail, err := utils.TypeConverter[model.DirectChatDetail](saved.ChatDetail)
	require.NoError(s.T(), err)

	require.True(s.T(), chatDetail.HasSide(detailModel.Sides[0]))
	require.True(s.T(), chatDetail.HasSide(detailModel.Sides[1]))
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

func (s *ChatTestSuite) TestGetUserChats() {
	ctx := context.TODO()

	userModel := &s.users[0]

	userModel.ChatsListIDs = []model.ChatID{
		s.createdChannelChatID,
		s.createdGroupChatID,
		s.createdDirectChatID,
	}

	userChatsList, varror := s.service.GetUserChats(ctx, userModel.UserID)
	require.Nil(s.T(), varror)

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

func (s *ChatTestSuite) TestJoinChat() {
	ctx := context.TODO()

	// Create a plain channel chat
	detailModel := model.ChannelChatDetail{
		Title:       "Channel3",
		Username:    "channel3",
		Owner:       s.users[0].UserID,
		Members:     []model.UserID{},
		Admins:      []model.UserID{},
		Description: "Channel3:Description",
	}
	channelChat, err := s.chatRepo.Create(ctx, model.Chat{
		ChatID:     model.NewChatID(),
		ChatType:   model.TypeChannel,
		ChatDetail: detailModel,
	})
	require.NoError(s.T(), err)

	userID := s.users[1].UserID

	joinResult, varror := s.service.JoinChat(ctx, channelChat.ChatID, userID)
	if varror != nil {
		s.T().Log(varror.Error)
	}
	require.Nil(s.T(), varror)
	require.True(s.T(), joinResult.Joined)
	require.NotEmpty(s.T(), joinResult.UpdatedChat)
}

func TestChatSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ChatTestSuite))
}
