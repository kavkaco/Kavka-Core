package service

import (
	"context"
	"testing"
	"time"

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

	authToken, verifyEmailToken string
}

func (s *AuthTestSuite) SetupSuite() {
	s.lem = lorem.New()

	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	authManager := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		RefreshTokenExpiration: time.Second * 10,
		AccessTokenExpiration:  time.Second * 3,
		PrivateKey:             "private-key",
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
}

func (s *AuthTestSuite) TestB_VerifyEmail() {
	ctx := context.TODO()

	err := s.service.VerifyEmail(ctx, s.verifyEmailToken)
	require.NoError(s.T(), err)
}

// ANCHOR - write test for login method

// func (s *AuthTestSuite) TestC_Authenticate() {
// 	ctx := context.TODO()

// 	name := s.lem.FirstName(0)
// 	lastName := s.lem.LastName()
// 	username := s.lem.Word(1, 10)
// 	email := s.lem.Email()
// 	password := "strong-password"

// 	user, verifyEmailToken, err := s.service.Register(ctx, name, lastName, username, email, password)
// 	require.NoError(s.T(), err)

// 	require.NotEmpty(s.T(), verifyEmailToken)
// 	require.Equal(s.T(), user.Name, name)
// 	require.Equal(s.T(), user.LastName, lastName)
// 	require.Equal(s.T(), user.Username, username)
// 	require.Equal(s.T(), user.Email, email)
// }

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
