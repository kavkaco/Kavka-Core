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

	userID         model.UserID
	chatID         model.ChatID
	savedMessageID model.MessageID
}

func (s *MessageTestSuite) SetupSuite() {
	ctx := context.TODO()

	chatRepo := repository_mongo.NewChatMongoRepository(db)
	messageRepo := repository_mongo.NewMessageMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	s.service = service.NewMessageService(nil, messageRepo, chatRepo, userRepo, nil)

	user, err := userRepo.Create(ctx, &model.User{
		UserID:       fmt.Sprintf("%d", random.GenerateUserID()),
		Name:         "User4:Name",
		LastName:     "User4:LastName",
		Email:        "user4@kavka.org",
		Username:     "user4",
		Biography:    "User4:biography",
		ChatsListIDs: []model.ChatID{},
	})
	require.NoError(s.T(), err)

	chat, err := chatRepo.Create(ctx, *model.NewChat(model.TypeChannel, model.ChannelChatDetail{
		Title:       "Channel2",
		Username:    "channel2",
		Owner:       user.UserID,
		Members:     []model.UserID{user.UserID},
		Admins:      []model.UserID{user.UserID},
		Description: "Channel2:Description",
	}))
	require.NoError(s.T(), err)

	s.userID = user.UserID
	s.chatID = chat.ChatID

	err = messageRepo.Create(ctx, s.chatID)
	require.NoError(s.T(), err)
}

func (s *MessageTestSuite) TestA_SendTextMessage() {
	ctx := context.TODO()

	messageContent := "Hello from kavka's integration tests"
	messageGetter, varror := s.service.SendTextMessage(ctx, s.chatID, s.userID, messageContent)
	require.Nil(s.T(), varror)

	savedMessageContent, err := utils.TypeConverter[model.TextMessage](messageGetter.Message.Content)
	require.NoError(s.T(), err)

	require.Equal(s.T(), messageGetter.Message.SenderID, s.userID)
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
