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

	const userID = "my-user1"

	var wg sync.WaitGroup
	eventsCh := make(chan stream.Event)

	go func() {
		defer wg.Done()

		for {
			event := <-eventsCh

			fmt.Println("EventName", event.Name)
			fmt.Println("Data", event.Data)
			fmt.Println("")
		}
	}()
	wg.Add(1)

	cons, err := stream_consumers.NewBroadcastConsumer(context.TODO(), kafkaConfig.Brokers, *kafkaConfig.Sarama)
	if err != nil {
		log.Fatal(err)
	}

	err = cons.SubscribeForUser(userID, eventsCh)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
