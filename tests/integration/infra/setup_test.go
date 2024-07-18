package infra

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
)

var kafkaConfig config.Kafka

func TestMain(m *testing.M) {
	config := config.Read()
	kafkaConfig = config.Kafka

	m.Run()
}
