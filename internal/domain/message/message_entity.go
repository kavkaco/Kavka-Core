package message

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

// Interfaces

type Repository interface {
	Insert(chatID primitive.ObjectID, msg *Message) (*Message, error)
	Update(chatID primitive.ObjectID, messageID primitive.ObjectID, fieldsToUpdate bson.M) error
	Delete(chatID primitive.ObjectID, messageID primitive.ObjectID) error
}

type Service interface {
	InsertTextMessage(chatID primitive.ObjectID, staticID primitive.ObjectID, messageContent string) (*Message, error)
	DeleteMessage(chatID primitive.ObjectID, messageID primitive.ObjectID) error
}
