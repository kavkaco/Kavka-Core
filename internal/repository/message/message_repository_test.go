package repository

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/message"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMessageRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("test insert text message", func(mt *mtest.T) {
		messageRepo := NewRepository(mt.DB)

		expectedResult := []bson.E{
			{Key: "NModified", Value: 1},
			{Key: "N", Value: 1},
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(expectedResult...))

		ownerStaticID := primitive.NewObjectID()
		chatModel := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
			Title:       "MyChannel",
			Username:    "my_channel",
			Description: "Example description",
			Members:     []primitive.ObjectID{ownerStaticID},
			Admins:      []primitive.ObjectID{ownerStaticID},
			Owner:       &ownerStaticID,
		})

		textMessageModel := &message.TextMessage{Message: "Hello World!"}
		messageModel := message.NewMessage(ownerStaticID, message.TypeTextMessage, textMessageModel)

		savedMessageModel, err := messageRepo.Insert(chatModel.ChatID, messageModel)
		assert.NoError(t, err)

		textMessageContentModel, err := utils.TypeConverter[message.TextMessage](savedMessageModel.Content)
		assert.NoError(t, err)

		assert.Equal(t, textMessageModel.Message, textMessageContentModel.Message)
	})

	mt.Run("test delete message", func(mt *mtest.T) {
		messageRepo := NewRepository(mt.DB)

		ownerStaticID := primitive.NewObjectID()
		chatModel := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
			Title:       "MyChannel",
			Username:    "my_channel",
			Description: "Example description",
			Members:     []primitive.ObjectID{ownerStaticID},
			Admins:      []primitive.ObjectID{ownerStaticID},
			Owner:       &ownerStaticID,
		})

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		textMessageModel := &message.TextMessage{Message: "Hello World!"}
		messageModel := message.NewMessage(ownerStaticID, message.TypeTextMessage, textMessageModel)

		savedMessageModel, err := messageRepo.Insert(chatModel.ChatID, messageModel)
		assert.NoError(t, err)

		expectedResult := []bson.E{{Key: "NModified", Value: 1}, {Key: "N", Value: 1}}
		mt.AddMockResponses(mtest.CreateSuccessResponse(expectedResult...))

		err = messageRepo.Delete(chatModel.ChatID, savedMessageModel.MessageID)
		assert.NoError(t, err)
	})
	mt.Run("test update message", func(mt *mtest.T) {
		messageRepo := NewRepository(mt.DB)

		ownerStaticID := primitive.NewObjectID()
		chatModel := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
			Title:       "MyChannel",
			Username:    "my_channel",
			Description: "Example description",
			Members:     []primitive.ObjectID{ownerStaticID},
			Admins:      []primitive.ObjectID{ownerStaticID},
			Owner:       &ownerStaticID,
		})

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		textMessageModel := &message.TextMessage{Message: "Hello World!"}
		messageModel := message.NewMessage(ownerStaticID, message.TypeTextMessage, textMessageModel)

		savedMessageModel, err := messageRepo.Insert(chatModel.ChatID, messageModel)
		assert.NoError(t, err)

		fieldsToUpdate := bson.M{
			"message": "Message changed!",
		}
		expectedResult := []bson.E{{Key: "NModified", Value: 1}, {Key: "N", Value: 1}}
		mt.AddMockResponses(mtest.CreateSuccessResponse(expectedResult...))

		err = messageRepo.Update(chatModel.ChatID, savedMessageModel.MessageID, fieldsToUpdate)
		assert.NoError(t, err)
	})
}
