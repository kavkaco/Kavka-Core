package grpc_handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	stream_consumers "github.com/kavkaco/Kavka-Core/infra/stream/consumers"
	"github.com/kavkaco/Kavka-Core/internal/model"
	eventsv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1/eventsv1connect"
)

type eventsHandler struct {
	broadcastConsumer stream_consumers.BroadcastConsumer
}

func NewEventsGrpcHandler(broadcastConsumer stream_consumers.BroadcastConsumer) eventsv1connect.EventsServiceHandler {
	return &eventsHandler{broadcastConsumer}
}

func (e *eventsHandler) SubscribeEventsStream(ctx context.Context, req *connect.Request[eventsv1.EventStreamRequest], str *connect.ServerStream[eventsv1.EventStreamResponse]) error {
	userID := ctx.Value(interceptor.UserIDKey{}).(model.UserID)

	ch := make(chan stream.Event)

	go func() {
		e.broadcastConsumer.SubscribeForUser(userID, ch)
	}()

	for {
		event := <-ch

		dataJson, err := json.Marshal(event.Data)
		if err != nil {
			// FIXME - logger
			fmt.Println(err)
			continue
		}

		str.Send(&eventsv1.EventStreamResponse{
			Name:    event.Name,
			RawJson: string(dataJson),
		})
	}
}
