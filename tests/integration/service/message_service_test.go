package service

import (
	"context"
	"fmt"
	"testing"

	lorem "github.com/bozaro/golorem"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	service "github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
	service service.MessageService
	lem     *lorem.Lorem

	chatID          model.ChatID
	userID          model.UserID
	recipientUserID model.UserID
	savedMessageID  model.MessageID
}

func (s *MessageTestSuite) SetupSuite() {
	ctx := context.TODO()

	s.lem = lorem.New()

	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	s.service = service.NewMessageService(messageRepo, chatRepo)

	s.userID = fmt.Sprintf("%d", random.GenerateUserID())
	s.recipientUserID = fmt.Sprintf("%d", random.GenerateUserID())

	// Create sample channel

	chatDetail := model.DirectChatDetail{
		Sides: [2]model.UserID{s.userID, s.recipientUserID},
	}
	chatModel := model.NewChat(model.TypeDirect, chatDetail)

	chat, err := chatRepo.Create(ctx, *chatModel)
	require.NoError(s.T(), err)

	err = messageRepo.Create(ctx, chat.ChatID)
	require.NoError(s.T(), err)

	s.chatID = chat.ChatID
}

func (s *MessageTestSuite) TestInsertTextMessage() {
	ctx := context.TODO()

	messageContent := s.lem.Paragraph(1, 20)
	message, err := s.service.InsertTextMessage(ctx, s.chatID, s.userID, messageContent)
	require.NoError(s.T(), err)

	savedMessageContent, err := utils.TypeConverter[model.TextMessage](message.Content)
	require.NoError(s.T(), err)

	require.Equal(s.T(), message.SenderID, s.userID)
	require.Equal(s.T(), savedMessageContent.Data, messageContent)

	s.savedMessageID = message.MessageID
}

func (s *MessageTestSuite) TestDeleteMessage() {
	ctx := context.TODO()

	err := s.service.DeleteMessage(ctx, s.chatID, s.userID, s.savedMessageID)
	require.NoError(s.T(), err)
}

func TestMessageSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(MessageTestSuite))
}
