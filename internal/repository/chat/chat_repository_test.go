package repository

import (
	"Kavka/config"
	"Kavka/database"
	"Kavka/internal/domain/chat"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var StaticID = primitive.NewObjectID()

var SampleSides = [2]primitive.ObjectID{
	primitive.NewObjectID(),
	primitive.NewObjectID(),
}

type MyTestSuite struct {
	suite.Suite
	chatRepo        *ChatRepository
	savedChat       *chat.Chat
	savedDirectChat *chat.Chat
}

func (s *MyTestSuite) SetupSuite() {
	// Load configs
	configs := config.Read()

	configs.Mongo.DBName = "test"

	mongoClient, connErr := database.GetMongoDBInstance(configs.Mongo)
	if connErr != nil {
		panic(connErr)
	}

	mongoClient.Collection(database.ChatsCollection).Drop(context.Background())

	s.chatRepo = NewChatRepository(mongoClient)
}

func (s *MyTestSuite) TestA_CreateChannel() {
	savedChat, saveErr := s.chatRepo.Create(chat.ChatTypeChannel, chat.ChannelChatDetail{
		Members:  []*primitive.ObjectID{&StaticID},
		Admins:   []*primitive.ObjectID{&StaticID},
		Username: "sample",
	})

	assert.NoError(s.T(), saveErr)

	assert.NotEmpty(s.T(), savedChat.ChatID, "ChatID is empty!")

	s.savedChat = savedChat
}

func (s *MyTestSuite) TestB_CreateDirect() {
	savedChat, saveErr := s.chatRepo.Create(chat.ChatTypeDirect, chat.DirectChatDetail{
		Sides: [2]*primitive.ObjectID{
			&SampleSides[0],
			&SampleSides[1],
		},
	})

	assert.NoError(s.T(), saveErr)
	assert.NotEmpty(s.T(), savedChat.ChatID, "ChatID is empty!")

	s.savedDirectChat = savedChat
}

func (s *MyTestSuite) TestC_FindByID() {
	chat, err := s.chatRepo.FindByID(s.savedChat.ChatID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), chat.ChatID, s.savedChat.ChatID)
}
func (s *MyTestSuite) TestD_FindChatOrSidesByStaticID() {
	findByChatID, findByChatIDErr := s.chatRepo.FindChatOrSidesByStaticID(&s.savedDirectChat.ChatID)

	assert.NoError(s.T(), findByChatIDErr)
	assert.Equal(s.T(), findByChatID.ChatID, s.savedDirectChat.ChatID)

	findBySides, findBySidesErr := s.chatRepo.FindChatOrSidesByStaticID(&SampleSides[0])

	assert.NoError(s.T(), findBySidesErr)
	assert.Equal(s.T(), findBySides.ChatID, s.savedDirectChat.ChatID)
}

func (s *MyTestSuite) TestE_FindBySides() {
	chat, err := s.chatRepo.FindBySides([2]*primitive.ObjectID{
		&SampleSides[0],
		&SampleSides[1],
	})

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), chat.ChatID, s.savedDirectChat.ChatID)
}

func (s *MyTestSuite) TestF_Destroy() {
	// destroy only created chat for test
	err := s.chatRepo.Destroy(s.savedChat.ChatID)
	assert.NoError(s.T(), err)
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
