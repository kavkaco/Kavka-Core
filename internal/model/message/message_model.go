package message

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TypeTextMessage  = "text"
	TypeImageMessage = "image"
)

type Message struct {
	MessageID primitive.ObjectID `bson:"message_id" json:"messageId"`
	SenderID  primitive.ObjectID `bson:"sender_id"  json:"senderId"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	Edited    bool               `bson:"edited" json:"edited"`
	Seen      bool               `bson:"seen" json:"seen"`

	// MessageType
	Type    string      `bson:"type" json:"type"`
	Content interface{} `bson:"content" json:"content"`
}

type TextMessage struct {
	Data string `bson:"data" json:"data"`
}

type ImageMessage struct {
	ImageURL string `bson:"image_url" json:"imageUrl"`
	Caption  string `bson:"caption" json:"caption"`
}

func NewMessage(senderID primitive.ObjectID, messageType string, content interface{}) *Message {
	m := &Message{}

	m.Type = messageType
	m.Content = content
	m.MessageID = primitive.NewObjectID()
	m.SenderID = senderID

	// set timestamps
	now := time.Now()
	m.CreatedAt = now

	return m
}
