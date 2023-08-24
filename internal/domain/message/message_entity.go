package message

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	MessageID primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time
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
	ImageUrl string
	Caption  string
}

func NewMessage(messageType string, content interface{}) *Message {
	m := &Message{}

	m.Content = content

	// set timestamps
	now := time.Now()
	m.CreatedAt = now

	return m
}
