package repository

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap"
)

func TestChatRepository(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // nolint

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test create channel", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		chatRepo := NewRepository(logger, mt.DB)

		ownerStaticID := primitive.NewObjectID()
		model := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
			Title:       "MyChannel",
			Username:    "my_channel",
			Description: "Example description",
			Members:     []primitive.ObjectID{ownerStaticID},
			Admins:      []primitive.ObjectID{ownerStaticID},
			Owner:       &ownerStaticID,
		})

		savedModel, err := chatRepo.Create(*model)
		assert.NoError(t, err)

		chatDetail, err := utils.TypeConverter[chat.ChannelChatDetail](savedModel.ChatDetail)
		assert.NoError(t, err)

		assert.Equal(t, chatDetail.Title, "MyChannel")
		assert.Equal(t, chatDetail.Username, "my_channel")
		assert.Equal(t, chatDetail.Owner.Hex(), ownerStaticID.Hex())
		assert.True(t, savedModel.IsMember(ownerStaticID))
		assert.True(t, savedModel.IsAdmin(ownerStaticID))
	})

	mt.Run("test create direct", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		chatRepo := NewRepository(logger, mt.DB)

		user1StaticID := primitive.NewObjectID()
		user2StaticID := primitive.NewObjectID()
		model := chat.NewChat(chat.TypeChannel, &chat.DirectChatDetail{
			Sides: [2]primitive.ObjectID{user1StaticID, user2StaticID},
		})

		savedModel, err := chatRepo.Create(*model)
		assert.NoError(t, err)

		chatDetail, err := utils.TypeConverter[chat.DirectChatDetail](savedModel.ChatDetail)
		assert.NoError(t, err)

		assert.True(t, chatDetail.HasSide(user1StaticID))
		assert.True(t, chatDetail.HasSide(user2StaticID))
	})

	mt.Run("test find by id", func(mt *mtest.T) {
		chatRepo := NewRepository(logger, mt.DB)

		chatID := primitive.NewObjectID()
		ownerStaticID := primitive.NewObjectID()
		expectedDoc := bson.D{
			{Key: "chat_id", Value: chatID},
			{Key: "chat_type", Value: "channel"},
			{
				Key: "chat_detail",
				Value: &chat.ChannelChatDetail{
					Title:       "MyChannel",
					Username:    "my_channel",
					Description: "Example description",
					Members:     []primitive.ObjectID{ownerStaticID},
					Admins:      []primitive.ObjectID{ownerStaticID},
					Owner:       &ownerStaticID,
				},
			},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "myDB.chats", mtest.FirstBatch, expectedDoc))

		model, err := chatRepo.FindByID(chatID)
		assert.NoError(t, err)

		assert.Equal(t, model.ChatType, expectedDoc.Map()["chat_type"])
	})

	mt.Run("test find chat or sides by static id", func(mt *mtest.T) {
		chatRepo := NewRepository(logger, mt.DB)

		chatID := primitive.NewObjectID()
		user1StaticID := primitive.NewObjectID()
		user2StaticID := primitive.NewObjectID()
		expectedDoc := bson.D{
			{Key: "chat_id", Value: chatID},
			{Key: "chat_type", Value: "direct"},
			{
				Key: "chat_detail",
				Value: &chat.DirectChatDetail{
					Sides: [2]primitive.ObjectID{user1StaticID, user2StaticID},
				},
			},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "myDB.chats", mtest.FirstBatch, expectedDoc))

		// Find by one of the sides
		model, err := chatRepo.FindChatOrSidesByStaticID(user1StaticID)
		assert.NoError(t, err)

		chatDetail, err := utils.TypeConverter[chat.DirectChatDetail](model.ChatDetail)
		assert.NoError(t, err)

		assert.True(t, chatDetail.HasSide(user1StaticID))
		assert.True(t, chatDetail.HasSide(user2StaticID))
		assert.Equal(t, model.ChatType, expectedDoc.Map()["chat_type"])
	})

	mt.Run("test find by sides", func(mt *mtest.T) {
		chatRepo := NewRepository(logger, mt.DB)

		chatID := primitive.NewObjectID()
		user1StaticID := primitive.NewObjectID()
		user2StaticID := primitive.NewObjectID()
		expectedDoc := bson.D{
			{Key: "chat_id", Value: chatID},
			{Key: "chat_type", Value: "direct"},
			{
				Key: "chat_detail",
				Value: &chat.DirectChatDetail{
					Sides: [2]primitive.ObjectID{user1StaticID, user2StaticID},
				},
			},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "myDB.chats", mtest.FirstBatch, expectedDoc))

		model, err := chatRepo.FindBySides([2]primitive.ObjectID{user1StaticID, user2StaticID})
		assert.NoError(t, err)

		chatDetail, err := utils.TypeConverter[chat.DirectChatDetail](model.ChatDetail)
		assert.NoError(t, err)

		assert.True(t, chatDetail.HasSide(user1StaticID))
		assert.True(t, chatDetail.HasSide(user2StaticID))
		assert.Equal(t, model.ChatType, expectedDoc.Map()["chat_type"])
	})
}
