package service

import (
	"context"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	service "github.com/kavkaco/Kavka-Core/internal/service/auth"
	"github.com/kavkaco/Kavka-Core/pkg/email"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	auth_manager "github.com/tahadostifam/go-auth-manager"
)

type AuthTestSuite struct {
	suite.Suite
	service service.AuthService

	userID                   model.UserID
	verifyEmailToken         string
	email, password          string
	accessToken              string
	refreshToken             string
	resetPasswordToken       string
	verifyEmailRedirectUrl   string
	resetPasswordRedirectUrl string
}

func (s *AuthTestSuite) SetupSuite() {
	authRepo := repository_mongo.NewAuthMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: "private-key",
	})

	emailService := email.NewEmailDevelopmentService()

	s.verifyEmailRedirectUrl = "example.com"
	s.resetPasswordRedirectUrl = "example.com"

	hashManager := hash.NewHashManager(hash.DefaultHashParams)
	s.service = service.NewAuthService(authRepo, userRepo, authManager, hashManager, emailService)
}

func (s *AuthTestSuite) TestA_Register() {
	ctx := context.TODO()

	name := "John"
	lastName := "Doe"
	username := "john_doe"
	email := "john_doe@kavka.org"
	password := "12345678"

	verifyEmailToken, varror := s.service.Register(ctx, name, lastName, username, email, password, s.verifyEmailRedirectUrl)
	require.Nil(s.T(), varror)

	s.verifyEmailToken = verifyEmailToken
	s.email = email
	s.password = password
}

func (s *AuthTestSuite) TestB_VerifyEmail() {
	ctx := context.TODO()

	varror := s.service.VerifyEmail(ctx, s.verifyEmailToken)
	require.Nil(s.T(), varror)
}

func (s *AuthTestSuite) TestC_Login() {
	ctx := context.TODO()

	user, accessToken, refreshToken, varror := s.service.Login(ctx, s.email, s.password)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), accessToken)
	require.NotEmpty(s.T(), refreshToken)
	require.NotEmpty(s.T(), user)
	require.Equal(s.T(), user.Email, s.email)

	s.accessToken = accessToken
	s.refreshToken = refreshToken
	s.userID = user.UserID
}

func (s *AuthTestSuite) TestD_ChangePassword() {
	ctx := context.TODO()

	newPassword := "12341234"

	varror := s.service.ChangePassword(ctx, s.userID, s.password, newPassword)
	require.Nil(s.T(), varror)

	// Login again with new password to be sure that's changed!

	user, _, _, varror := s.service.Login(ctx, s.email, newPassword)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), user)
	require.Equal(s.T(), user.Email, s.email)

	s.password = newPassword
}

func (s *AuthTestSuite) TestE_Authenticate() {
	ctx := context.TODO()

	user, varror := s.service.Authenticate(ctx, s.accessToken)
	require.Nil(s.T(), varror)

	require.Equal(s.T(), user.Email, s.email)
}

func (s *AuthTestSuite) TestF_RefreshToken() {
	ctx := context.TODO()

	accessToken, varror := s.service.RefreshToken(ctx, s.userID, s.refreshToken)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), accessToken)
	require.NotEqual(s.T(), accessToken, s.accessToken)

	s.accessToken = accessToken
}

func (s *AuthTestSuite) TestG_SendResetPassword() {
	ctx := context.TODO()

	resetPasswordToken, timeout, varror := s.service.SendResetPassword(ctx, s.email, s.resetPasswordRedirectUrl)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), timeout)
	require.NotEmpty(s.T(), resetPasswordToken)

	s.resetPasswordToken = resetPasswordToken
}

func (s *AuthTestSuite) TestH_SubmitResetPassword() {
	ctx := context.TODO()

	newPassword := "98769876"

	varror := s.service.SubmitResetPassword(ctx, s.resetPasswordToken, newPassword)
	require.Nil(s.T(), varror)

	// Login again with new password to be sure that's changed!

	user, _, _, varror := s.service.Login(ctx, s.email, newPassword)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), user)
	require.Equal(s.T(), user.Email, s.email)

	s.password = newPassword
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
