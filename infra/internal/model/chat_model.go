package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatID = primitive.ObjectID

const (
	TypeChannel = "channel"
	TypeGroup   = "group"
	TypeDirect  = "direct"
)

type Chat struct {
	ChatID     ChatID      `bson:"_id" json:"chatId"`
	ChatType   string      `bson:"chat_type" json:"chatType"`
	ChatDetail interface{} `bson:"chat_detail" json:"chatDetail"`
}

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
	UserID          UserID `bson:"user_id" json:"userId"`
	RecipientUserID UserID `bson:"recipient_user_id" json:"recipientUserId"`
}

func (c *ChannelChatDetail) IsMember(userID UserID) bool {
	for _, memberUserID := range c.Members {
		if memberUserID == userID {
			return true
		}
	}

	return false
}

func (c *ChannelChatDetail) IsAdmin(userID UserID) bool {
	for _, adminUserID := range c.Admins {
		if adminUserID == userID {
			return true
		}
	}

	return false
}

func (c *GroupChatDetail) IsMember(userID UserID) bool {
	for _, memberUserID := range c.Members {
		if memberUserID == userID {
			return true
		}
	}

	return false
}

func (c *GroupChatDetail) IsAdmin(userID UserID) bool {
	for _, adminUserID := range c.Admins {
		if adminUserID == userID {
			return true
		}
	}

	return false
}

func (d *DirectChatDetail) HasSide(userID UserID) bool {
	if d.RecipientUserID == userID {
		return true
	} else if d.UserID == userID {
		return true
	}

	return false
}

func (d *DirectChatDetail) GetRecipient(userID UserID) UserID {
	if d.UserID == userID {
		return d.RecipientUserID
	} else if d.RecipientUserID == userID {
		return d.UserID
	}

	return ""
}

// Safe means no duplication
func (d *ChannelChatDetail) AddMemberSafely(userID UserID) {
	isMember := d.IsMember(userID)
	if !isMember {
		d.Members = append(d.Members, userID)
	}
}

// Safe means no duplication
func (d *GroupChatDetail) AddMemberSafely(userID UserID) {
	isMember := d.IsMember(userID)
	if !isMember {
		d.Members = append(d.Members, userID)
	}
}

// Safe means no duplication
func (d *ChannelChatDetail) AddAdminSafely(userID UserID) {
	isAdmin := d.IsAdmin(userID)
	if !isAdmin {
		d.Admins = append(d.Admins, userID)
	}
}

// Safe means no duplication
func (d *GroupChatDetail) AddAdminSafely(userID UserID) {
	isAdmin := d.IsAdmin(userID)
	if !isAdmin {
		d.Admins = append(d.Admins, userID)
	}
}

func NewChat(chatType string, chatDetail interface{}) *Chat {
	m := &Chat{}

	m.ChatType = chatType
	m.ChatDetail = chatDetail
	m.ChatID = primitive.NewObjectID()

	return m
}

func ParseChatID(chatID string) (ChatID, error) {
	return primitive.ObjectIDFromHex(chatID)
}

func NewChatID() ChatID {
	return primitive.NewObjectID()
}
