package stream

import (
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/model"
	eventsv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1"
)

var ErrPublishEvent = errors.New("publishing event went wrong")

type StreamSubscribedUser struct {
	UserID   model.UserID
	UserPipe chan *eventsv1.EventStreamResponse
}
