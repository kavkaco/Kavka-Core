package repository

import (
	"context"
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/message"
	repository "github.com/kavkaco/Kavka-Core/internal/repository/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MyTestSuite struct {
	suite.Suite
	db              *mongo.Database
	messageRepo     message.Repository
	chatRepo        chat.Repository
	senderStaticID  primitive.ObjectID
	channelChatID   primitive.ObjectID
	sampleMessageID primitive.ObjectID
}

func (s *MyTestSuite) SetupSuite() {
	// Connecting to test database!
	cfg := config.Read()
	cfg.Mongo.DBName = "test"
	db, connErr := database.GetMongoDBInstance(cfg.Mongo)
	assert.NoError(s.T(), connErr)
	s.db = db

	// Drop test db
	err := s.db.Drop(context.TODO())
	assert.NoError(s.T(), err)

	s.messageRepo = NewMessageRepository(db)
	s.chatRepo = repository.NewRepository(db)

	s.senderStaticID = primitive.NewObjectID()

	newChannelChat := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
		Title:       "New Channel",
		Username:    "sample_channel",
		Description: "This is a new channel created from unit-test.",
		Members:     []primitive.ObjectID{s.senderStaticID},
		Admins:      []primitive.ObjectID{s.senderStaticID},
	})

	newChannelChat, createErr := s.chatRepo.Create(*newChannelChat)
	assert.NoError(s.T(), createErr)

	s.channelChatID = newChannelChat.ChatID
}

func (s *MyTestSuite) TestA_InsertTextMessage() {
	chatID := s.channelChatID
	messageContent := "Hello World"
	newMessage := message.NewMessage(s.senderStaticID, message.TypeTextMessage, message.TextMessage{Message: messageContent})

	createdMessage, err := s.messageRepo.Insert(chatID, newMessage)
	assert.NoError(s.T(), err)

	// Reading the memory for the recently updated store.
	foundChat, findErr := s.chatRepo.FindByID(chatID)
	assert.NoError(s.T(), findErr)
	chatMessages := foundChat.Messages

	assert.Equal(s.T(), len(chatMessages), 1)

	// Get the created message from the store.
	assert.Equal(s.T(), createdMessage.Type, message.TypeTextMessage)
	assert.Equal(s.T(), createdMessage.Content.(message.TextMessage).Message, messageContent)
	assert.NotEmpty(s.T(), createdMessage.MessageID)

	s.sampleMessageID = createdMessage.MessageID
}

func (s *MyTestSuite) TestB_UpdateMessage() {
	// Update the fields of a text-message
	update := bson.M{
		"messages.$.seen": true,
	}

	err := s.messageRepo.Update(s.channelChatID, s.sampleMessageID, update)

	assert.NoError(s.T(), err)

	// Reading the memory for the recently updated store.
	foundChat, findErr := s.chatRepo.FindByID(s.channelChatID)
	assert.NoError(s.T(), findErr)
	chatMessages := foundChat.Messages

	createdMessage := chatMessages[0]

	assert.Equal(s.T(), createdMessage.Seen, true)
}

func (s *MyTestSuite) TestC_DeleteMessage() {
	// Delete a text-message from the first chat of the list.
	err := s.messageRepo.Delete(s.channelChatID, s.sampleMessageID)

	assert.NoError(s.T(), err)

	// Reading the memory for the recently updated store.
	foundChat, findErr := s.chatRepo.FindByID(s.channelChatID)
	assert.NoError(s.T(), findErr)
	chatMessages := foundChat.Messages

	assert.Equal(s.T(), len(chatMessages), 0)
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
