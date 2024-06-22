package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	lorem "github.com/bozaro/golorem"
	chatv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/chat/v1/chatv1connect"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ChatTestSuite struct {
	suite.Suite
	client chatv1connect.ChatServiceClient
	l      lorem.Lorem
}

func (s *ChatTestSuite) SetupSuite() {
	s.client = chatv1connect.NewChatServiceClient(http.DefaultClient, BaseUrl)
}

func (s *ChatTestSuite) TestA_CreateChannel() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	req := connect.NewRequest(&chatv1.CreateChannelRequest{
		Title:       "My channel",
		Username:    "my_channel",
		Description: "Sample channel only for test",
	})
	req.Header().Set("X-Access-Token", "sample000")

	resp, err := s.client.CreateChannel(ctx, req)
	require.NoError(s.T(), err)
	s.T().Log(resp)
}

func TestChatSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ChatTestSuite))
}
