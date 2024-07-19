package infra_consumers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	stream_consumers "github.com/kavkaco/Kavka-Core/infra/stream/consumers"
	"github.com/stretchr/testify/require"
)

func Test_Chat(t *testing.T) {
	kafkaConfig := config.Read().Kafka

	messageEncoder := stream.NewMessageJsonEncoder()

	eventsChan, err := stream_consumers.NewChatStreamConsumer(context.TODO(), kafkaConfig, messageEncoder)
	require.NoError(t, err)

	for event := range eventsChan {
		fmt.Println(event)
	}
}
