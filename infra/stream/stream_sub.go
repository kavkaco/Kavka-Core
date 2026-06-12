package stream

import (
	"fmt"
	"sync"
	"time"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/log"
	eventsv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/events/v1"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const eventStreamSubject = "events"
const eventStreamName = "KAVKA_EVENTS_STREAM"

type StreamSubscriber interface {
	UserSubscribe(userID model.UserID, userCh chan *eventsv1.SubscribeEventsStreamResponse)
	UserUnsubscribe(userID model.UserID)
}

type sub struct {
	nc              *nats.Conn
	js              nats.JetStreamContext
	logger          *log.SubLogger
	mu              sync.RWMutex
	subscribedUsers []StreamSubscribedUser
}

func NewStreamSubscriber(adapter *NATSAdapter, logger *log.SubLogger) (StreamSubscriber, error) {
	subInstance := &sub{
		nc:              adapter.Conn,
		js:              adapter.JetStream,
		logger:          logger,
		mu:              sync.RWMutex{},
		subscribedUsers: []StreamSubscribedUser{},
	}

	// create the stream if it doesn't exist
	if subInstance.js != nil {
		if err := ensureStreamExists(subInstance.js, logger); err != nil {
			logger.Error("failed to ensure stream exists: " + err.Error())
		}
	}

	subscribeFn := func(msg *nats.Msg) {
		go func() {
			var event eventsv1.StreamEvent
			err := proto.Unmarshal(msg.Data, &event)
			if err != nil {
				logger.Error("proto unmarshal error when decoding incoming msg of the broker: " + err.Error())
				return
			}

			var payload eventsv1.SubscribeEventsStreamResponse
			err = proto.Unmarshal(event.Payload, &payload)
			if err != nil {
				logger.Error("proto unmarshal error when decoding msg payload of the broker event: " + err.Error())
				return
			}

			subInstance.mu.RLock()
			for _, receiverUserID := range event.ReceiversUserId {
				for _, su := range subInstance.subscribedUsers {
					if su.UserID == receiverUserID {
						if su.UserPipe == nil {
							logger.Error("event stream skipped broken user pipe")
							continue
						}

						su.UserPipe <- &payload
					}
				}
			}
			subInstance.mu.RUnlock()

			if msg.Reply != "" {
				msg.Ack()
			}
		}()
	}

	subOpts := []nats.SubOpt{
		nats.ManualAck(),
		nats.DeliverNew(),
	}

	subOpts = append(subOpts, nats.Durable("kavka-events-subscriber"))

	if subInstance.js != nil {
		_, err := subInstance.js.QueueSubscribe(eventStreamSubject, "event-workers", subscribeFn, subOpts...)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := subInstance.nc.Subscribe(eventStreamSubject, subscribeFn)
		if err != nil {
			return nil, err
		}
	}

	return subInstance, nil
}

func (p *sub) UserSubscribe(userID model.UserID, userCh chan *eventsv1.SubscribeEventsStreamResponse) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribedUsers = append(p.subscribedUsers, StreamSubscribedUser{UserID: userID, UserPipe: userCh})
}

func (p *sub) UserUnsubscribe(userID model.UserID) {
	p.mu.Lock()
	defer p.mu.Unlock()
	idx := -1

	for i, su := range p.subscribedUsers {
		if su.UserID == userID {
			idx = i
			break
		}
	}

	if idx != -1 {
		p.subscribedUsers = append(p.subscribedUsers[:idx], p.subscribedUsers[idx+1:]...)
	}
}

func ensureStreamExists(js nats.JetStreamContext, logger *log.SubLogger) error {
	streamInfo, err := js.StreamInfo(eventStreamName)
	if err == nil {
		logger.Info(fmt.Sprintf("JetStream stream already exists: %s (subjects: %v)", eventStreamName, streamInfo.Config.Subjects))
		return nil
	}

	if err != nats.ErrStreamNotFound {
		return fmt.Errorf("failed to check stream info: %w", err)
	}

	logger.Info(fmt.Sprintf("Creating JetStream stream: %s for subject: %s", eventStreamName, eventStreamSubject))

	streamConfig := &nats.StreamConfig{
		Name:      eventStreamName,
		Subjects:  []string{eventStreamSubject},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxMsgs:   -1,                 // No limit on number of messages
		MaxBytes:  -1,                 // No limit on total bytes
		MaxAge:    7 * 24 * time.Hour, // Keep messages for 7 days
		Discard:   nats.DiscardOld,

		// FIXME: make this configurable
		Replicas: 1, // For development, use 1 replica

		Duplicates:  time.Minute, // Duplicate detection window
		AllowRollup: true,        // Allow rollup messages
	}

	_, err = js.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	logger.Info(fmt.Sprintf("JetStream stream created successfully: %s", eventStreamName))
	return nil
}
