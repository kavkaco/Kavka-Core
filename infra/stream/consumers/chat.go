package stream_consumers

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/infra/stream"
)

type ChatStreamConsumerEvents interface {
	ChatCreated(data map[string]interface{})
	ErrorOccurred(err error)
}

type ConsumerGroupHandler struct {
	encoder stream.MessageEncoder
	eventCh chan map[string]interface{}
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
	var wg sync.WaitGroup

	go func() {
		defer wg.Done()

		for msg := range claim.Messages() {
			decoded, err := cg.encoder.Decode(string(msg.Value))
			if err != nil {
				continue
			}

			decoded["eventName"] = string(msg.Key)

			go func() {
				cg.eventCh <- decoded
			}()

			sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
		}
	}()
	wg.Add(1)
	wg.Wait()

	return nil
}

func NewChatStreamConsumer(ctx context.Context, config config.Kafka, encoder stream.MessageEncoder, eventCh chan map[string]interface{}) error {
	consumer, err := sarama.NewConsumerGroup(config.Brokers, "consumer-group-default", config.Sarama)
	if err != nil {
		return err
	}

	handler := &ConsumerGroupHandler{encoder, eventCh}

	err = consumer.Consume(context.TODO(), []string{stream.KafkaTopics().ChatTopic}, handler)
	if err != nil {
		return err
	}

	return nil
}
