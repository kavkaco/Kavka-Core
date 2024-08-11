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

	// User
	userID          model.UserID
	email, password string

	// Tokens
	verifyEmailToken   string
	accessToken        string
	refreshToken       string
	resetPasswordToken string

	// Urls
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

	// FIXME
	s.verifyEmailRedirectUrl = "example.com"
	s.resetPasswordRedirectUrl = "example.com"

	hashManager := hash.NewHashManager(hash.DefaultHashParams)
	s.service = service.NewAuthService(authRepo, userRepo, authManager, hashManager, emailService)
}

func (s *AuthTestSuite) quickLogin(email string, password string) {
	ctx := context.TODO()

	user, _, _, varror := s.service.Login(ctx, email, password)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), user)
	require.Equal(s.T(), user.Email, email)
}

func (s *AuthTestSuite) TestA_Register() {
	ctx := context.TODO()

	user := userTestModels[0]
	s.email = user.Email
	s.password = "plain-password"

	verifyEmailToken, varror := s.service.Register(ctx, user.Name, user.LastName, user.Username, user.Email, s.password, s.verifyEmailRedirectUrl)
	require.Nil(s.T(), varror)

	s.verifyEmailToken = verifyEmailToken
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

	newPassword := "password-must-be-changed"

	varror := s.service.ChangePassword(ctx, s.userID, s.password, newPassword)
	require.Nil(s.T(), varror)

	s.quickLogin(s.email, newPassword)
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

	newPassword := "password-reset-must-work"

	varror := s.service.SubmitResetPassword(ctx, s.resetPasswordToken, newPassword)
	require.Nil(s.T(), varror)

	s.quickLogin(s.email, newPassword)
	s.password = newPassword
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
