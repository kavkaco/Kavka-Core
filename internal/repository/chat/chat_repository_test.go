package repository

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const SampleChatUsername = "sample_chat"

type MyTestSuite struct {
	suite.Suite
	chatRepo              chat.Repository
	sampleChannelChat     chat.Chat
	sampleDirectChat      chat.Chat
	sampleDirectChatSides [2]primitive.ObjectID
	creatorID             primitive.ObjectID
}

func (s *MyTestSuite) SetupSuite() {
	s.creatorID = primitive.NewObjectID()

	// Create the clients who going to chat with each other in direct-chat.
	user1StaticID := primitive.NewObjectID()
	user2StaticID := primitive.NewObjectID()
	s.sampleDirectChatSides = [2]primitive.ObjectID{user1StaticID, user2StaticID}

	s.chatRepo = NewMockRepository()
}

func (s *MyTestSuite) TestA_Create() {
	// Create a channel chat
	newChannelChat := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
		Title:       "New Channel",
		Username:    SampleChatUsername,
		Description: "This is a new channel created from unit-test.",
		Members:     []primitive.ObjectID{s.creatorID},
		Admins:      []primitive.ObjectID{s.creatorID},
	})

	newChannelChat, err := s.chatRepo.Create(*newChannelChat)
	channelChatDetail, _ := utils.TypeConverter[chat.ChannelChatDetail](newChannelChat.ChatDetail)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), channelChatDetail.Username, SampleChatUsername)
	assert.True(s.T(), newChannelChat.IsMember(s.creatorID))
	assert.True(s.T(), newChannelChat.IsAdmin(s.creatorID))

	s.sampleChannelChat = *newChannelChat

	// Create a direct chat
	newDirectChat := chat.NewChat(chat.TypeDirect, &chat.DirectChatDetail{Sides: s.sampleDirectChatSides})

	newDirectChat, err = s.chatRepo.Create(*newDirectChat)

	chatDetail, _ := utils.TypeConverter[chat.DirectChatDetail](newDirectChat.ChatDetail)
	assert.Equal(s.T(), chatDetail.Sides, s.sampleDirectChatSides)

	assert.NoError(s.T(), err)

	s.sampleDirectChat = *newDirectChat
}

func (s *MyTestSuite) TestB_FindByID() {
	cases := []struct {
		name     string
		staticID primitive.ObjectID
		success  bool
	}{
		{
			name:     "Should rise chat not found error",
			staticID: primitive.NilObjectID,
			success:  false,
		},
		{
			name:     "Should be found",
			staticID: s.sampleChannelChat.ChatID,
			success:  true,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			foundChat, err := s.chatRepo.FindByID(tt.staticID)

			if tt.success {
				assert.NoError(s.T(), err)
				assert.NotEqual(s.T(), foundChat, nil)
			} else {
				assert.Error(s.T(), err)
			}
		})
	}
}

func (s *MyTestSuite) TestC_FindChatOrSidesByStaticID() {
	cases := []struct {
		name     string
		success  bool
		staticID primitive.ObjectID
	}{
		{
			name:     "Find the direct chat by side",
			staticID: s.sampleDirectChatSides[0],
			success:  true,
		},
		{
			name:     "Find the direct chat by staticID",
			staticID: s.sampleDirectChat.ChatID,
			success:  true,
		},
		{
			name:     "Find the channel chat by staticID",
			staticID: s.sampleChannelChat.ChatID,
			success:  true,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			user, err := s.chatRepo.FindChatOrSidesByStaticID(tt.staticID)

			if tt.success {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), user)
			} else {
				assert.Empty(s.T(), user)
			}
		})
	}
}

func (s *MyTestSuite) TestD_FindBySides() {
	cases := []struct {
		name    string
		success bool
		sides   [2]primitive.ObjectID
	}{
		{
			name:    "Must find the created chat using by the sides",
			sides:   s.sampleDirectChatSides,
			success: true,
		},
		{
			name:    "Should not find anything",
			sides:   [2]primitive.ObjectID{primitive.NilObjectID, primitive.NilObjectID},
			success: false,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			user, err := s.chatRepo.FindBySides(tt.sides)

			if tt.success {
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), user)
			} else {
				assert.Empty(s.T(), user)
			}
		})
	}
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
