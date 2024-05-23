package presenters

import (
	"github.com/kavkaco/Kavka-Core/internal/model/message"
	"github.com/kavkaco/Kavka-Core/utils"
)

func MessageAsJSON(obj message.Message) (*message.Message, error) {
	var messageContent interface{}

	switch obj.Type {
	case message.TypeTextMessage:
		localMessageContent, err := utils.TypeConverter[message.TextMessage](obj.Content)
		if err != nil {
			return nil, err
		}

		messageContent = localMessageContent
	}

	obj.Content = messageContent

	return &obj, nil
}
