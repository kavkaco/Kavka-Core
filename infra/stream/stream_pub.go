package stream

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

const subjEvent = "events"

type StreamPublisher interface {
	Publish(event StreamEvent) error
}

type pub struct {
	nc *nats.Conn
}

func NewStreamPublisher(nc *nats.Conn) (StreamPublisher, error) {
	return &pub{nc}, nil
}

func (p *pub) Publish(event StreamEvent) error {
	dataBytes, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	err = p.nc.Publish(subjEvent, dataBytes)
	if err != nil {
		return err
	}

	p.nc.Flush()

	return nil
}
