package stream

import (
	"time"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/nats-io/nats.go"
)

func NewNATSAdapter(config *config.Config, logger *log.SubLogger) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectHandler(func(c *nats.Conn) {
			logger.Error("nats stream publisher disconnected")
		}),
		nats.ConnectHandler(func(c *nats.Conn) {
			logger.Info("nats stream publisher connected")
		}),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			logger.Error("nats raised a error: " + err.Error())
		}),
	}

	nc, err := nats.Connect(config.Nats.Url, opts...)
	if err != nil {
		return nil, err
	}

	return nc, err
}
