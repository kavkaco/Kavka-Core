package stream

import (
	"errors"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

var ErrPublishEvent = errors.New("publishing event went wrong")

type StreamEvent struct {
	SenderUserID     model.UserID   `json:"senderUserId"`
	ReceiversUserIDs []model.UserID `json:"receiversUserId"`
	Name             string         `json:"name"`
	DataJson         string         `json:"dataJson"`
}

type StreamSubscribedUser struct {
	UserID   model.UserID
	UserPipe chan StreamEvent
}
