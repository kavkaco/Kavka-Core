package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageID = primitive.ObjectID

const (
	TypeTextMessage  = "text"
	TypeImageMessage = "image"
	TypeLabelMessage = "label"
)

type Message struct {
	MessageID MessageID   `bson:"message_id" json:"messageId"`
	SenderID  UserID      `bson:"sender_id"  json:"senderId"`
	CreatedAt time.Time   `bson:"created_at" json:"createdAt"`
	Edited    bool        `bson:"edited" json:"edited"`
	Seen      bool        `bson:"seen" json:"seen"`
	Type      string      `bson:"type" json:"type"`
	Content   interface{} `bson:"content" json:"content"`
}

type MessageGetter struct {
	Sender  *MessageSenderDTO `bson:"sender" json:"sender"`
	Message *Message          `bson:"message" json:"message"`
}

type TextMessage struct {
	Text string `bson:"text" json:"text"`
}

type LabelMessage struct {
	Text string `bson:"text" json:"text"`
}

type ImageMessage struct {
	ImageURL string `bson:"image_url" json:"imageUrl"`
	Caption  string `bson:"caption" json:"caption"`
}

func NewMessage(messageType string, content interface{}, senderID UserID) *Message {
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
