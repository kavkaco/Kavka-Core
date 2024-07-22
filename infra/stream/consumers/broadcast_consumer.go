package stream_consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sync"

	"github.com/IBM/sarama"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/model"
)

// type BroadcastConsumerEvents interface {
// 	ChatCreated(data map[string]interface{})
// 	ErrorOccurred(err error)
// }

type consumerGroupHandler struct {
	encoder stream.MessageEncoder
	localCh chan eventReceivers
}

func (cg *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	if config.CurrentEnv == config.Development {
		fmt.Println("Consumer joined group! =)")
	}

	return nil
}

func (cg *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	if config.CurrentEnv == config.Development {
		fmt.Println("Consumer left group.")
	}

	return nil
}

func (cg *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var wg sync.WaitGroup

	go func() {
		defer wg.Done()

		for msg := range claim.Messages() {
			var decodedMsg stream.MessagePayload
			err := json.Unmarshal(msg.Value, &decodedMsg)
			if err != nil {
				// FIXME - logger
				fmt.Println(err)
				continue
			}

			go func() {
				cg.localCh <- eventReceivers{
					event: stream.Event{
						Name: string(msg.Key),
						Data: decodedMsg.Data,
					},
					receivers: decodedMsg.Receivers,
				}
			}()

			sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
		}
	}()
	wg.Add(1)
	wg.Wait()

	return nil
}

type BroadcastConsumer interface {
	SubscribeForUser(userID model.UserID, ch chan stream.Event) error
}

type eventReceivers struct {
	event     stream.Event
	receivers []stream.Receiver
}

type subscriber struct {
	sch    chan stream.Event
	userID model.UserID
}

type broadcastConsumer struct {
	encoder     stream.MessageEncoder
	consumer    sarama.ConsumerGroup
	subscribers []subscriber
}

func NewBroadcastConsumer(ctx context.Context, brokers []string, saramaConfig sarama.Config) (BroadcastConsumer, error) {
	consumer, err := sarama.NewConsumerGroup(brokers, "consumer-group-default", &saramaConfig)
	if err != nil {
		return nil, err
	}

	bc := &broadcastConsumer{
		consumer: consumer,
		encoder:  stream.NewMessageJsonEncoder(),
	}

	// Create local channel get all of the incoming events (consuming)
	localCh := make(chan eventReceivers)

	go func() {
		for {
			m := <-localCh

			for _, s := range bc.subscribers {
				if slices.Contains(m.receivers, stream.Receiver{UUID: s.userID}) {
					s.sch <- m.event
				}
			}
		}
	}()

	go func() {
		handler := &consumerGroupHandler{stream.NewMessageJsonEncoder(), localCh}
		err = consumer.Consume(context.TODO(), []string{stream.KafkaTopics().ChatTopic}, handler)
		if err != nil {
			// FIXME - logger
			fmt.Println(err)
		}
	}()

	return bc, nil
}

func (c *broadcastConsumer) SubscribeForUser(userID model.UserID, ch chan stream.Event) error {
	c.subscribers = append(c.subscribers, subscriber{
		userID: userID,
		sch:    ch,
	})

	return nil
}
