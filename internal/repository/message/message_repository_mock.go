package repository

import (
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/message"
	"github.com/kavkaco/Kavka-Core/utils/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRepository struct {
	chats []*chat.Chat
}

func (repo MockRepository) Insert(chatID primitive.ObjectID, msg *message.Message) (*message.Message, error) {
	for i, c := range repo.chats {
		if c.ChatID.Hex() == chatID.Hex() {
			repo.chats[i].Messages = append(c.Messages, msg)
			return msg, nil
		}
	}

	return nil, ErrChatNotFound
}

func (repo MockRepository) Update(chatID primitive.ObjectID, messageID primitive.ObjectID, fieldsToUpdate bson.M) error {
	for i, c := range repo.chats {
		if c.ChatID.Hex() == chatID.Hex() {
			chatMessages := c.Messages

			for j, m := range chatMessages {
				if m.MessageID.Hex() == messageID.Hex() {
					// Change value in the message
					for key, value := range fieldsToUpdate {
						msg := repo.chats[i].Messages[j]

						setErr := structs.SetFieldByBSON(msg, key, value)
						if setErr != nil {
							return setErr
						}
					}

					return nil
				}
			}
		}
	}

	return ErrChatNotFound
}

func (repo MockRepository) Delete(chatID primitive.ObjectID, messageID primitive.ObjectID) error {
	for i, c := range repo.chats {
		if c.ChatID.Hex() == chatID.Hex() { // Chat found!
			for j, m := range c.Messages {
				if m.MessageID.Hex() == messageID.Hex() {
					messagesList := repo.chats[i].Messages
					repo.chats[i].Messages = append(messagesList[:j], messagesList[j+1:]...)
					return nil
				}
			}

			return ErrMsgNotFound
		}
	}

	return ErrChatNotFound
}

// NewMockRepository takes the @existingChats to mock repo because
// we do not want to focus on creating a chat here in this mock-repository.
func NewMockRepository(existingChats []*chat.Chat) message.Repository {
	return &MockRepository{existingChats}
}
