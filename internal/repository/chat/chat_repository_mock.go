package repository

import (
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockChatRepository struct {
	chats []*chat.Chat
}

func NewMockChatRepository() chat.ChatRepository {
	return &MockChatRepository{}
}

func (repo *MockChatRepository) Create(newChat *chat.Chat) (*chat.Chat, error) {
	repo.chats = append(repo.chats, newChat)

	return newChat, nil
}

func (repo *MockChatRepository) Where(filter any) ([]*chat.Chat, error) {
	return nil, nil
}

func (repo *MockChatRepository) Destroy(chatID primitive.ObjectID) error {
	return nil
}

func (repo *MockChatRepository) FindByID(staticID primitive.ObjectID) (*chat.Chat, error) {
	return nil, nil
}

func (repo *MockChatRepository) FindChatOrSidesByStaticID(staticID *primitive.ObjectID) (*chat.Chat, error) {
	return nil, nil
}

func (repo *MockChatRepository) FindBySides(sides [2]*primitive.ObjectID) (*chat.Chat, error) {
	return nil, nil
}
