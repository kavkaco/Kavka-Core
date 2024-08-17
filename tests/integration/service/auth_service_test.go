package service

import (
	"context"
	"testing"
	"time"

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
	user            *model.User
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

	//auth manger
	authManager auth_manager.AuthManager
}

func (s *AuthTestSuite) SetupSuite() {
	authRepo := repository_mongo.NewAuthMongoRepository(db)
	userRepo := repository_mongo.NewUserMongoRepository(db)
	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: "private-key",
	})

	emailService := email.NewEmailDevelopmentService()
	s.authManager = authManager

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

func (s *AuthTestSuite) TestAA_Register() {
	ctx := context.TODO()

	user := model.User{
		Name:     "User1:Name",
		LastName: "User1:LastName",
		Email:    "user1@kavka.org",
		Username: "user1",
	}
	s.email = user.Email
	s.password = "plain-password"

	verifyEmailToken, varror := s.service.Register(ctx, user.Name, user.LastName, user.Username, user.Email, s.password, s.verifyEmailRedirectUrl)
	require.Nil(s.T(), varror)

	s.verifyEmailToken = verifyEmailToken
}

// Ok
func (s *AuthTestSuite) TestAB_InvalidInputsRegister() {
	ctx := context.TODO()

	_, varror := s.service.Register(ctx, "", "", "", "", "", "")
	require.NotNil(s.T(), varror)
}

// Ok
func (s *AuthTestSuite) TestAC_DuplicatedEmailRegister() {
	ctx := context.TODO()

	user := model.User{
		Name:     "User2:Name",
		LastName: "User2:LastName",
		Email:    "user1@kavka.org",
		Username: "user2",
	}
	_, varror := s.service.Register(ctx, user.Name, user.LastName, user.Username, user.Email, s.password, s.verifyEmailRedirectUrl)
	require.Equal(s.T(), varror.Error, service.ErrEmailAlreadyExist)
}

// Ok
func (s *AuthTestSuite) TestAD_DuplicatedUsernameRegister() {
	ctx := context.TODO()

	user := model.User{
		Name:     "User2:Name",
		LastName: "User2:LastName",
		Email:    "user2@kavka.org",
		Username: "user1",
	}
	_, varror := s.service.Register(ctx, user.Name, user.LastName, user.Username, user.Email, s.password, s.verifyEmailRedirectUrl)
	require.Equal(s.T(), varror.Error, service.ErrUsernameAlreadyExist)
}

// Ok
func (s *AuthTestSuite) TestAE_EmailNotVerifiedLogin() {
	ctx := context.TODO()

	_, _, _, varror := s.service.Login(ctx, s.email, s.password)
	require.Equal(s.T(), varror.Error, service.ErrEmailNotVerified)
}

// New
func (s *AuthTestSuite) TestAF_NotVerifiedEmailPassword() {
	ctx := context.TODO()

	_, _, varror := s.service.SendResetPassword(ctx, s.email, s.resetPasswordRedirectUrl)
	require.Equal(s.T(), varror.Error, service.ErrEmailNotVerified)
}
func (s *AuthTestSuite) TestBA_VerifyEmail() {
	ctx := context.TODO()

	varror := s.service.VerifyEmail(ctx, s.verifyEmailToken)
	require.Nil(s.T(), varror)
}

// ok
func (s *AuthTestSuite) TestBB_InvalidInputVerifyEmail() {
	ctx := context.TODO()

	varror := s.service.VerifyEmail(ctx, "")
	require.NotNil(s.T(), varror)
}

// ok
func (s *AuthTestSuite) TestBC_InvalidTokenVerifyEmail() {
	ctx := context.TODO()

	varror := s.service.VerifyEmail(ctx, "invalid-token")
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)
}

// fail
// func (s *AuthTestSuite) TestBD_InvalidTokenUUIDVerifyEmail() {
// 	ctx := context.TODO()

// 	vet, err := s.authManager.GenerateToken(ctx, auth_manager.VerifyEmail,
// 		&auth_manager.TokenPayload{
// 			UUID:      "invalid",
// 			TokenType: auth_manager.VerifyEmail,
// 			CreatedAt: time.Now(),
// 		},
// 		service.VerifyEmailTokenExpr,
// 	)
// 	require.Nil(s.T(), err)

