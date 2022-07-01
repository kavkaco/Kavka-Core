package models

type Message struct {
	MessageID   int    `bson:"message_id"`
	SenderID    int    `bson:"sender_id"`
	MessageType int    `bson:"message_type"`
	SendTime    int    `bson:"send_time"`
	Edited      bool   `bson:"edited"`
	SeenState   string `bson:"seen_state"`
	// TextMessage
	TextContent string `bson:"text_content"`
	// ImageMessage
	Image string `bson:"image"`
}
