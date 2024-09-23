package grpc_handlers

import (
	"context"

	"connectrpc.com/connect"
	grpc_helpers "github.com/kavkaco/Kavka-Core/delivery/grpc/helpers"
	"github.com/kavkaco/Kavka-Core/internal/model/proto_model_transformer"
	"github.com/kavkaco/Kavka-Core/internal/service/search"
	"github.com/kavkaco/Kavka-Core/log"
	searchv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/search/v1"
	"github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/search/v1/searchv1connect"
)

type searchHandler struct {
	logger        *log.SubLogger
	searchService *search.SearchService
}

func NewSearchGrpcHandler(logger *log.SubLogger, searchService *search.SearchService) searchv1connect.SearchServiceHandler {
	return &searchHandler{logger, searchService}
}

func (s *searchHandler) Search(ctx context.Context, req *connect.Request[searchv1.SearchRequest]) (*connect.Response[searchv1.SearchResponse], error) {
	result, varror := s.searchService.Search(ctx, req.Msg.Input)
	if varror != nil {
		return nil, grpc_helpers.GrpcVarror(varror, connect.CodeUnavailable)
	}

	chats, err := proto_model_transformer.ChatsToProto(result.Chats)
	if err != nil {
		return nil, err
	}

	users := proto_model_transformer.UsersToProto(result.Users)

	res := &connect.Response[searchv1.SearchResponse]{
		Msg: &searchv1.SearchResponse{
			Result: &searchv1.SearchResponse_SearchResult{
				Chats: chats,
				Users: users,
			},
		},
	}

	return res, nil
}
