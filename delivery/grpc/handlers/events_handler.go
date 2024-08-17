package grpc_handlers

import (
	"context"

	"connectrpc.com/connect"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/log"
	eventsv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1/eventsv1connect"
)

type eventsHandler struct {
	logger   *log.SubLogger
	streamer stream.StreamSubscriber
}

func NewEventsGrpcHandler(logger *log.SubLogger, streamer stream.StreamSubscriber) eventsv1connect.EventsServiceHandler {
	return &eventsHandler{logger, streamer}
}

func (e *eventsHandler) SubscribeEventsStream(ctx context.Context, req *connect.Request[eventsv1.SubscribeEventsStreamRequest], str *connect.ServerStream[eventsv1.SubscribeEventsStreamResponse]) error {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)

	done := ctx.Done()
	userCh := make(chan *eventsv1.SubscribeEventsStreamResponse)
	e.streamer.UserSubscribe(userID, userCh)

	e.logger.Debug("user stream established")

	for {
		if str == nil {
			e.logger.Error("user stream is closed")
			return nil
		}

		select {
		case <-done:
			e.logger.Debug("user disconnected!")
			e.streamer.UserUnsubscribe(userID)
			return nil
		case event, ok := <-userCh:
			if !ok {
				e.logger.Error("user channel closed in user-subscribe method")
				continue
			}

			e.logger.Debug("events-handler", "event-name", event.Name)

			err := str.Send(event)
			if err != nil {
				log.Error("unable to send message with grpc: " + err.Error())
				continue
			}
		}
	}
}
