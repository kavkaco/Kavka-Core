package repository

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type MockRepository struct {
	chats []chat.Chat
}

func NewMockRepository() chat.Repository {
	return &MockRepository{}
}

func (repo *MockRepository) Create(newChat chat.Chat) (*chat.Chat, error) {
	repo.chats = append(repo.chats, newChat)

	return &newChat, nil
}

func (repo *MockRepository) Where(filter bson.M) ([]chat.Chat, error) {
	var filterKey string
	var filterValue interface{}

	for k, v := range filter {
		filterKey = k
		filterValue = v
	}

	var result []chat.Chat

	if len(repo.chats) == 0 {
		return result, nil
	}

	for _, row := range repo.chats {
		// Check filter
		fields := structs.Fields(row)

		for _, field := range fields {
			tag := field.Tag("bson")
			fieldValue := field.Value()
			fieldValueType := reflect.TypeOf(fieldValue).Name()

			if filterKey == tag {
				if (fieldValue == filterValue) || (fieldValueType == "ObjectID" && fieldValue.(primitive.ObjectID).Hex() == filterValue.(primitive.ObjectID).Hex()) {
					result = append(result, row)
				}
			}
		}
	}

	return result, nil
}

func (repo *MockRepository) Destroy(chatID primitive.ObjectID) error {
	for index, row := range repo.chats {
		if row.ChatID == chatID {
			repo.chats = append(repo.chats[:index], repo.chats[index+1:]...)
			break
		}
	}

	return ErrChatNotFound
}

func (repo *MockRepository) findBy(filter bson.M) (*chat.Chat, error) {
	result, err := repo.Where(filter)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		row := result[len(result)-1]

		return &row, nil
	}

	return nil, ErrChatNotFound
}

func (repo *MockRepository) FindByID(staticID primitive.ObjectID) (*chat.Chat, error) {
	filter := bson.M{"id": staticID}
	return repo.findBy(filter)
}

func (repo *MockRepository) FindChatOrSidesByStaticID(staticID primitive.ObjectID) (*chat.Chat, error) {
	// Check as a group or channel with StaticID
	resultByID, err := repo.FindByID(staticID)
	if errors.Is(err, ErrChatNotFound) {
		// This condition means it is not a channel or group
		// So we must look up in the Sides of direct-chats.
		for _, c := range repo.chats {
			if c.ChatType == chat.TypeDirect {
				chatDetail, err := utils.TypeConverter[chat.DirectChatDetail](c.ChatDetail)
				if err != nil {
					return nil, err
				}

				if chatDetail.HasSide(staticID) {
					return &c, nil
				}
			}
		}
	} else {
		return resultByID, nil
	}

	return nil, ErrChatNotFound
}
