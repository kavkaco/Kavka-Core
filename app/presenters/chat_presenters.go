package presenters

import (
	"log"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HttpChannelChatDetail struct {
	Title    string `json:"title"`
	Username string `json:"username"`
}
type HttpGroupChatDetail struct {
	Title    string `json:"title"`
	Username string `json:"username"`
}

type HttpDirectChatDetail struct {
	UserID   primitive.ObjectID `json:"userId"`
	Name     string             `json:"name"`
	LastName string             `json:"lastName"`
}

func ChatAsJSON(chat model.Chat, userID model.UserID) (interface{}, error) {
	// Determine the specific ChatDetail type based on chatType
	var httpChatDetail interface{}

	switch chat.ChatType {
	case model.TypeDirect:
		log.Fatal("Not implemented yet")
		// chatDetailLocal, err := utils.TypeConverter[map[string]interface{}](chat.ChatDetail)
		// if err != nil {
		// 	return nil, err
		// }

		// fetchedUsers, err := UnmarshalFetchedUsers((*chatDetailLocal)["fetchedUsers"].(primitive.A))
		// if err != nil {
		// 	return nil, err
		// }

		// if len(fetchedUsers) != 2 {
		// 	return nil, errors.New("invalid length of fetched users")
		// }

		// recipientUserID := model.DetectRecipient(fetchedUsers, userID)

		// httpChatDetail = HttpDirectChatDetail{
		// 	UserID:   recipientUserID,
		// 	Name:     fetchedUsers,
		// 	LastName: target.LastName,
		// }
	case model.TypeChannel:
		chatDetailLocal, err := utils.TypeConverter[model.ChannelChatDetail](chat.ChatDetail)
		if err != nil {
			return nil, err
		}

		httpChatDetail = HttpChannelChatDetail{
			Title:    chatDetailLocal.Title,
			Username: chatDetailLocal.Username,
		}
	case model.TypeGroup:
		chatDetailLocal, err := utils.TypeConverter[model.GroupChatDetail](chat.ChatDetail)
		if err != nil {
			return nil, err
		}

		httpChatDetail = HttpGroupChatDetail{
			Title:    chatDetailLocal.Title,
			Username: chatDetailLocal.Username,
		}
	}

	chat.ChatDetail = httpChatDetail

	// for i, msg := range chat.Messages {
	// 	messageJson, err := MessageAsJSON(*msg)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	chat.Messages[i] = messageJson
	// }

	return chat, nil
}
