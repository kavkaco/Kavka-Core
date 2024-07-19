package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kavkaco/Kavka-Core/config"
	stream_producers "github.com/kavkaco/Kavka-Core/infra/stream/producer"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	kafkaConfig := config.Read().Kafka

	// Create producer
	chatProducer, err := stream_producers.NewChatStreamProducer(&kafkaConfig)
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

	time.Sleep(time.Second * 1)

	fmt.Println("Message sent!")
}
