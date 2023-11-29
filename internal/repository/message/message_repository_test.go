package repository

import (
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type MyTestSuite struct {
	suite.Suite
	messageRepo    message.Repository
	chats          []*chat.Chat
	senderStaticID primitive.ObjectID
}

func (s *MyTestSuite) SetupSuite() {
	s.senderStaticID = primitive.NewObjectID()

	s.chats = []*chat.Chat{
		chat.NewChat(chat.TypeChannel, chat.ChannelChatDetail{}),
	}

	s.messageRepo = NewMockRepository(s.chats)
}

func (s *MyTestSuite) TestA_InsertTextMessage() {
	// Inserting a text-message into the first chat of the list.
	selectedChatIndex := 0
	chatID := s.chats[selectedChatIndex].ChatID
	messageContent := "Hello World"
	newMessage := message.NewMessage(s.senderStaticID, message.TypeTextMessage, message.TextMessage{Message: messageContent})

	_, err := s.messageRepo.Insert(chatID, newMessage)

	assert.NoError(s.T(), err)

	// Reading the memory for the recently updated store.
	chatMessages := s.chats[selectedChatIndex].Messages
	assert.Equal(s.T(), len(chatMessages), 1)

	// Get the created message from the store.
	createdMessage := chatMessages[0]

	assert.Equal(s.T(), createdMessage.Type, message.TypeTextMessage)
	assert.Equal(s.T(), createdMessage.Content.(message.TextMessage).Message, messageContent)
	assert.NotEmpty(s.T(), createdMessage.MessageID)
}

func (s *MyTestSuite) TestB_UpdateMessage() {
	// Update the fields of a text-message
	c := s.chats[0]
	messageID := s.chats[0].Messages[0].MessageID

	update := bson.M{"seen": true}

	err := s.messageRepo.Update(c.ChatID, messageID, update)

	assert.NoError(s.T(), err)

	// Reading the memory for the recently updated store.
	msg := s.chats[0].Messages[0]
	assert.Equal(s.T(), msg.Seen, true)
}

func (s *MyTestSuite) TestC_DeleteMessage() {
	// Delete a text-message from the first chat of the list.
	chatID := s.chats[0].ChatID
	messageID := s.chats[0].Messages[0].MessageID

	err := s.messageRepo.Delete(chatID, messageID)

	assert.NoError(s.T(), err)

	// Reading the memory for the recently updated store.
	chatMessages := s.chats[0].Messages
	assert.Equal(s.T(), len(chatMessages), 0)
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
