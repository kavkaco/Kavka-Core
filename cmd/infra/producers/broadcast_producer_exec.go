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

	const userID = "my-user1"

	// Create producer
	prod, err := stream_producers.NewBroadcastStreamProducer(&kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}

	chatModel := model.Chat{
		ChatID:   primitive.NewObjectID(),
		ChatType: "channel",
	}

	err = prod.ChatCreated(userID, chatModel)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 1)

	fmt.Println("Message sent!")
}
