package infra_consumers_test

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
)

var kafkaConfig config.Kafka

func TestMain(m *testing.M) {
	kafkaConfig = config.Read().Kafka

	m.Run()
}
