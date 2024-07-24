package stream

import (
	"encoding/json"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/nats-io/nats.go"
)

const eventStreamSubject = "events"

type StreamSubscriber interface {
	UserSubscribe(userID model.UserID, userCh chan StreamEvent)
	UserUnsubscribe(userID model.UserID)
}

type sub struct {
	nc                *nats.Conn
	logger            *log.SubLogger
	globalEventStream chan StreamEvent
	subscribedUsers   []StreamSubscribedUser
}

func NewStreamSubscriber(nc *nats.Conn, logger *log.SubLogger) (StreamSubscriber, error) {
	globalEventStream := make(chan StreamEvent)

	nc.Subscribe(eventStreamSubject, func(msg *nats.Msg) {
		go func() {
			var pe StreamEvent
			err := json.Unmarshal(msg.Data, &pe)
			if err != nil {
				logger.Error("unable to decode msg data in global subscribe of stream: " + err.Error())
				return
			}

			globalEventStream <- pe

		}()
	})
	err := nc.Flush()
	if err != nil {
		logger.Error("nats flush error: " + err.Error())
	}

	subInstance := &sub{nc, logger, globalEventStream, []StreamSubscribedUser{}}

	// Matcher engine
	go func() {
		for {
			pe := <-globalEventStream

			// // Broadcast event to receivers by their pipe
			for _, receiverUserID := range pe.ReceiversUserIDs {
				if su := MatchUserSubscription(receiverUserID, subInstance.subscribedUsers); su != nil {
					if su.UserPipe == nil {
						logger.Error("global event stream skipped broken user pipe")
						continue
					}

					go func() {
						su.UserPipe <- pe
					}()
				}
			}
		}
	}()

	return subInstance, nil
}

func (p *sub) UserSubscribe(userID model.UserID, userCh chan StreamEvent) {
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
