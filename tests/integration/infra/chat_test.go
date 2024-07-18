package infra

import (
	"log"
	"testing"

	stream_producer "github.com/kavkaco/Kavka-Core/infra/stream/producer"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestChat(t *testing.T) {
	chatProducer, err := stream_producer.NewChatStreamProducer(&kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}

	chatModel := model.Chat{
		ChatID:   primitive.NewObjectID(),
		ChatType: "channel",
	}

	err = chatProducer.ChatCreated(nil, chatModel)
	if err != nil {
		log.Fatal(err)
	}
}
