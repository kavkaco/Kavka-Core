package service

import (
	"context"
	"testing"

	lorem "github.com/bozaro/golorem"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	service "github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/pkg/auth_manager"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	service service.AuthService
	lem     *lorem.Lorem

	verifyEmailToken   string
	email, password    string
	accessToken        string
	refreshToken       string
	resetPasswordToken string
}

func (s *AuthTestSuite) SetupSuite() {
	s.lem = lorem.New()

	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: "private-key",
	})
	hashManager := hash.NewHashManager(hash.DefaultHashParams)
	s.service = service.NewAuthService(authRepo, userRepo, authManager, hashManager)
}

func (s *AuthTestSuite) TestA_Register() {
	ctx := context.TODO()

	name := s.lem.FirstName(0)
	lastName := s.lem.LastName()
	username := s.lem.Word(1, 10)
	email := s.lem.Email()
	password := "strong-password"

	user, verifyEmailToken, err := s.service.Register(ctx, name, lastName, username, email, password)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), verifyEmailToken)
	require.Equal(s.T(), user.Name, name)
	require.Equal(s.T(), user.LastName, lastName)
	require.Equal(s.T(), user.Username, username)
	require.Equal(s.T(), user.Email, email)

	s.verifyEmailToken = verifyEmailToken
	s.email = email
	s.password = password
}

func (s *AuthTestSuite) TestB_VerifyEmail() {
	ctx := context.TODO()

	err := s.service.VerifyEmail(ctx, s.verifyEmailToken)
	require.NoError(s.T(), err)
}

func (s *AuthTestSuite) TestC_Login() {
	ctx := context.TODO()

	user, accessToken, refreshToken, err := s.service.Login(ctx, s.email, s.password)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), accessToken)
	require.NotEmpty(s.T(), refreshToken)
	require.NotEmpty(s.T(), user)
	require.Equal(s.T(), user.Email, s.email)

	s.accessToken = accessToken
	s.refreshToken = refreshToken
}

func (s *AuthTestSuite) TestD_ChangePassword() {
	ctx := context.TODO()

	newPassword := "password-changed"

	err := s.service.ChangePassword(ctx, s.accessToken, s.password, newPassword)
	require.NoError(s.T(), err)

	// Login again with new password to be sure that's changed!

	user, _, _, err := s.service.Login(ctx, s.email, newPassword)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), user)
	require.Equal(s.T(), user.Email, s.email)

	s.password = newPassword
}

func (s *AuthTestSuite) TestE_Authenticate() {
	ctx := context.TODO()

	user, err := s.service.Authenticate(ctx, s.accessToken)
	require.NoError(s.T(), err)
	require.Equal(s.T(), user.Email, s.email)
}

func (s *AuthTestSuite) TestF_RefreshToken() {
	ctx := context.TODO()

	newAccessToken, err := s.service.RefreshToken(ctx, s.refreshToken, s.accessToken)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), newAccessToken)
	require.NotEqual(s.T(), newAccessToken, s.accessToken)
}

func (s *AuthTestSuite) TestG_SendResetPasswordVerification() {
	ctx := context.TODO()

	resetPasswordToken, timeout, err := s.service.SendResetPasswordVerification(ctx, s.email)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), timeout)
	require.NotEmpty(s.T(), resetPasswordToken)

	s.resetPasswordToken = resetPasswordToken
}

func (s *AuthTestSuite) TestH_SubmitResetPassword() {
	ctx := context.TODO()

	newPassword := "reset-password"

	err := s.service.SubmitResetPassword(ctx, s.resetPasswordToken, newPassword)
	require.NoError(s.T(), err)

	// Login again with new password to be sure that's changed!

	user, _, _, err := s.service.Login(ctx, s.email, newPassword)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), user)
	require.Equal(s.T(), user.Email, s.email)

	s.password = newPassword
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
