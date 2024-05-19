package user

import (
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/utils/random"
)

type UserID = string
type User struct {
	UserID       UserID       `bson:"user_id" json:"user_id"`
	Name         string       `bson:"name" json:"name"`
	LastName     string       `bson:"last_name" json:"lastName"`
	Email        string       `bson:"email" json:"email"`
	Username     string       `bson:"username" json:"username"`
	Profile      Profile      `bson:"profile" json:"profile"`
	OnlineStatus interface{}  `bson:"online_status" json:"onlineStatus"`
	ChatsList    []*chat.Chat `bson:"chats_list" json:"chatsList"`
}

func NewUser(email, username string) *User {
	user := User{}
	user.UserID = string(rune(random.GenerateUserID()))
	user.Username = random.GenerateUsername()
	user.Email = email

	return &user
}

func (u *User) FullName() string {
	return u.Name + " " + u.LastName
}
