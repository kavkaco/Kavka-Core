package stream

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
)

type NATSAdapter struct {
	Conn        *nats.Conn
	JetStream   nats.JetStreamContext
	StreamConfig JetStreamConfig
}

func NewNATSAdapter(cfg *config.Nats, logger *log.SubLogger, jsCfg JetStreamConfig) (*NATSAdapter, error) {
	opts := []nats.Option{
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectHandler(func(c *nats.Conn) {
			logger.Error("nats stream publisher disconnected")
		}),
		nats.ConnectHandler(func(c *nats.Conn) {
			logger.Info("nats stream publisher connected")
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			logger.Info("nats reconnected")
		}),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			logger.Error("nats raised an error: " + err.Error())
		}),
	}

	nc, err := nats.Connect(cfg.Url, opts...)
	if err != nil {
		return nil, err
	}

	adapter := &NATSAdapter{
		Conn:         nc,
		StreamConfig: jsCfg,
	}

	if jsCfg.Enabled {
		js, err := nc.JetStream()
		if err != nil {
			return nil, fmt.Errorf("failed to create jetstream context: %w", err)
		}
		adapter.JetStream = js

		err = adapter.configureStream()
		if err != nil {
			logger.Error("failed to configure jetstream stream: " + err.Error())
		}
	}

	return adapter, nil
}

func (a *NATSAdapter) configureStream() error {
	stream, err := a.JetStream.StreamInfo(EventsStreamName)
	if err != nil {
		_, err = a.JetStream.AddStream(&nats.StreamConfig{
			Name:      EventsStreamName,
			Subjects:  []string{EventsStreamSubject},
			Storage:   a.StreamConfig.StorageType,
			MaxAge:    a.StreamConfig.MaxAge,
			MaxMsgs:   a.StreamConfig.MaxMsgs,
			MaxBytes:  a.StreamConfig.MaxBytes,
			Replicas:  a.StreamConfig.Replicas,
		})
		if err != nil {
			return fmt.Errorf("failed to add jetstream stream: %w", err)
		}
	}

	_ = stream
	return nil
}

func (a *NATSAdapter) Close() {
	if a.Conn != nil {
		a.Conn.Close()
	}
}

func GetNATSTestInstance(callback func(*NATSAdapter)) {
	var adapter *NATSAdapter

	err := os.Setenv("ENV", "test")
	if err != nil {
		stdlog.Fatalf("Could not set the environment variable to test: %s", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		stdlog.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		stdlog.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("nats", "latest", []string{})
	if err != nil {
		stdlog.Fatalf("Could not start resource: %s", err)
	}

	ipAddr := resource.Container.NetworkSettings.IPAddress + ":4222"

	defer func() {
		if err = pool.Purge(resource); err != nil {
			stdlog.Fatalf("Could not purge resource: %s", err)
		}
	}()

	err = pool.Retry(func() error {
		logger := log.NewSubLogger("nats-test-instance")

		adapter, err = NewNATSAdapter(&config.Nats{
			Url: ipAddr,
		}, logger, DefaultJetStreamConfig())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		stdlog.Fatalf("Could not connect to nats: %s", err)
	}

	fmt.Printf("Docker nats container network ip address: %s\n\n", ipAddr)

	callback(adapter)

	adapter.Close()
}
