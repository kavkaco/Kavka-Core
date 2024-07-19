package stream_consumers

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/infra/stream"
)

var wg sync.WaitGroup

type ChatStreamConsumerEvents interface {
	ChatCreated(data map[string]interface{})
	ErrorOccurred(err error)
}

type ConsumerGroupHandler struct {
	encoder    stream.MessageEncoder
	eventsChan chan map[string]interface{}
}

func (cg *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	if config.CurrentEnv == config.Development {
		fmt.Println("Consumer joined group! =)")
	}

	return nil
}

func (cg *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	if config.CurrentEnv == config.Development {
		fmt.Println("Consumer left group.")
	}

	return nil
}

func (cg *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	go func() {
		defer wg.Done()

		for msg := range claim.Messages() {
			decoded, err := cg.encoder.Decode(string(msg.Value))
			if err != nil {
				// FIXME - Add logger later
				fmt.Println(err)
				continue
			}

			decoded["eventName"] = string(msg.Key)

			cg.eventsChan <- decoded

			sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
		}
	}()

	return nil
}

func NewChatStreamConsumer(ctx context.Context, config config.Kafka, encoder stream.MessageEncoder) (chan map[string]interface{}, error) {
	eventsChan := make(chan map[string]interface{})

	c, err := sarama.NewConsumerGroup(config.Brokers, "1", config.Sarama)
	if err != nil {
		return nil, err
	}

	handler := &ConsumerGroupHandler{encoder, eventsChan}

	err = c.Consume(context.TODO(), []string{stream.KafkaTopics().ChatTopic}, handler)
	if err != nil {
		return nil, err
	}

	wg.Wait()

	return eventsChan, nil
}
