package grpc

import (
	"net/http"

	"connectrpc.com/connect"
	grpc_handlers "github.com/kavkaco/Kavka-Core/delivery/grpc/handlers"
	"github.com/kavkaco/Kavka-Core/delivery/grpc/interceptor"
	"github.com/kavkaco/Kavka-Core/infra/stream"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/internal/service/chat"
	"github.com/kavkaco/Kavka-Core/internal/service/message"
	"github.com/kavkaco/Kavka-Core/internal/service/search"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1/authv1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1/chatv1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/events/v1/eventsv1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/message/v1/messagev1connect"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/search/v1/searchv1connect"
)

type Services struct {
	AuthService      *auth.AuthService
	ChatService      *chat.ChatService
	MessageService   *message.MessageService
	SearchService    *search.SearchService
	StreamSubscriber stream.StreamSubscriber
}

func NewGrpcServer(router *http.ServeMux, services *Services) {
	authInterceptor := interceptor.NewAuthInterceptor(services.AuthService)
	interceptors := connect.WithInterceptors(authInterceptor)

	authGrpcHandler := grpc_handlers.NewAuthGrpcHandler(services.AuthService)
	authGrpcRoute, authGrpcRouter := authv1connect.NewAuthServiceHandler(authGrpcHandler)

	chatGrpcHandler := grpc_handlers.NewChatGrpcHandler(log.NewSubLogger("chats-handler"), services.ChatService)
	chatGrpcRoute, chatGrpcRouter := chatv1connect.NewChatServiceHandler(chatGrpcHandler, interceptors)

	eventsGrpcHandler := grpc_handlers.NewEventsGrpcHandler(log.NewSubLogger("events-handler"), services.StreamSubscriber)
	eventsGrpcRoute, eventsGrpcRouter := eventsv1connect.NewEventsServiceHandler(eventsGrpcHandler, interceptors)

	messageGrpcHandler := grpc_handlers.NewMessageGrpcHandler(log.NewSubLogger("message-handler"), services.MessageService)
	messageGrpcRoute, messageGrpcRouter := messagev1connect.NewMessageServiceHandler(messageGrpcHandler, interceptors)

	searchGrpcHandler := grpc_handlers.NewSearchGrpcHandler(log.NewSubLogger("message-handler"), services.SearchService)
	searchGrpcRoute, searchGrpcRouter := searchv1connect.NewSearchServiceHandler(searchGrpcHandler, interceptors)

	router.Handle(authGrpcRoute, authGrpcRouter)
	router.Handle(chatGrpcRoute, chatGrpcRouter)
	router.Handle(eventsGrpcRoute, eventsGrpcRouter)
	router.Handle(messageGrpcRoute, messageGrpcRouter)
	router.Handle(searchGrpcRoute, searchGrpcRouter)
}
