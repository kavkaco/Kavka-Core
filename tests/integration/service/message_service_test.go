package service

import (
	"context"
	"fmt"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/model"
	service "github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageTestSuite struct {
	suite.Suite
	service *service.MessageService

	userID         model.UserID
	chatID         model.ChatID
	savedMessageID model.MessageID
}

func (s *MessageTestSuite) SetupSuite() {
	ctx := context.TODO()

	chatRepo := repository_mongo.NewChatMongoRepository(db)
	messageRepo := repository_mongo.NewMessageMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	streamPublisher, err := stream.NewStreamPublisher(natsClient)
	require.NoError(s.T(), err)

	s.service = service.NewMessageService(nil, messageRepo, chatRepo, userRepo, streamPublisher)

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

	testCases := []struct {
		chatID         primitive.ObjectID
		userID         string
		messageContent string
		Valid          bool
		Error          error
	}{
		{
			chatID:         model.NewChatID(),
			userID:         "",
			messageContent: "fail",
			Valid:          false,
		},
		{
			chatID:         s.chatID,
			userID:         s.userID,
			messageContent: "",
			Valid:          false,
		},
		{
			chatID:         model.NewChatID(),
			userID:         "invalid",
			messageContent: "fail",
			Valid:          false,
		},
		{
			chatID:         model.NewChatID(),
			userID:         s.userID,
			messageContent: "fail",
			Error:          service.ErrChatNotFound,
			Valid:          false,
		},
		{
			chatID:         s.chatID,
			userID:         s.userID,
			messageContent: "pass",
			Valid:          true,
		},
	}

	for _, tc := range testCases {
		messageGetter, varror := s.service.SendTextMessage(ctx, tc.chatID, tc.userID, tc.messageContent)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}
			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)

			savedMessageContent, err := utils.TypeConverter[model.TextMessage](messageGetter.Message.Content)
			require.NoError(s.T(), err)

			require.Equal(s.T(), messageGetter.Message.SenderID, tc.userID)
			require.Equal(s.T(), savedMessageContent.Text, tc.messageContent)

			s.savedMessageID = messageGetter.Message.MessageID
		} else {
			require.Fail(s.T(), "not specific")
		}
	}

	// messageContent := "Hello from kavka's integration tests"
	// messageGetter, varror := s.service.SendTextMessage(ctx, s.chatID, s.userID, messageContent)
	// require.Nil(s.T(), varror)

	// savedMessageContent, err := utils.TypeConverter[model.TextMessage](messageGetter.Message.Content)
	// require.NoError(s.T(), err)

	// require.Equal(s.T(), messageGetter.Message.SenderID, s.userID)
	// require.Equal(s.T(), savedMessageContent.Text, messageContent)

	// s.savedMessageID = messageGetter.Message.MessageID
}

// func (s *MessageTestSuite) TestB_DeleteMessage() {
// 	ctx := context.TODO()

// 	testCases := []struct {
// 		chatID    primitive.ObjectID
// 		userID    string
// 		messageID primitive.ObjectID
// 		Valid     bool
// 		Error     error
// 	}{
// 		{
// 			chatID:    model.NewChatID(),
// 			userID:    "",
// 			messageID: model.NewChatID(),
// 			Valid:     false,
// 		},
// 		{
// 			chatID:    s.chatID,
// 			userID:    s.userID,
// 			messageID: model.NewChatID(),
// 			Error:     service.ErrNotFound,
// 			Valid:     false,
// 		},
// 		{
// 			chatID:    model.NewChatID(),
// 			userID:    s.userID,
// 			messageID: s.savedMessageID,
// 			Error:     service.ErrChatNotFound,
// 			Valid:     false,
// 		},
// 		{
// 			chatID:    s.chatID,
// 			userID:    s.userID,
// 			messageID: s.savedMessageID,
// 			Valid:     true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		varror := s.service.DeleteMessage(ctx, tc.chatID, tc.userID, tc.messageID)
// 		if !tc.Valid {
// 			if tc.Error != nil {
// 				require.Equal(s.T(), tc.Error, varror.Error)
// 				continue
// 			}
// 			require.NotNil(s.T(), varror)
// 		} else if tc.Valid {
// 			require.Nil(s.T(), varror)
// 		} else {
// 			require.Fail(s.T(), "not specific")
// 		}
// 	}

// 	// varror := s.service.DeleteMessage(ctx, s.chatID, s.userID, s.savedMessageID)
// 	// require.Nil(s.T(), varror)
// }

func (s *MessageTestSuite) TestC_UpdateMessage() {
	ctx := context.TODO()

	defer func() {
		r := recover()
		if r == nil {
			require.Fail(s.T(), "should panic")
		}
	}()

	s.service.UpdateTextMessage(ctx, s.chatID, "hello")
}

func (s *MessageTestSuite) TestD_FetchMessages() {
	ctx := context.TODO()

	_, varror := s.service.FetchMessages(ctx, s.chatID)
	require.Nil(s.T(), varror)
}

// func (s *MessageTestSuite) TestE_EmptyContentSendMessage() {
// 	ctx := context.TODO()

// 	_, varror := s.service.SendTextMessage(ctx, s.chatID, s.userID, "")
// 	log.Println(varror)
// 	require.NotNil(s.T(), varror)
// }

// func (s *MessageTestSuite) TestF_InvalidChatIDSendMessage() {
// 	ctx := context.TODO()

// 	_, varror := s.service.SendTextMessage(ctx, model.NewChatID(), s.userID, "test")
// 	log.Println(varror)
// 	require.Equal(s.T(), varror, &vali.ValiErr{Error: service.ErrChatNotFound})
// }

// func (s *MessageTestSuite) TestG_InvalidUserIDSendMessage() {
// 	ctx := context.TODO()

// 	_, varror := s.service.SendTextMessage(ctx, s.chatID, "invalid", "test")
// 	log.Println(varror)
// 	require.NotNil(s.T(), varror)
// }

// func (s *MessageTestSuite) TestH_InvalidValuesDeleteChat() {
// 	ctx := context.TODO()

// 	varror := s.service.DeleteMessage(ctx, model.NewChatID(), "", model.NewChatID())
// 	log.Println(varror)
// 	require.NotNil(s.T(), varror)
// }

// func (s *MessageTestSuite) TestI_InvalidChatIDDeleteChat() {
// 	ctx := context.TODO()

// 	varror := s.service.DeleteMessage(ctx, model.NewChatID(), s.userID, model.NewChatID())
// 	require.Equal(s.T(), varror.Error, service.ErrChatNotFound)
// }

func TestMessageSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(MessageTestSuite))
}
