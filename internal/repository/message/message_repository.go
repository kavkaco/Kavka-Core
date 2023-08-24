package message

import (
	"Kavka/database"

	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct {
	messagesCollection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{
		db.Collection(database.MessagesCollection),
	}
}
