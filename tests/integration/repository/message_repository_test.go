package repository

import (
	"context"
	"fmt"
	"testing"

	lorem "github.com/bozaro/golorem"
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
	lem  *lorem.Lorem

	chatID         model.ChatID
	senderID       model.UserID
	savedMessageID model.MessageID
}

func (s *MessageTestSuite) SetupSuite() {
	ctx := context.TODO()
	chatRepo := repository_mongo.NewChatMongoRepository(db)
	s.repo = repository_mongo.NewMessageMongoRepository(db)
	s.lem = lorem.New()
	s.senderID = fmt.Sprintf("%d", random.GenerateUserID())

	// Create a sample chat
	chatModel := model.NewChat(model.TypeChannel, model.ChannelChatDetail{
		Title:       s.lem.Word(1, 10),
		Username:    s.lem.LastName(),
		Description: s.lem.Paragraph(1, 2),
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

	messageContentModel := model.TextMessage{Text: s.lem.Sentence(1, 3)}
	messageModel := model.NewMessage(model.TypeTextMessage, messageContentModel, s.senderID)
	saved, err := s.repo.Insert(ctx, s.chatID, messageModel)
	require.NoError(s.T(), err)

	require.Equal(s.T(), saved.Content, messageContentModel)
	require.Equal(s.T(), saved.SenderID, s.senderID)
	require.Equal(s.T(), saved.Type, messageModel.Type)

	s.savedMessageID = saved.MessageID
}

// func (s *MessageTestSuite) TestB_FetchMessages() {
// 	ctx := context.TODO()

// 	// messages, err := s.repo.FetchMessages(ctx, s.chatID)
// 	// require.NoError(s.T(), err)

// }

// func (s *MessageTestSuite) TestC_FindMessage() {
// 	ctx := context.TODO()

// 	message, err := s.repo.FindMessage(ctx, s.chatID, s.savedMessageID)
// 	require.NoError(s.T(), err)

// 	require.NotEmpty(s.T(), message)
// 	require.Equal(s.T(), message.MessageID, s.savedMessageID)
// 	require.Equal(s.T(), message.SenderID, s.senderID)
// }

// func (s *MessageTestSuite) TestD_UpdateTextMessage() {
// 	ctx := context.TODO()

// 	newMessageContent := s.lem.Sentence(1, 2)
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
