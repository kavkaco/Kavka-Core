package grpc_service

import "github.com/kavkaco/Kavka-Core/internal/service/chat"

type ChatGrpcServer struct {
	chatService chat.ChatService
}

func NewChatGrpcHandler(chatService chat.ChatService) ChatGrpcServer {
	return ChatGrpcServer{chatService}
}
