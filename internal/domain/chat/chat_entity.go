package chat

import (
	"Kavka/internal/domain/message"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ChatTypeChannel = "channel"
	ChatTypeGroup   = "group"
	ChatTypeDirect  = "direct"
)

type Chat struct {
	ChatID     primitive.ObjectID `bson:"_id" json:"chat_id"`
	ChatType   string             `bson:"chat_type" json:"chat_type"`
	ChatDetail interface{}        `bson:"chat_detail" json:"chat_detail"`
	Messages   []message.Message  `json:"messages"`
}

// Chat Detail

type ChannelChatDetail struct {
	Title        string               `json:"title"`
	Members      []primitive.ObjectID `json:"members"`
	Admins       []primitive.ObjectID `json:"admins"`
	Owner        primitive.ObjectID   `json:"owner" bson:"owner"`
	RemovedUsers []primitive.ObjectID `bson:"removed_users" json:"removed_users"`
	Username     string               `json:"username"`
	Description  string               `json:"description"`
}

type GroupChatDetail struct {
	Title        string               `json:"title"`
	Members      []primitive.ObjectID `json:"members"`
	Admins       []primitive.ObjectID `json:"admins"`
	Owner        primitive.ObjectID   `json:"owner" bson:"owner"`
	RemovedUsers []primitive.ObjectID `bson:"removed_users" json:"removed_users"`
	Username     string               `json:"username"`
	Description  string               `json:"description"`
}

type DirectChatDetail struct {
	// ID of the users that chats with each other
	Sides [2]*primitive.ObjectID `json:"sides"`
}

func (d *DirectChatDetail) HasSide(staticID primitive.ObjectID) bool {
	has := false
	for _, v := range d.Sides {
		if *v == staticID {
			has = true
			break
		}
	}
	return has
}

func NewChat(chatType string, chatDetail interface{}) *Chat {
	m := &Chat{}

	m.ChatType = chatType
	m.ChatDetail = chatDetail
	m.ChatID = primitive.NewObjectID()
	m.Messages = []message.Message{}

	return m
}

func (c *Chat) GetMessage(messageID primitive.ObjectID) *message.Message {
	for _, v := range c.Messages {
		if v.MessageID == messageID {
			return &v
		}
	}

	return nil
}

func GetChatDetailBSON(chatDetail interface{}) ([]byte, error) {
	chatDetailD := chatDetail.(primitive.D)

	chatDetailBSON, marshalErr := bson.Marshal(chatDetailD)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return chatDetailBSON, nil
}
