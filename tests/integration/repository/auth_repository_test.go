package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/utils/hash"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	repo          repository.AuthRepository
	hashManager   *hash.HashManager
	userID        model.UserID
	plainPassword string
}

func (s *AuthTestSuite) SetupSuite() {
	// Set a user id to be used in tests
	s.userID = fmt.Sprintf("%d", random.GenerateUserID())

	s.hashManager = hash.NewHashManager(hash.DefaultHashParams)
	s.repo = repository.NewAuthRepository(db)

	// Set plain password to generate hash in auth creation
	s.plainPassword = "kavkaco"
}

func (s *AuthTestSuite) TestA_Create() {
	ctx := context.TODO()

	passwordHash, err := s.hashManager.HashPassword(s.plainPassword)
	require.NoError(s.T(), err)

	authModel := model.NewAuth(s.userID, passwordHash)

	auth, err := s.repo.Create(ctx, authModel)
	require.NoError(s.T(), err)
	require.Equal(s.T(), auth.UserID, s.userID)
	require.Equal(s.T(), auth.PasswordHash, passwordHash)
	require.False(s.T(), auth.EmailVerified)
}

func (s *AuthTestSuite) TestB_GetUserAuth() {
	ctx := context.TODO()

	userAuth, err := s.repo.GetUserAuth(ctx, s.userID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), userAuth.UserID, s.userID)
}

func (s *AuthTestSuite) TestC_IncrementFailedLoginAttempts() {
	ctx := context.TODO()

	err := s.repo.IncrementFailedLoginAttempts(ctx, s.userID)
	require.NoError(s.T(), err)

	// Fetch user to check the value of failed login attempts
	userAuth, err := s.repo.GetUserAuth(ctx, s.userID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), userAuth.FailedLoginAttempts, 1)
}

func (s *AuthTestSuite) TestD_VerifyEmail() {
	ctx := context.TODO()

	err := s.repo.VerifyEmail(ctx, s.userID)
	require.NoError(s.T(), err)

	// Fetch user to see email verified or not
	userAuth, err := s.repo.GetUserAuth(ctx, s.userID)
	require.NoError(s.T(), err)
	require.True(s.T(), userAuth.EmailVerified)
}

func (s *AuthTestSuite) TestE_ChangePassword() {
	ctx := context.TODO()

	newPlainPassword := "kavkaco-new" // nolint
	newPasswordHash, err := s.hashManager.HashPassword(newPlainPassword)
	require.NoError(s.T(), err)

	err = s.repo.ChangePassword(ctx, s.userID, newPasswordHash)
	require.NoError(s.T(), err)
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
