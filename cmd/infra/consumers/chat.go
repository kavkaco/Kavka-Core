package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	stream_consumers "github.com/kavkaco/Kavka-Core/infra/stream/consumers"
)

func main() {
	kafkaConfig := config.Read().Kafka

	messageEncoder := stream.NewMessageJsonEncoder()

	var wg sync.WaitGroup
	eventCh := make(chan map[string]interface{})

	go func() {
		defer wg.Done()

		for {
			event := <-eventCh

			fmt.Println("Event", event)
		}
	}()
	wg.Add(1)

	err := stream_consumers.NewChatStreamConsumer(context.TODO(), kafkaConfig, messageEncoder, eventCh)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
