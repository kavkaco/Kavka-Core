package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	lorem "github.com/bozaro/golorem"
	authv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1"
	"github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/auth/v1/authv1connect"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	client authv1connect.AuthServiceClient

	name, lastName, username, email, password   string
	verifyEmailToken, accessToken, refreshToken string //nolint
}

func (s *AuthTestSuite) SetupSuite() {
	// Generate random user info
	l := lorem.New()
	s.name = l.FirstName(0)
	s.lastName = l.LastName()
	s.username = l.Word(3, 10)
	s.email = l.Email()
	s.password = l.Word(10, 30)

	s.client = authv1connect.NewAuthServiceClient(http.DefaultClient, BaseUrl)
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
	require.NotEmpty(s.T(), resp.Msg.VerifyEmailToken)

	s.verifyEmailToken = resp.Msg.VerifyEmailToken
}

func (s *AuthTestSuite) TestB_VerifyEmail() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	_, err := s.client.VerifyEmail(ctx, &connect.Request[authv1.VerifyEmailRequest]{
		Msg: &authv1.VerifyEmailRequest{
			VerifyEmailToken: s.verifyEmailToken,
		},
	})
	require.NoError(s.T(), err)
}

func (s *AuthTestSuite) TestC_Login() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	resp, err := s.client.Login(ctx, &connect.Request[authv1.LoginRequest]{
		Msg: &authv1.LoginRequest{
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

func (s *AuthTestSuite) TestD_ChangePassword() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	newPassword := "new-password"
	_, err := s.client.ChangePassword(ctx, &connect.Request[authv1.ChangePasswordRequest]{
		Msg: &authv1.ChangePasswordRequest{
			AccessToken: s.accessToken,
			OldPassword: s.password,
			NewPassword: newPassword,
		},
	})
	require.NoError(s.T(), err)

	// Login again with new password to be sure that's changed!
	loginResp, err := s.client.Login(ctx, &connect.Request[authv1.LoginRequest]{
		Msg: &authv1.LoginRequest{
			Email:    s.email,
			Password: newPassword,
		},
	})

	require.NoError(s.T(), err)

	require.Equal(s.T(), loginResp.Msg.User.Name, s.name)
	require.Equal(s.T(), loginResp.Msg.User.LastName, s.lastName)
	require.Equal(s.T(), loginResp.Msg.User.Username, s.username)
	require.Equal(s.T(), loginResp.Msg.User.Email, s.email)
	s.password = newPassword
}

func (s *AuthTestSuite) TestE_Authenticate() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	resp, err := s.client.Authenticate(ctx, &connect.Request[authv1.AuthenticateRequest]{
		Msg: &authv1.AuthenticateRequest{
			AccessToken: s.accessToken,
		},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.Msg.User.Name, s.name)
	require.Equal(s.T(), resp.Msg.User.LastName, s.lastName)
	require.Equal(s.T(), resp.Msg.User.Username, s.username)
	require.Equal(s.T(), resp.Msg.User.Email, s.email)
}

func (s *AuthTestSuite) TestF_RefreshToken() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	resp, err := s.client.RefreshToken(ctx, &connect.Request[authv1.RefreshTokenRequest]{
		Msg: &authv1.RefreshTokenRequest{
			AccessToken:  s.accessToken,
			RefreshToken: s.refreshToken,
		},
	})
	require.NoError(s.T(), err)
	s.accessToken = resp.Msg.AccessToken
}

func (s *AuthTestSuite) TestG_SendResetPasswordVerification() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	resp, err := s.client.SendResetPasswordVerification(ctx, &connect.Request[authv1.SendResetPasswordVerificationRequest]{
		Msg: &authv1.SendResetPasswordVerificationRequest{
			Email: s.email,
		},
	})
	require.NoError(s.T(), err)
	s.verifyEmailToken = resp.Msg.VerifyEmailToken
}

func (s *AuthTestSuite) TestH_SubmitResetPassword() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	newPassword := "new-password2"
	_, err := s.client.SubmitResetPassword(ctx, &connect.Request[authv1.SubmitResetPasswordRequest]{
		Msg: &authv1.SubmitResetPasswordRequest{
			Token:       s.verifyEmailToken,
			NewPassword: newPassword,
		},
	})
	require.NoError(s.T(), err)
	// Login again with new password to be sure that's changed!
	loginResp, err := s.client.Login(ctx, &connect.Request[authv1.LoginRequest]{
		Msg: &authv1.LoginRequest{
			Email:    s.email,
			Password: newPassword,
		},
	})

	require.NoError(s.T(), err)

	require.Equal(s.T(), loginResp.Msg.User.Name, s.name)
	require.Equal(s.T(), loginResp.Msg.User.LastName, s.lastName)
	require.Equal(s.T(), loginResp.Msg.User.Username, s.username)
	require.Equal(s.T(), loginResp.Msg.User.Email, s.email)
	s.password = newPassword
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
