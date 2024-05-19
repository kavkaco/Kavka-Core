package repository

// func TestMessageRepository(t *testing.T) {
// 	logger, _ := zap.NewDevelopment()
// 	defer logger.Sync() // nolint

// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()
// 	mt.Run("test insert text message", func(mt *mtest.T) {
// 		messageRepo := NewRepository(logger, mt.DB)

// 		expectedResult := []bson.E{
// 			{Key: "NModified", Value: 1},
// 			{Key: "N", Value: 1},
// 		}
// 		mt.AddMockResponses(mtest.CreateSuccessResponse(expectedResult...))

// 		ownerStaticID := primitive.NewObjectID()
// 		chatModel := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
// 			Title:       "MyChannel",
// 			Username:    "my_channel",
// 			Description: "Example description",
// 			Members:     []primitive.ObjectID{ownerStaticID},
// 			Admins:      []primitive.ObjectID{ownerStaticID},
// 			Owner:       &ownerStaticID,
// 		})

// 		textMessageModel := &message.TextMessage{Data: "Hello World!"}
// 		messageModel := message.NewMessage(ownerStaticID, message.TypeTextMessage, textMessageModel)

// 		savedMessageModel, err := messageRepo.Insert(chatModel.ChatID, messageModel)
// 		assert.NoError(t, err)

// 		textMessageContentModel, err := utils.TypeConverter[message.TextMessage](savedMessageModel.Content)
// 		assert.NoError(t, err)

// 		assert.Equal(t, textMessageModel.Data, textMessageContentModel.Data)
// 	})

// 	mt.Run("test delete message", func(mt *mtest.T) {
// 		messageRepo := NewRepository(logger, mt.DB)

// 		ownerStaticID := primitive.NewObjectID()
// 		chatModel := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
// 			Title:       "MyChannel",
// 			Username:    "my_channel",
// 			Description: "Example description",
// 			Members:     []primitive.ObjectID{ownerStaticID},
// 			Admins:      []primitive.ObjectID{ownerStaticID},
// 			Owner:       &ownerStaticID,
// 		})

// 		mt.AddMockResponses(mtest.CreateSuccessResponse())

// 		textMessageModel := &message.TextMessage{Data: "Hello World!"}
// 		messageModel := message.NewMessage(ownerStaticID, message.TypeTextMessage, textMessageModel)

// 		savedMessageModel, err := messageRepo.Insert(chatModel.ChatID, messageModel)
// 		assert.NoError(t, err)

// 		expectedResult := []bson.E{{Key: "NModified", Value: 1}, {Key: "N", Value: 1}}
// 		mt.AddMockResponses(mtest.CreateSuccessResponse(expectedResult...))

// 		err = messageRepo.Delete(chatModel.ChatID, savedMessageModel.MessageID)
// 		assert.NoError(t, err)
// 	})
// 	mt.Run("test update message", func(mt *mtest.T) {
// 		messageRepo := NewRepository(logger, mt.DB)

// 		ownerStaticID := primitive.NewObjectID()
// 		chatModel := chat.NewChat(chat.TypeChannel, &chat.ChannelChatDetail{
// 			Title:       "MyChannel",
// 			Username:    "my_channel",
// 			Description: "Example description",
// 			Members:     []primitive.ObjectID{ownerStaticID},
// 			Admins:      []primitive.ObjectID{ownerStaticID},
// 			Owner:       &ownerStaticID,
// 		})

// 		mt.AddMockResponses(mtest.CreateSuccessResponse())

// 		textMessageModel := &message.TextMessage{Data: "Hello World!"}
// 		messageModel := message.NewMessage(ownerStaticID, message.TypeTextMessage, textMessageModel)

// 		savedMessageModel, err := messageRepo.Insert(chatModel.ChatID, messageModel)
// 		assert.NoError(t, err)

// 		fieldsToUpdate := bson.M{
// 			"message": "Message changed!",
// 		}
// 		expectedResult := []bson.E{{Key: "NModified", Value: 1}, {Key: "N", Value: 1}}
// 		mt.AddMockResponses(mtest.CreateSuccessResponse(expectedResult...))

// 		err = messageRepo.Update(chatModel.ChatID, savedMessageModel.MessageID, fieldsToUpdate)
// 		assert.NoError(t, err)
// 	})
// }
