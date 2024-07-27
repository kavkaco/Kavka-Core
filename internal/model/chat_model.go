package model

import (
	"github.com/kavkaco/Kavka-Core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatID = primitive.ObjectID

const (
	TypeChannel = "channel"
	TypeGroup   = "group"
	TypeDirect  = "direct"
)

type Chat struct {
	ChatID      ChatID       `bson:"_id" json:"chatId"`
	ChatType    string       `bson:"chat_type" json:"chatType"`
	ChatDetail  interface{}  `bson:"chat_detail" json:"chatDetail"`
	LastMessage *LastMessage `bson:"last_message" json:"lastMessage"`
}

type Member struct {
	UserID   UserID `bson:"user_id" json:"userID"`
	Name     string `bson:"name" json:"name"`
	LastName string `bson:"last_name" json:"lastName"`

	// We will add the profile photo here later =)
}

// Chat Detail

type ChannelChatDetail struct {
	Title        string   `bson:"title" json:"title"`
	Members      []UserID `bson:"members,omitempty" json:"members"`
	Admins       []UserID `bson:"admins,omitempty" json:"admins"`
	Owner        UserID   `bson:"owner,omitempty"         json:"owner"`
	RemovedUsers []UserID `bson:"removed_users,omitempty" json:"removedUsers"`
	Username     string   `bson:"username,omitempty" json:"username"`
	Description  string   `bson:"description,omitempty" json:"description"`
}

type GroupChatDetail struct {
	Title        string   `bson:"title" json:"title"`
	Members      []UserID `bson:"members,omitempty" json:"members"`
	Admins       []UserID `bson:"admins,omitempty" json:"admins"`
	Owner        UserID   `bson:"owner,omitempty"         json:"owner"`
	RemovedUsers []UserID `bson:"removed_users,omitempty" json:"removedUsers"`
	Username     string   `bson:"username,omitempty" json:"username"`
	Description  string   `bson:"description,omitempty" json:"description"`
}

type DirectChatDetail struct {
	// Chat partners
	Sides [2]UserID `bson:"sides" json:"sides"`
}

type DirectChatFetchedDetail struct {
	UserInfo User
}

func (c *Chat) IsMember(userID UserID) bool {
	d, _ := utils.TypeConverter[ChannelChatDetail](c.ChatDetail)
	for _, memberUserID := range d.Members {
		if memberUserID == userID {
			return true
		}
	}

	return false
}

func (c *Chat) IsAdmin(userID UserID) bool {
	d, _ := utils.TypeConverter[ChannelChatDetail](c.ChatDetail)
	for _, adminUserID := range d.Admins {
		if adminUserID == userID {
			return true
		}
	}

	return false
}

func (d *DirectChatDetail) HasSide(userID UserID) bool {
	has := false
	for _, v := range d.Sides {
		if v == userID {
			has = true
			break
		}
	}
	return has
}

func DetectRecipient(userIDs [2]UserID, currentUserID UserID) *UserID {
	if userIDs[0] == currentUserID {
		return &userIDs[1]
	}

	return &userIDs[0]
}

func NewChat(chatType string, chatDetail interface{}) *Chat {
	m := &Chat{}

	m.ChatType = chatType
	m.ChatDetail = chatDetail
	m.ChatID = primitive.NewObjectID()
	m.LastMessage = nil

	return m
}