// 	varror := s.service.VerifyEmail(ctx, vet)
// 	require.Equal(s.T(), varror.Error, service.ErrVerifyEmail)
// }

func (s *AuthTestSuite) TestCA_Login() {
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

	s.user = user
}

// OK
func (s *AuthTestSuite) TestCB_InvalidInputLogin() {
	ctx := context.TODO()

	_, _, _, varror := s.service.Login(ctx, "", "")
	require.NotNil(s.T(), varror)
}

// OK
func (s *AuthTestSuite) TestCC_InvalidEmailLogin() {
	ctx := context.TODO()

	_, _, _, varror := s.service.Login(ctx, "invalidemail@gmail.com", s.password)
	require.Equal(s.T(), varror.Error, service.ErrInvalidEmailOrPassword)
}

// Ok
func (s *AuthTestSuite) TestCD_InvalidPasswordLogin() {
	ctx := context.TODO()

	_, _, _, varror := s.service.Login(ctx, s.email, "invalidPassword")
	require.Equal(s.T(), varror.Error, service.ErrInvalidEmailOrPassword)
}

// ok
func (s *AuthTestSuite) TestCE_AccountLockLogin() {
	ctx := context.TODO()

	for i := 0; i <= 6; i++ {
		_, _, _, varror := s.service.Login(ctx, s.email, "invalidPassword")
		require.NotNil(s.T(), varror.Error)
	}
}

func (s *AuthTestSuite) TestDA_ChangePassword() {
	ctx := context.TODO()

	newPassword := "password-must-be-changed"

	varror := s.service.ChangePassword(ctx, s.userID, s.password, newPassword)
	require.Nil(s.T(), varror)

	s.quickLogin(s.email, newPassword)
	s.password = newPassword
}

// ok
func (s *AuthTestSuite) TestDB_InvalidInputsChangePassword() {
	ctx := context.TODO()

	varror := s.service.ChangePassword(ctx, "", "", "")
	require.NotNil(s.T(), varror)
}

// ok
func (s *AuthTestSuite) TestDC_InvalidUserIDChangePassword() {
	ctx := context.TODO()

	varror := s.service.ChangePassword(ctx, "invalid", s.password, "NewFuckingPassword")
	require.Equal(s.T(), varror.Error, service.ErrNotFound)
}

// ok
func (s *AuthTestSuite) TestDD_InvalidOldPasswordChangePassword() {
	ctx := context.TODO()

	varror := s.service.ChangePassword(ctx, s.userID, "invalid", "NewFuckingPassword")
	require.Equal(s.T(), varror.Error, service.ErrInvalidEmailOrPassword)
}

func (s *AuthTestSuite) TestEA_Authenticate() {
	ctx := context.TODO()

	user, varror := s.service.Authenticate(ctx, s.accessToken)
	require.Nil(s.T(), varror)

	require.Equal(s.T(), user.Email, s.email)
}

// ok
func (s *AuthTestSuite) TestEB_InvalidTokenAuthentication() {
	ctx := context.TODO()

	_, varror := s.service.Authenticate(ctx, "invalid-token")
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)
}

// ok
func (s *AuthTestSuite) TestEC_InvalidTokenAuthentication() {
	ctx := context.TODO()

	at, err := s.authManager.GenerateAccessToken(ctx, "", service.AccessTokenExpr)
	require.Nil(s.T(), err)

	_, varror := s.service.Authenticate(ctx, at)
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)
}

// ok
func (s *AuthTestSuite) TestED_InvalidUUIDAuthentication() {
	ctx := context.TODO()

	at, err := s.authManager.GenerateAccessToken(ctx, "invalid", service.AccessTokenExpr)
	require.Nil(s.T(), err)

	_, varror := s.service.Authenticate(ctx, at)
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)
}

// ok
func (s *AuthTestSuite) TestEE_InvalidInputAuthentication() {
	ctx := context.TODO()

	_, varror := s.service.Authenticate(ctx, "")
	require.NotNil(s.T(), varror)
}

func (s *AuthTestSuite) TestFA_RefreshToken() {
	ctx := context.TODO()

	accessToken, varror := s.service.RefreshToken(ctx, s.userID, s.refreshToken)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), accessToken)
	require.NotEqual(s.T(), accessToken, s.accessToken)

	s.accessToken = accessToken
}

