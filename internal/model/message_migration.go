package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageDocumentV2 struct {
	ID          MessageID  `bson:"_id" json:"id"`
	ChatID      ChatID     `bson:"chat_id" json:"chatId"`
	SenderID    UserID     `bson:"sender_id" json:"senderId"`
	SenderName  string     `bson:"sender_name" json:"senderName"`
	SenderLastName string  `bson:"sender_last_name" json:"senderLastName"`
	SenderUsername string  `bson:"sender_username" json:"senderUsername"`
	CreatedAt   time.Time  `bson:"created_at" json:"createdAt"`
	Edited      bool       `bson:"edited" json:"edited"`
	Seen        bool       `bson:"seen" json:"seen"`
	Type        string     `bson:"type" json:"type"`
	Content     interface{} `bson:"content" json:"content"`
}

func NewMessageDocumentV2(chatID ChatID, sender *MessageSenderDTO, msg *Message) *MessageDocumentV2 {
	return &MessageDocumentV2{
		ID:             msg.MessageID,
		ChatID:         chatID,
		SenderID:       sender.UserID,
		SenderName:     sender.Name,
		SenderLastName: sender.LastName,
		SenderUsername: sender.Username,
		CreatedAt:      msg.CreatedAt,
		Edited:         msg.Edited,
		Seen:           msg.Seen,
		Type:           msg.Type,
		Content:        msg.Content,
	}
}

func (m *MessageDocumentV2) ToMessageGetter() *MessageGetter {
	return &MessageGetter{
		Sender: &MessageSenderDTO{
			UserID:   m.SenderID,
			Name:     m.SenderName,
			LastName: m.SenderLastName,
			Username: m.SenderUsername,
		},
		Message: &Message{
			MessageID: m.ID,
			SenderID:  m.SenderID,
			CreatedAt: m.CreatedAt,
			Edited:    m.Edited,
			Seen:      m.Seen,
			Type:      m.Type,
			Content:   m.Content,
		},
	}
}

type MigrationState struct {
	CollectionName string
	Migrated       bool
	TotalDocs      int64
	LastMigratedID primitive.ObjectID
}
