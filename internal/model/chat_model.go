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
	Sides [2]UserID `bson:"sides" json:"sides"`
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
