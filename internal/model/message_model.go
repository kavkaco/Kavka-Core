package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageID = primitive.ObjectID

const (
	TypeTextMessage  = "text"
	TypeImageMessage = "image"
)

type LastMessage struct {
	MessageType    string `bson:"type" json:"type"`
	MessageCaption string `bson:"caption" json:"caption"`
}

type MessageStore struct {
	ChatID   ChatID    `bson:"chat_id"`
	Messages []Message `bson:"messages"`
}

type Message struct {
	MessageID MessageID   `bson:"message_id" json:"messageId"`
	SenderID  UserID      `bson:"sender_id"  json:"senderId"`
	CreatedAt time.Time   `bson:"created_at" json:"createdAt"`
	Edited    bool        `bson:"edited" json:"edited"`
	Seen      bool        `bson:"seen" json:"seen"`
	Type      string      `bson:"type" json:"type"`
	Content   interface{} `bson:"content" json:"content"`
}

type TextMessage struct {
	Data string `bson:"data" json:"data"`
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

func NewLastMessage(messageType string, messageCaption string) *LastMessage {
	return &LastMessage{
		MessageType:    messageType,
		MessageCaption: messageCaption,
	}
}
