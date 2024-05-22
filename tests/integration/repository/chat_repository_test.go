package repository

import (
	"context"
	"fmt"
	"testing"

	lorem "github.com/bozaro/golorem"
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
	lem      *lorem.Lorem

	userID               model.UserID
	createdChannelChatID model.ChatID
	createdGroupChatID   model.ChatID
	createdDirectChatID  model.ChatID
	recipientUserID      model.UserID
}

func (s *ChatTestSuite) SetupSuite() {
	s.lem = lorem.New()
	s.repo = repository.NewChatRepository(db)
	s.userRepo = repository.NewUserRepository(db)

	s.userID = fmt.Sprintf("%d", random.GenerateUserID())
	s.recipientUserID = fmt.Sprintf("%d", random.GenerateUserID())
}

func (s *ChatTestSuite) TestCreateChannel() {
	ctx := context.TODO()

	chatDetail := model.ChannelChatDetail{
		Title:       s.lem.Word(3, 6),
		Username:    s.lem.LastName(),
		Members:     []model.UserID{s.userID},
		Admins:      []model.UserID{s.userID},
		Owner:       &s.userID,
		Description: s.lem.Paragraph(1, 4),
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

func (s *ChatTestSuite) TestCreateGroup() {
	ctx := context.TODO()

	chatDetail := model.GroupChatDetail{
		Title:       s.lem.Word(3, 6),
		Username:    s.lem.LastName(),
		Members:     []model.UserID{s.userID},
		Admins:      []model.UserID{s.userID},
		Owner:       &s.userID,
		Description: s.lem.Paragraph(1, 4),
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

func (s *ChatTestSuite) TestCreateDirect() {
	ctx := context.TODO()

	chatDetail := model.DirectChatDetail{
		Sides: [2]model.UserID{s.userID, s.recipientUserID},
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

func (s *ChatTestSuite) TestFindBySides() {
	ctx := context.TODO()

	sides := [2]model.UserID{s.userID, s.recipientUserID}
	chat, err := s.repo.FindBySides(ctx, sides)
	require.NoError(s.T(), err)

	chatDetail, err := utils.TypeConverter[model.DirectChatDetail](chat.ChatDetail)
	require.NoError(s.T(), err)

	require.True(s.T(), chatDetail.HasSide(s.userID))
	require.True(s.T(), chatDetail.HasSide(s.recipientUserID))
	require.Equal(s.T(), chat.ChatType, model.TypeDirect)
	require.Equal(s.T(), chat.ChatID, s.createdDirectChatID)
}

func (s *ChatTestSuite) TestUpdateChatLastMessage() {
	ctx := context.TODO()

	// Create message model
	messageContent := model.TextMessage{Data: "Sample message..."}
	model.NewMessage(model.TypeTextMessage, messageContent, s.userID)

	// Create last message model and update it in chat repository
	lastMessageModel := model.NewLastMessage(model.TypeTextMessage, messageContent.Data)
	err := s.repo.UpdateChatLastMessage(ctx, s.createdDirectChatID, *lastMessageModel)
	require.NoError(s.T(), err)

	// Get the chat to be sure thats changed
	chat, err := s.repo.FindByID(ctx, s.createdDirectChatID)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), chat.LastMessage)
	require.Equal(s.T(), chat.LastMessage.MessageCaption, lastMessageModel.MessageCaption)
	require.Equal(s.T(), chat.LastMessage.MessageType, lastMessageModel.MessageType)
}

func (s *ChatTestSuite) TestFindMany() {
	ctx := context.TODO()

	chatIDs := []model.ChatID{s.createdChannelChatID, s.createdGroupChatID, s.createdDirectChatID}
	chats, err := s.repo.FindMany(ctx, chatIDs)
	require.NoError(s.T(), err)

	require.Len(s.T(), chats, 3)
}

func TestChatSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ChatTestSuite))
}
