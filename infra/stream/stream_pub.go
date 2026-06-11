package stream

import (
	"fmt"

	eventsv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/events/v1"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const subjEvent = "events"

type StreamPublisher interface {
	Publish(event *eventsv1.StreamEvent) error
}

type pub struct {
	nc  *nats.Conn
	js  nats.JetStreamContext
}

func NewStreamPublisher(adapter *NATSAdapter) (StreamPublisher, error) {
	if adapter == nil {
		return nil, fmt.Errorf("nats adapter is nil")
	}
	return &pub{nc: adapter.Conn, js: adapter.JetStream}, nil
}

func (p *pub) Publish(event *eventsv1.StreamEvent) error {
	eventBuf, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	if p.js != nil {
		_, err = p.js.Publish(subjEvent, eventBuf)
		return err
	}

	return p.nc.Publish(subjEvent, eventBuf)
}
