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
	MessageID primitive.ObjectID `bson:"message_id" json:"message_id"`
	SenderID  primitive.ObjectID `bson:"sender_id"  json:"sender_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Edited    bool
	Seen      bool

	// MessageType
	Type    string
	Content interface{}
}

type TextMessage struct {
	Message string
}

type ImageMessage struct {
	ImageURL string
	Caption  string
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
