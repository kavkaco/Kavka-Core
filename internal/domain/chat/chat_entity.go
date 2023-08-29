package chat

import (
	"Kavka/internal/domain/message"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ChatTypeChannel = "channel"
	ChatTypeGroup   = "group"
	ChatTypeDirect  = "direct"
)

type Chat struct {
	ChatID     primitive.ObjectID `bson:"_id"`
	ChatType   string             `bson:"chat_type"`
	ChatDetail interface{}        `bson:"chat_detail"`
	Messages   []*message.Message
}

// Chat Detail

type ChannelChatDetail struct {
	Members      []*primitive.ObjectID
	Admins       []*primitive.ObjectID
	RemovedUsers []*primitive.ObjectID `bson:"removed_users"`
	Username     string
	Description  string
}

type GroupChatDetail struct {
	Members      []*primitive.ObjectID
	Admins       []*primitive.ObjectID
	RemovedUsers []*primitive.ObjectID `bson:"removed_users"`
	Username     string
	Description  string
}

type DirectChatDetail struct {
	// ID of the users that chats with each other
	Sides [2]*primitive.ObjectID
}

func (c *Chat) GetMessageByID(id primitive.ObjectID) *message.Message {
	for _, v := range c.Messages {
		if v.MessageID == id {
			return v
		}
	}

	return nil
}

func NewChat(chatType string, chatDetail interface{}) *Chat {
	m := &Chat{}

	m.ChatType = chatType
	m.ChatDetail = chatDetail
	m.ChatID = primitive.NewObjectID()
	m.Messages = []*message.Message{}

	return m
}
