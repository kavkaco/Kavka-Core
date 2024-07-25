package stream

import (
	eventsv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const subjEvent = "events"

type StreamPublisher interface {
	Publish(event eventsv1.StreamEvent) error
}

type pub struct {
	nc *nats.Conn
}

func NewStreamPublisher(nc *nats.Conn) (StreamPublisher, error) {
	return &pub{nc}, nil
}

func (p *pub) Publish(event eventsv1.StreamEvent) error {
	eventBuf, err := proto.Marshal(&event)
	if err != nil {
		return err
	}

	err = p.nc.Publish(subjEvent, eventBuf)
	if err != nil {
		return err
	}

	p.nc.Flush()

	return nil
}
