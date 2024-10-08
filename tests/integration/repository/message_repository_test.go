package repository

import (
	"context"
	"fmt"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
	repo repository.MessageRepository

	chatID         model.ChatID
	senderID       model.UserID
	savedMessageID model.MessageID
}

func (s *MessageTestSuite) SetupSuite() {
	ctx := context.TODO()

	chatRepo := repository_mongo.NewChatMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	s.repo = repository_mongo.NewMessageMongoRepository(db)

	user, err := userRepo.Create(ctx, &model.User{
		UserID:       fmt.Sprintf("%d", random.GenerateUserID()),
		Name:         "User2:Name",
		LastName:     "User2:LastName",
		Email:        "user2@kavka.org",
		Username:     "user2",
		Biography:    "User2:biography",
		ChatsListIDs: []model.ChatID{},
	})
	require.NoError(s.T(), err)

	s.senderID = user.UserID

	// Create a sample chat
	chatModel := model.NewChat(model.TypeChannel, model.ChannelChatDetail{
		Title:       "Chat 1",
		Username:    "Username_1",
		Description: "Test Chat 1",
		Owner:       s.senderID,
		Members:     []model.UserID{s.senderID},
		Admins:      []model.UserID{s.senderID},
	})
	chat, err := chatRepo.Create(ctx, *chatModel)
	require.NoError(s.T(), err)

	// Create message store for our chat
	err = s.repo.Create(ctx, chat.ChatID)
	require.NoError(s.T(), err)

	s.chatID = chat.ChatID
}

func (s *MessageTestSuite) TestA_InsertTextMessage() {
	ctx := context.TODO()

	messageContentModel := model.TextMessage{Text: "Text message"}
	messageModel := model.NewMessage(model.TypeTextMessage, messageContentModel, s.senderID)
	saved, err := s.repo.Insert(ctx, s.chatID, messageModel)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.Content, messageContentModel)
	require.Equal(s.T(), saved.SenderID, s.senderID)
	require.Equal(s.T(), saved.Type, messageModel.Type)

	s.savedMessageID = saved.MessageID
}

func (s *MessageTestSuite) TestB_FetchMessages() {
	ctx := context.TODO()

	messages, err := s.repo.FetchMessages(ctx, s.chatID)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), messages)
	require.Len(s.T(), messages, 1)
	require.Equal(s.T(), messages[0].Message.MessageID, s.savedMessageID)
	require.Equal(s.T(), messages[0].Sender.UserID, s.senderID)
}

func (s *MessageTestSuite) TestC_FetchLastMessage() {
	ctx := context.TODO()

	lastMessage, err := s.repo.FetchLastMessage(ctx, s.chatID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), lastMessage.MessageID, s.savedMessageID)
}

func (s *MessageTestSuite) TestC_FetchMessage() {
	ctx := context.TODO()

	message, err := s.repo.FetchMessage(ctx, s.chatID, s.savedMessageID)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), message)
	require.Equal(s.T(), message.MessageID, s.savedMessageID)
	// require.Equal(s.T(), message.Message.SenderID, s.senderID)
}

// func (s *MessageTestSuite) TestD_UpdateTextMessage() {
// 	ctx := context.TODO()

// 	newMessageContent := "Test message updated"
// 	err := s.repo.UpdateMessageContent(ctx, s.chatID, s.savedMessageID, newMessageContent)
// 	require.NoError(s.T(), err)

// 	// Fetch message from chat
// 	messages, err := s.repo.FetchMessages(ctx, s.chatID)
// 	require.NoError(s.T(), err)

// 	lastMessageContent, err := utils.TypeConverter[model.TextMessage](messages[0].Content)
// 	require.NoError(s.T(), err)

// 	updatedMessageContent := lastMessageContent.Text

// 	require.Equal(s.T(), newMessageContent, updatedMessageContent)
// }

// func (s *MessageTestSuite) TestE_DeleteMessage() {
// 	ctx := context.TODO()

// 	err := s.repo.Delete(ctx, s.chatID, s.savedMessageID)
// 	require.NoError(s.T(), err)

// 	// Fetch message from chat
// 	messages, err := s.repo.FetchMessages(ctx, s.chatID)
// 	require.NoError(s.T(), err)

// 	require.Len(s.T(), messages, 0)
// }

func TestMessageSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(MessageTestSuite))
}
