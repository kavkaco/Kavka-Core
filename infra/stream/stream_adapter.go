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

func NewNATSAdapter(config *config.Nats, logger *log.SubLogger) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectHandler(func(c *nats.Conn) {
			logger.Error("nats stream publisher disconnected")
		}),
		nats.ConnectHandler(func(c *nats.Conn) {
			logger.Info("nats stream publisher connected")
		}),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			logger.Error("nats raised an error: " + err.Error())
		}),
	}

	nc, err := nats.Connect(config.Url, opts...)
	if err != nil {
		return nil, err
	}

	return nc, err
}

func GetNATSTestInstance(callback func(*nats.Conn)) {
	var conn *nats.Conn

	dockerContainerEnvVariables := []string{}

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

	resource, err := pool.Run("nats", "latest", dockerContainerEnvVariables)
	if err != nil {
		stdlog.Fatalf("Could not start resource: %s", err)
	}

	ipAddr := resource.Container.NetworkSettings.IPAddress + ":4222"

	// Kill the container
	defer func() {
		if err = pool.Purge(resource); err != nil {
			stdlog.Fatalf("Could not purge resource: %s", err)
		}
	}()

	err = pool.Retry(func() error {
		logger := log.NewSubLogger("nats-test-instance")

		conn, err = NewNATSAdapter(&config.Nats{
			Url: ipAddr,
		}, logger)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		stdlog.Fatalf("Could not connect to nats: %s", err)
	}

	fmt.Printf("Docker nats container network ip address: %s\n\n", ipAddr)

	callback(conn)

	conn.Close()
}
