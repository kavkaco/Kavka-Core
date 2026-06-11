package stream

import (
	"errors"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
	eventsv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/events/v1"
	"github.com/nats-io/nats.go"
)

var (
	ErrPublishEvent       = errors.New("publishing event went wrong")
	ErrStreamNotConfigured = errors.New("nats jetstream not configured")
)

type StreamSubscribedUser struct {
	UserID   model.UserID
	UserPipe chan *eventsv1.SubscribeEventsStreamResponse
}

const (
	EventsStreamName    = "kavka-events"
	EventsStreamSubject = "events.>"
	MaxEventsAge        = 7 * 24 * time.Hour
)

type JetStreamConfig struct {
	Enabled      bool
	StorageType  nats.StorageType
	MaxAge       time.Duration
	MaxMsgs      int64
	MaxBytes     int64
	Replicas     int
}

func DefaultJetStreamConfig() JetStreamConfig {
	return JetStreamConfig{
		Enabled:     true,
		StorageType: nats.FileStorage,
		MaxAge:      MaxEventsAge,
		MaxMsgs:     -1,
		MaxBytes:    -1,
		Replicas:    1,
	}
}
