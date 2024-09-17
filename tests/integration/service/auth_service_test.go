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
	service *service.AuthService

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

	// auth manger
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

func (s *AuthTestSuite) TestA_Register() {
	ctx := context.TODO()

	testCases := []struct {
		Name                   string
		LastName               string
		Username               string
		Email                  string
		Password               string
		VerifyEmailRedirectUrl string
		Valid                  bool
		Error                  error
	}{
		{
			Name:                   "",
			LastName:               "",
			Username:               "",
			Email:                  "",
			Password:               "",
			VerifyEmailRedirectUrl: "",
			Valid:                  false,
		},
		{
			Name:                   " ",
			LastName:               " ",
			Username:               "fk",
			Email:                  "email",
			Password:               "5285",
			VerifyEmailRedirectUrl: "example.com",
			Valid:                  false,
		},
		{
			Name:     "User1:Name",
			LastName: "User1:LastName",
			Email:    "user1@kavka.org",
			Username: "user1",
			Password: "plain-password",
			Valid:    true,
		},
		{
			Name:                   "User2:Name",
			LastName:               "User2:LastName",
			Email:                  "user1@kavka.org",
			Username:               "user2",
			Password:               "plain-password",
			VerifyEmailRedirectUrl: "example.com",
			Valid:                  false,
			Error:                  service.ErrEmailAlreadyExist,
		},
		{
			Name:                   "User2:Name",
			LastName:               "User2:LastName",
			Email:                  "user2@kavka.org",
			Username:               "user1",
			Password:               "plain-password",
			VerifyEmailRedirectUrl: "example.com",
			Valid:                  false,
			Error:                  service.ErrUsernameAlreadyExist,
		},
	}

	for _, tc := range testCases {
		verifyEmailToken, varror := s.service.Register(ctx, tc.Name, tc.LastName, tc.Username, tc.Email, tc.Password, tc.VerifyEmailRedirectUrl)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}
			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)
			s.email = tc.Email
			s.password = tc.Password
			s.verifyEmailToken = verifyEmailToken

			_, _, _, varror := s.service.Login(ctx, s.email, s.password)
			require.Equal(s.T(), varror.Error, service.ErrEmailNotVerified)

			_, _, varror = s.service.SendResetPassword(ctx, s.email, s.resetPasswordRedirectUrl)
			require.Equal(s.T(), varror.Error, service.ErrEmailNotVerified)
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestB_VerifyEmail() {
	ctx := context.TODO()

	testCases := []struct {
		VerifyEmailToken string
		Valid            bool
		Error            error
	}{
		{
			VerifyEmailToken: "",
			Valid:            false,
		},
		{
			VerifyEmailToken: "invalidToken",
			Error:            service.ErrAccessDenied,
			Valid:            false,
		},
		{
			VerifyEmailToken: s.verifyEmailToken,
			Valid:            true,
		},
	}

	for _, tc := range testCases {
		varror := s.service.VerifyEmail(ctx, tc.VerifyEmailToken)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}
			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestC_Login() {
	ctx := context.TODO()

	testCases := []struct {
		Email    string
		Password string
		Valid    bool
		Error    error
	}{
		{
			Email:    "",
			Password: "",
			Valid:    false,
		},
		{
			Email:    "email",
			Password: "5285",
			Valid:    false,
		},
		{
			Email:    "invalid@kavka.org",
			Password: s.password,
			Error:    service.ErrInvalidEmailOrPassword,
			Valid:    false,
		},
		{
			Email:    s.email,
			Password: "invalidpassword",
			Valid:    false,
			Error:    service.ErrInvalidEmailOrPassword,
		},
		{
			Email:    "invalidpp@kafka.com",
			Password: "invalid-password",
			Valid:    false,
		},
		{
			Email:    s.email,
			Password: s.password,
			Valid:    true,
		},
	}

	for _, tc := range testCases {
		user, accessToken, refreshToken, varror := s.service.Login(ctx, tc.Email, tc.Password)

		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}

			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)

			require.NotEmpty(s.T(), accessToken)
			require.NotEmpty(s.T(), refreshToken)
			require.NotEmpty(s.T(), user)
			require.Equal(s.T(), user.Email, s.email)

			s.accessToken = accessToken
			s.refreshToken = refreshToken
			s.userID = user.UserID

			s.user = user
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestD_ChangePassword() {
	ctx := context.TODO()

	testCases := []struct {
		UserID      string
		OldPassword string
		NewPassword string
		Valid       bool
		Error       error
	}{
		{
			UserID:      "",
			OldPassword: "",
			NewPassword: "",
			Valid:       false,
		},
		{
			UserID:      "1857",
			OldPassword: "48796552",
			NewPassword: "7485",
			Valid:       false,
		},
		{
			UserID:      "invalid",
			OldPassword: s.password,
			NewPassword: "7485896554",
			Error:       service.ErrNotFound,
			Valid:       false,
		},
		{
			UserID:      s.userID,
			OldPassword: "48796552",
			NewPassword: "okaiojdkOKS17",
			Error:       service.ErrInvalidEmailOrPassword,
			Valid:       false,
		},
		{
			UserID:      s.userID,
			OldPassword: s.password,
			NewPassword: "validpassword85",
			Valid:       true,
		},
	}

	for _, tc := range testCases {
		varror := s.service.ChangePassword(ctx, tc.UserID, tc.OldPassword, tc.NewPassword)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}

			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)

			s.quickLogin(s.email, tc.NewPassword)
			s.password = tc.NewPassword
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestE_Authenticate() {
	ctx := context.TODO()

	emptyPayloadAT, err := s.authManager.GenerateAccessToken(ctx, "", service.AccessTokenExpr)
	require.Nil(s.T(), err)

	invalidPayloadAT, err := s.authManager.GenerateAccessToken(ctx, "invalid", service.AccessTokenExpr)
	require.Nil(s.T(), err)

	testCases := []struct {
		AccessToken string
		Valid       bool
		Error       error
	}{
		{
			AccessToken: "",
			Valid:       false,
		},
		{
			AccessToken: "invalid-Access-Token",
			Error:       service.ErrAccessDenied,
			Valid:       false,
		},
		{
			AccessToken: emptyPayloadAT,
			Error:       service.ErrAccessDenied,
			Valid:       false,
		},
		{
			AccessToken: invalidPayloadAT,
			Error:       service.ErrAccessDenied,
			Valid:       false,
		},
		{
			AccessToken: s.accessToken,
			Valid:       true,
		},
	}

	for _, tc := range testCases {
		user, varror := s.service.Authenticate(ctx, tc.AccessToken)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}

			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)

			require.Equal(s.T(), user.Email, s.email)
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestF_RefreshToken() {
	ctx := context.TODO()

	testCases := []struct {
		UserID       string
		RefreshToken string
		Valid        bool
		Error        error
	}{
		{
			UserID:       "",
			RefreshToken: "",
			Valid:        false,
		},
		{
			UserID:       "",
			RefreshToken: "invalid",
			Valid:        false,
		},
		{
			UserID:       "12345",
			RefreshToken: "",
			Valid:        false,
		},
		{
			UserID:       "invalid",
			RefreshToken: s.refreshToken,
			Error:        service.ErrAccessDenied,
			Valid:        false,
		},
		{
			UserID:       s.userID,
			RefreshToken: "invalid",
			Error:        service.ErrAccessDenied,
			Valid:        false,
		},
		{
			UserID:       s.userID,
			RefreshToken: s.refreshToken,
			Valid:        true,
		},
	}

	for _, tc := range testCases {
		accessToken, varror := s.service.RefreshToken(ctx, tc.UserID, tc.RefreshToken)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}

			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)

			require.NotEmpty(s.T(), accessToken)
			require.NotEqual(s.T(), accessToken, s.accessToken)

			s.accessToken = accessToken
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestG_SendResetPassword() {
	ctx := context.TODO()

	testCases := []struct {
		Email            string
		ResetPasswordURL string
		Valid            bool
		Error            error
	}{
		{
			Email:            "",
			ResetPasswordURL: "",
			Valid:            false,
		},
		{
			Email:            "invalid",
			ResetPasswordURL: "invalid",
			Valid:            false,
		},
		{
			Email:            "invalid@gmail.com",
			ResetPasswordURL: "invalid",
			Error:            service.ErrNotFound,
			Valid:            false,
		},
		{
			Email:            s.email,
			ResetPasswordURL: s.resetPasswordRedirectUrl,
			Valid:            true,
		},
	}

	for _, tc := range testCases {
		resetPasswordToken, timeout, varror := s.service.SendResetPassword(ctx, tc.Email, tc.ResetPasswordURL)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}

			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)

			require.NotEmpty(s.T(), timeout)
			require.NotEmpty(s.T(), resetPasswordToken)

			require.Nil(s.T(), varror)

			s.resetPasswordToken = resetPasswordToken
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestH_SubmitResetPassword() {
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

	testCases := []struct {
		NewPassword              string
		SubmitResetPasswordToken string
		Valid                    bool
		Error                    error
	}{
		{
			NewPassword:              "",
			SubmitResetPasswordToken: "",
			Valid:                    false,
		},
		{
			NewPassword:              "5689",
			SubmitResetPasswordToken: s.resetPasswordToken,
			Valid:                    false,
		},
		{
			NewPassword:              "valid-password8568",
			SubmitResetPasswordToken: "invalid",
			Error:                    service.ErrAccessDenied,
			Valid:                    false,
		},
		{
			NewPassword:              "valid-password8568",
			SubmitResetPasswordToken: rpt,
			Error:                    service.ErrAccessDenied,
			Valid:                    false,
		},
		{
			NewPassword:              "valid-password-anfaj",
			SubmitResetPasswordToken: s.resetPasswordToken,
			Valid:                    true,
		},
	}

	for _, tc := range testCases {
		varror := s.service.SubmitResetPassword(ctx, tc.SubmitResetPasswordToken, tc.NewPassword)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}

			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)

			s.quickLogin(s.email, tc.NewPassword)

			s.password = tc.NewPassword
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func (s *AuthTestSuite) TestI_InvalidPasswordDeleteAccount() {
	ctx := context.TODO()

	testCases := []struct {
		UserID   string
		Password string
		Error    error
		Valid    bool
	}{
		{
			UserID:   "invalid",
			Password: s.password,
			Error:    service.ErrNotFound,
			Valid:    false,
		},
		{
			UserID:   s.userID,
			Password: "invalid",
			Error:    service.ErrInvalidEmailOrPassword,
			Valid:    false,
		},
		{
			UserID:   s.userID,
			Password: s.password,
			Valid:    true,
		},
	}

	for _, tc := range testCases {
		varror := s.service.DeleteAccount(ctx, tc.UserID, tc.Password)
		if !tc.Valid {
			if tc.Error != nil {
				require.Equal(s.T(), tc.Error, varror.Error)
				continue
			}
			require.NotNil(s.T(), varror)
		} else if tc.Valid {
			require.Nil(s.T(), varror)
		} else {
			require.Fail(s.T(), "not specific")
		}
	}
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
