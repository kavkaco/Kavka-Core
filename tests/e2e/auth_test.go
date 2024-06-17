package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	lorem "github.com/bozaro/golorem"
	authv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/proto/auth/v1"
	"github.com/kavkaco/Kavka-ProtoBuf/gen/go/proto/auth/v1/authv1connect"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	client authv1connect.AuthServiceClient

	name, lastName, username, email, password string
	accessToken, refreshToken                 string //nolint
}

func (s *AuthTestSuite) SetupSuite() {
	// Generate random user info
	l := lorem.New()
	s.name = l.FirstName(0)
	s.lastName = l.LastName()
	s.username = l.Word(3, 10)
	s.email = l.Email()
	s.password = l.Word(10, 30)

	s.client = authv1connect.NewAuthServiceClient(http.DefaultClient, "http://127.0.0.1:8000")
}

func (s *AuthTestSuite) TestA_Register() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	resp, err := s.client.Register(ctx, &connect.Request[authv1.RegisterRequest]{
		Msg: &authv1.RegisterRequest{
			Name:     s.name,
			LastName: s.lastName,
			Username: s.username,
			Email:    s.email,
			Password: s.password,
		},
	})
	require.NoError(s.T(), err)

	require.Equal(s.T(), resp.Msg.User.Name, s.name)
	require.Equal(s.T(), resp.Msg.User.LastName, s.lastName)
	require.Equal(s.T(), resp.Msg.User.Username, s.username)
	require.Equal(s.T(), resp.Msg.User.Email, s.email)
}

// func (s *AuthTestSuite) TestB_Login() {
// 	ctx, _ := context.WithTimeout(context.TODO(), time.Second*5)

// 	resp, err := s.client.Login(ctx, &connect.Request[authv1.LoginRequest]{
// 		Msg: &authv1.LoginRequest{
// 			Email:    "mr.tahadostifam@gmail.com",
// 			Password: "12345678",
// 		},
// 	})
// 	require.NoError(s.T(), err)

// 	s.T().Log(resp.Msg.User)
// }

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
