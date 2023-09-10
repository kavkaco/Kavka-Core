package repository

import (
	"Kavka/config"
	"Kavka/database"
	"Kavka/internal/domain/chat"
	"Kavka/internal/domain/message"
	repository "Kavka/internal/repository/chat"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var StaticID = primitive.NewObjectID()

type MyTestSuite struct {
	suite.Suite
	messageRepo *MessageRepository
	chatRepo    *repository.ChatRepository
	savedChat   *chat.Chat
	savedMsg    *message.Message
}

func (s *MyTestSuite) SetupSuite() {
	// Load configs
	configs := config.Read()

	configs.Mongo.DBName = "test"

	mongoClient, connErr := database.GetMongoDBInstance(configs.Mongo)
	if connErr != nil {
		panic(connErr)
	}

	s.messageRepo = NewMessageRepository(mongoClient)
	s.chatRepo = repository.NewChatRepository(mongoClient)
}

func (s *MyTestSuite) TestA_CreateChat() {
	// Creating a sample channel
	chat, saveChatErr := s.chatRepo.Create(chat.ChatTypeChannel, chat.ChannelChatDetail{
		Members: []*primitive.ObjectID{&StaticID},
		Admins:  []*primitive.ObjectID{&StaticID},
	})

	s.T().Log("ChatID:", chat.ChatID)

	assert.NoError(s.T(), saveChatErr)

	s.savedChat = chat
}

func (s *MyTestSuite) TestB_Insert() {
	msg := message.NewMessage(StaticID, message.TypeTextMessage, message.TextMessage{
		Message: "Hello World",
	})

	msg2 := message.NewMessage(StaticID, message.TypeTextMessage, message.TextMessage{
		Message: "Hello World 2",
	})

	msg3 := message.NewMessage(StaticID, message.TypeTextMessage, message.TextMessage{
		Message: "Hello World 3",
	})

	savedMsg, saveErr := s.messageRepo.Insert(s.savedChat.ChatID, msg)
	s.messageRepo.Insert(s.savedChat.ChatID, msg2)
	s.messageRepo.Insert(s.savedChat.ChatID, msg3)

	s.savedMsg = savedMsg

	assert.NoError(s.T(), saveErr)

	assert.Equal(
		s.T(),
		savedMsg.Content.(message.TextMessage).Message,
		msg.Content.(message.TextMessage).Message,
	)

	// Update chat
	chat, _ := s.chatRepo.FindByID(s.savedChat.ChatID)
	assert.Equal(s.T(), len(chat.Messages), 3)
}

// func (s *MyTestSuite) TestC_Update() {
// 	update := bson.M{"messages.$.seen": true}
// 	err := s.messageRepo.Update(s.savedChat.ChatID, s.savedMsg.MessageID, update)

// 	assert.NoError(s.T(), err)

// 	// Update chat
// 	chat, _ := s.chatRepo.FindByID(s.savedChat.ChatID)
// 	s.T().Log(chat)
// 	// seen := chat.GetMessageByID(s.savedMsg.MessageID).Seen
// 	// assert.Equal(s.T(), seen, true)
// }

func (s *MyTestSuite) TestD_Delete() {
	err := s.messageRepo.Delete(s.savedChat.ChatID, s.savedMsg.MessageID)

	assert.NoError(s.T(), err)
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
