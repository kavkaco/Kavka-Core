package model

import (
	"fmt"

	"github.com/kavkaco/Kavka-Core/utils/random"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserID = string
	User   struct {
		UserID        UserID         `bson:"user_id" json:"userID"`
		Name          string         `bson:"name" json:"name"`
		LastName      string         `bson:"last_name" json:"lastName"`
		Email         string         `bson:"email" json:"email"`
		Username      string         `bson:"username" json:"username"`
		ChatsListIDs  []ChatID       `bson:"chats_list_ids"`
		Biography     string         `bson:"biography" json:"biography"`
		ProfilePhotos []ProfilePhoto `bson:"profile_photos" json:"profilePhotos"`
	}
)

func (u *User) IncludesChatID(chatID ChatID) bool {
	for _, v := range u.ChatsListIDs {
		if v.Hex() == chatID.Hex() {
			return true
		}
	}

	return false
}

func NewUser(name, lastName, email, username string) *User {
	user := User{}

	user.UserID = fmt.Sprintf("%d", random.GenerateUserID())
	user.Name = name
	user.LastName = lastName
	user.Username = username
	user.Email = email
	user.ChatsListIDs = []primitive.ObjectID{}

	return &user
}

func (u *User) FullName() string {
	return u.Name + " " + u.LastName
}
