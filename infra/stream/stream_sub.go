package stream

import (
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/log"
	eventsv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const eventStreamSubject = "events"

type StreamSubscriber interface {
	UserSubscribe(userID model.UserID, userCh chan *eventsv1.EventStreamResponse)
	UserUnsubscribe(userID model.UserID)
}

type sub struct {
	nc              *nats.Conn
	logger          *log.SubLogger
	subscribedUsers []StreamSubscribedUser
}

func NewStreamSubscriber(nc *nats.Conn, logger *log.SubLogger) (StreamSubscriber, error) {
	subInstance := &sub{nc, logger, []StreamSubscribedUser{}}

	_, err := nc.Subscribe(eventStreamSubject, func(msg *nats.Msg) {
		go func() {
			var event eventsv1.StreamEvent
			err := proto.Unmarshal(msg.Data, &event)
			if err != nil { // || msgMap["payload"] == nil
				logger.Error("proto unmarshal error when decoding incoming msg of the broker: " + err.Error())
				return
			}

			var payload eventsv1.EventStreamResponse
			err = proto.Unmarshal(event.Payload, &payload)
			if err != nil {
				logger.Error("proto unmarshal error when decoding msg payload of the broker event: " + err.Error())
				return
			}

			// // Broadcast event to receivers by their pipe
			for _, receiverUserID := range event.ReceiversUserId {
				if su := MatchUserSubscription(receiverUserID, subInstance.subscribedUsers); su != nil {
					if su.UserPipe == nil {
						logger.Error("global event stream skipped broken user pipe")
						continue
					}

					go func() {
						su.UserPipe <- &payload
					}()
				}
			}
		}()
	})
	if err != nil {
		return nil, err
	}

	err = nc.Flush()
	if err != nil {
		logger.Error("nats flush error: " + err.Error())
	}

	return subInstance, nil
}

func (p *sub) UserSubscribe(userID model.UserID, userCh chan *eventsv1.EventStreamResponse) {
	p.logger.Debug("user stream established")
	p.subscribedUsers = append(p.subscribedUsers, StreamSubscribedUser{UserID: userID, UserPipe: userCh})
}

func (p *sub) UserUnsubscribe(userID model.UserID) {
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
