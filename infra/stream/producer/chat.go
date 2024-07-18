package stream_producer

import (
	"github.com/IBM/sarama"
	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const topicName = "chats"

type ChatProducer interface {
	MessageSent(chatID model.ChatID, messageID model.MessageID, message model.Message) error
	MessageDeleted(chatID model.ChatID, messageID model.MessageID) error
	ChatCreated(eventReceivers []model.UserID, chat model.Chat) error
	ChatDeleted(eventReceivers []model.UserID, chatID model.ChatID) error
}

type producer struct {
	kafkaConfig    *config.Kafka
	producer       sarama.AsyncProducer
	messageEncoder stream.MessageEncoder
}

func NewChatStreamProducer(kafkaConfig *config.Kafka) (ChatProducer, error) {
	p, err := sarama.NewAsyncProducer(kafkaConfig.Brokers, kafkaConfig.Sarama)
	if err != nil {
		return nil, err
	}
	// defer p.AsyncClose()

	// go func() {
	// 	errs := p.Errors()

	// 	for err := range errs {
	// 		log.Fatalln(err)
	// 	}
	// }()

	messageEncoder := stream.NewMessageJsonEncoder()

	return &producer{kafkaConfig, p, messageEncoder}, nil
}

func (p *producer) ChatCreated(eventReceivers []string, chat model.Chat) error {
	encodedModel, err := p.messageEncoder.Encode(chat)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic:     topicName,
		Partition: 0,
		Offset:    0,
		Key:       sarama.StringEncoder("chat_created"),
		Value:     sarama.ByteEncoder(encodedModel),
	}

	p.producer.Input() <- &msg

	return nil
}

func (p *producer) ChatDeleted(eventReceivers []string, chatID primitive.ObjectID) error {
	panic("unimplemented")
}

func (p *producer) MessageDeleted(chatID primitive.ObjectID, messageID primitive.ObjectID) error {
	panic("unimplemented")
}

func (p *producer) MessageSent(chatID primitive.ObjectID, messageID primitive.ObjectID, message model.Message) error {
	panic("unimplemented")
}
