package stream

type KafkaTopicsStruct struct {
	ChatTopic string
}

func KafkaTopics() *KafkaTopicsStruct {
	return &KafkaTopicsStruct{
		ChatTopic: "chats",
	}
}
