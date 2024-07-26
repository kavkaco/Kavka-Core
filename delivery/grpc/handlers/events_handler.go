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

const maximumConnectionErrorCount = 5

type eventsHandler struct {
	logger   *log.SubLogger
	streamer stream.StreamSubscriber
}

func NewEventsGrpcHandler(logger *log.SubLogger, streamer stream.StreamSubscriber) eventsv1connect.EventsServiceHandler {
	return &eventsHandler{logger, streamer}
}

func (e *eventsHandler) SubscribeEventsStream(ctx context.Context, req *connect.Request[eventsv1.EventStreamRequest], str *connect.ServerStream[eventsv1.EventStreamResponse]) error {
	userID := ctx.Value(interceptor.UserID{}).(model.UserID)

	userCh := make(chan *eventsv1.EventStreamResponse)
	e.streamer.UserSubscribe(userID, userCh)

	occurredErrorsCount := 0

	for {
		if str == nil {
			e.logger.Error("user stream is closed")
			return nil
		}

		event, ok := <-userCh
		if !ok {
			e.logger.Error("user channel closed in user-subscribe method")
			return nil
		}

		e.logger.Debug("events-handler", "event-name")

		err := str.Send(event)
		if err != nil {
			occurredErrorsCount++
		}

		// resource releasing after achieving maximum error count
		if occurredErrorsCount >= maximumConnectionErrorCount {
			e.logger.Debug("stream resource released")
			e.streamer.UserUnsubscribe(userID)
			close(userCh)
			return nil
		}
	}
}
