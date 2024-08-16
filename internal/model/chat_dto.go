package model

type ChatDTO struct {
	ChatID      ChatID      `bson:"_id" json:"chatId"`
	ChatType    string      `bson:"chat_type" json:"chatType"`
	ChatDetail  interface{} `bson:"chat_detail" json:"chatDetail"`
	LastMessage *Message    `bson:"last_message" json:"lastMessage"`
}

type Member struct {
	UserID   UserID `bson:"user_id" json:"userID"`
	Name     string `bson:"name" json:"name"`
	LastName string `bson:"last_name" json:"lastName"`

	// We will add the profile photo here later =)
}

func NewChatDTO(chatModel *Chat) *ChatDTO {
	m := &ChatDTO{}

	m.ChatID = chatModel.ChatID
	m.ChatType = chatModel.ChatType
	m.ChatDetail = chatModel.ChatDetail

	return m
}
