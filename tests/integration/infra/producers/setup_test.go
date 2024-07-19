package infra_producers_test

import (
	"testing"
	"time"

	"github.com/kavkaco/Kavka-Core/config"
)

var kafkaConfig config.Kafka

func TestMain(m *testing.M) {
	kafkaConfig = config.Read().Kafka

	m.Run()

	time.Sleep(2 * time.Second)
}