// New
func (s *AuthTestSuite) TestFB_InvalidInputsRefreshToken() {
	ctx := context.TODO()

	_, varror := s.service.RefreshToken(ctx, "", "")
	require.NotNil(s.T(), varror)

}

// New
func (s *AuthTestSuite) TestFC_InvalidTokenRefreshToken() {
	ctx := context.TODO()

	_, varror := s.service.RefreshToken(ctx, s.userID, "Invalid")
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)

}

// New
func (s *AuthTestSuite) TestFD_InvalidUserIDRefreshToken() {
	ctx := context.TODO()

	_, varror := s.service.RefreshToken(ctx, "invalid", s.accessToken)
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)
}
func (s *AuthTestSuite) TestGA_SendResetPassword() {
	ctx := context.TODO()

	resetPasswordToken, timeout, varror := s.service.SendResetPassword(ctx, s.email, s.resetPasswordRedirectUrl)
	require.Nil(s.T(), varror)

	require.NotEmpty(s.T(), timeout)
	require.NotEmpty(s.T(), resetPasswordToken)

	require.Nil(s.T(), varror)

	s.resetPasswordToken = resetPasswordToken
}

// New
func (s *AuthTestSuite) TestGB_InvalidInputsSendResetPassword() {
	ctx := context.TODO()

	_, _, varror := s.service.SendResetPassword(ctx, "", "")
	require.NotNil(s.T(), varror)
}

//New

func (s *AuthTestSuite) TestGC_InvalidEmailSendResetPassword() {
	ctx := context.TODO()

	_, _, varror := s.service.SendResetPassword(ctx, "invalidEmail@gmail.com", s.resetPasswordRedirectUrl)
	require.Equal(s.T(), varror.Error, service.ErrNotFound)
}

func (s *AuthTestSuite) TestHA_SubmitResetPassword() {
	ctx := context.TODO()

	newPassword := "password-reset-must-work"

	varror := s.service.SubmitResetPassword(ctx, s.resetPasswordToken, newPassword)
	require.Nil(s.T(), varror)

	s.quickLogin(s.email, newPassword)
	s.password = newPassword
}

// New
func (s *AuthTestSuite) TestHB_InvalidInputsSubmitResetPassword() {
	ctx := context.TODO()

	varror := s.service.SubmitResetPassword(ctx, "", "")
	require.NotNil(s.T(), varror)
}

// New
func (s *AuthTestSuite) TestHC_InvalidTokenSubmitResetPassword() {
	ctx := context.TODO()

	varror := s.service.SubmitResetPassword(ctx, "invalid", "new-password")
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)
}

// New
func (s *AuthTestSuite) TestHD_InvalidTokenUUIDSubmitResetPassword() {
	ctx := context.TODO()
	rpt, err := s.authManager.GenerateToken(ctx, auth_manager.ResetPassword,
		&auth_manager.TokenPayload{
			UUID:      "invalid",
			TokenType: auth_manager.ResetPassword,
			CreatedAt: time.Now(),
		},
		service.ResetPasswordTokenExpr,
	)
	require.Nil(s.T(), err)

	varror := s.service.SubmitResetPassword(ctx, rpt, "new-password")
	require.Equal(s.T(), varror.Error, service.ErrAccessDenied)
}

// New
func (s *AuthTestSuite) TestIA_InvalidPasswordDeleteAccount() {
	ctx := context.TODO()

	varror := s.service.DeleteAccount(ctx, s.userID, "invalid")
	require.Equal(s.T(), varror.Error, service.ErrInvalidEmailOrPassword)
}

// New
func (s *AuthTestSuite) TestIB_InvalidUserIDDeleteAccount() {
	ctx := context.TODO()

	varror := s.service.DeleteAccount(ctx, "invalid", s.password)
	require.Equal(s.T(), varror.Error, service.ErrNotFound)
}

// FIXME
// func (s *AuthTestSuite) TestI_InvalidInputsDeleteAccount() {
// 	ctx := context.TODO()

//		varror := s.service.DeleteAccount(ctx, "", "")
//		require.NotNil(s.T(), varror)
//	}
//
// New
func (s *AuthTestSuite) TestIC_ShouldDeleteAccount() {
	ctx := context.TODO()

	varror := s.service.DeleteAccount(ctx, s.userID, s.password)
	require.Nil(s.T(), varror)
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
