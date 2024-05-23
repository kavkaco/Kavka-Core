package presenters

import (
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/internal/model/user"
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

func UnmarshalFetchedUsers(fetchedUsers primitive.A) ([]user.User, error) {
	users := []user.User{}

	for _, v := range fetchedUsers {
		currentUser, err := utils.TypeConverter[user.User](v)
		if err != nil {
			return nil, err
		}

		users = append(users, *currentUser)
	}

	return users, nil
}

func ChatAsJSON(obj chat.ChatC, userStaticID primitive.ObjectID) (interface{}, error) {
	// Determine the specific ChatDetail type based on chatType
	var httpChatDetail interface{}

	switch obj.ChatType {
	case chat.TypeDirect:
		chatDetailLocal, err := utils.TypeConverter[map[string]interface{}](obj.ChatDetail)
		if err != nil {
			return nil, err
		}

		fetchedUsers, err := UnmarshalFetchedUsers((*chatDetailLocal)["fetchedUsers"].(primitive.A))
		if err != nil {
			return nil, err
		}

		if len(fetchedUsers) != 2 {
			return nil, errors.New("invalid length of fetched users")
		}

		// FIXME
		// target := chat.DetectTarget(fetchedUsers, userStaticID)

		// httpChatDetail = HttpDirectChatDetail{
		// 	UserID:   target.StaticID,
		// 	Name:     target.Name,
		// 	LastName: target.LastName,
		// }
	case chat.TypeChannel:
		chatDetailLocal, err := utils.TypeConverter[chat.ChannelChatDetail](obj.ChatDetail)
		if err != nil {
			return nil, err
		}

		httpChatDetail = HttpChannelChatDetail{
			Title:    chatDetailLocal.Title,
			Username: chatDetailLocal.Username,
		}
	case chat.TypeGroup:
		chatDetailLocal, err := utils.TypeConverter[chat.GroupChatDetail](obj.ChatDetail)
		if err != nil {
			return nil, err
		}

		httpChatDetail = HttpGroupChatDetail{
			Title:    chatDetailLocal.Title,
			Username: chatDetailLocal.Username,
		}
	}

	obj.ChatDetail = httpChatDetail

	for i, msg := range obj.Messages {
		messageJson, err := MessageAsJSON(*msg)
		if err != nil {
			return nil, err
		}

		obj.Messages[i] = messageJson
	}

	return obj, nil
}
