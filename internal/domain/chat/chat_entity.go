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
	ChatType   string
	ChatDetail interface{}
	Messages   []*message.Message
}

// Chat Detail

type ChannelChatDetail struct {
	Members []*primitive.ObjectID
	Admins  []*primitive.ObjectID
}

type GroupChatDetail struct {
	Members []*primitive.ObjectID
	Admins  []*primitive.ObjectID
}

type DirectChatDetail struct {
	// ID of the users that chats with each other
	Sides [2]*primitive.ObjectID
}
