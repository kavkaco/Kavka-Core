package presenters

import (
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/utils"
)

func MessageAsJSON(obj model.Message) (*model.Message, error) {
	var messageContent interface{}

	switch obj.Type {
	case model.TypeTextMessage:
		localMessageContent, err := utils.TypeConverter[model.TextMessage](obj.Content)
		if err != nil {
			return nil, err
		}

		messageContent = localMessageContent
	}

	obj.Content = messageContent

	return &obj, nil
}
