package service

import (
	"context"
	"fmt"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	service "github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
	service service.MessageService

	chatID          model.ChatID
	userID          model.UserID
	recipientUserID model.UserID
	savedMessageID  model.MessageID
}

func (s *MessageTestSuite) SetupSuite() {
	ctx := context.TODO()

	chatRepo := repository_mongo.NewChatMongoRepository(db)
	messageRepo := repository_mongo.NewMessageMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	s.service = service.NewMessageService(nil, messageRepo, chatRepo, userRepo, nil)

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

func (s *MessageTestSuite) TestA_SendTextMessage() {
	ctx := context.TODO()

	messageContent := "Hello from kavka's integration tests"
	messageGetter, varror := s.service.SendTextMessage(ctx, s.chatID, s.userID, messageContent)
	require.Nil(s.T(), varror)

	savedMessageContent, err := utils.TypeConverter[model.TextMessage](messageGetter.Message.Content)
	require.NoError(s.T(), err)

	require.Equal(s.T(), messageGetter.Message, s.userID)
	require.Equal(s.T(), savedMessageContent.Text, messageContent)

	s.savedMessageID = messageGetter.Message.MessageID
}

func (s *MessageTestSuite) TestB_DeleteMessage() {
	ctx := context.TODO()

	varror := s.service.DeleteMessage(ctx, s.chatID, s.userID, s.savedMessageID)
	require.Nil(s.T(), varror)
}

func TestMessageSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(MessageTestSuite))
}
