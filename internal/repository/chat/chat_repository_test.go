package repository

import (
	"Kavka/config"
	"Kavka/database"
	"Kavka/internal/domain/chat"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var StaticID = primitive.NewObjectID()
var SampleUsername = "sample"

type MyTestSuite struct {
	suite.Suite
	chatRepo  *ChatRepository
	savedChat *chat.Chat
}

func (s *MyTestSuite) SetupSuite() {
	// Load configs
	configs := config.Read()

	configs.Mongo.DBName = "test"

	mongoClient, connErr := database.GetMongoDBInstance(configs.Mongo)
	if connErr != nil {
		panic(connErr)
	}

	s.chatRepo = NewChatRepository(mongoClient)
}

func (s *MyTestSuite) TestA_Create() {
	savedChat, saveErr := s.chatRepo.Create(chat.ChatTypeChannel, chat.ChannelChatDetail{
		Members:  []*primitive.ObjectID{&StaticID},
		Admins:   []*primitive.ObjectID{&StaticID},
		Username: "sample",
	})

	assert.NoError(s.T(), saveErr)

	assert.NotEmpty(s.T(), savedChat.ChatID, "ChatID is empty!")

	s.savedChat = savedChat
}

func (s *MyTestSuite) TestB_FindByID() {
	chat, err := s.chatRepo.FindByID(s.savedChat.ChatID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), chat.ChatID, s.savedChat.ChatID)
}
func (s *MyTestSuite) TestC_FindByUsername() {
	chat, err := s.chatRepo.FindByUsername(SampleUsername)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), chat.ChatID, s.savedChat.ChatID)
}

func (s *MyTestSuite) TestD_Destroy() {
	// Finally we have to test the Destroy created chat

	err := s.chatRepo.Destroy(s.savedChat.ChatID)
	assert.NoError(s.T(), err)
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
