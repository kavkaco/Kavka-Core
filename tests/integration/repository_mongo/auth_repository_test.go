package repository_mongo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kavkaco/Kavka-Core/internal/model"
	authRepository "github.com/kavkaco/Kavka-Core/internal/repository/auth"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	repo   authRepository.AuthRepository
	UserID model.UserID
}

func (s *AuthTestSuite) SetupTest() {
	// Set a user id to be used in tests
	s.UserID = uuid.NewString()

	s.repo = authRepository.NewRepository(db)
}

func (s *AuthTestSuite) TestCreate() {
	ctx := context.TODO()

	passwordHash := "hashed-password"

	auth, err := s.repo.Create(ctx, s.UserID, passwordHash)
	require.NoError(s.T(), err)
	require.Equal(s.T(), auth.UserID, s.UserID)
	require.Equal(s.T(), auth.PasswordHash, passwordHash)
}

func (s *AuthTestSuite) TestGetUserAuth() {
	ctx := context.TODO()

	auth, err := s.repo.Create(ctx, s.UserID, "hashed-password")
	require.NoError(s.T(), err)

	tests := []struct {
		name      string
		userID    string
		wantError error
	}{
		{
			name:   "success",
			userID: auth.UserID,
		},
		{
			name:      "should return error",
			userID:    "wrong-user-id",
			wantError: authRepository.ErrAuthNotFound,
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			foundAuth, err := s.repo.GetUserAuth(ctx, tc.userID)
			if err != tc.wantError {
				t.Errorf("Unexpected error: %v", err)
			}

			if tc.wantError == nil {
				require.Equal(t, foundAuth.UserID, tc.userID)
			}
		})
	}
}

func (s *AuthTestSuite) TestChangePassword() {
	ctx := context.TODO()

	hashedPassword := "hashed-password"
	newPassword := "new-password"

	_, err := s.repo.Create(ctx, s.UserID, hashedPassword)
	require.NoError(s.T(), err)

	ok, err := s.repo.ChangePassword(ctx, s.UserID, newPassword)
	require.NoError(s.T(), err)
	require.True(s.T(), ok)
}

func TestAuthSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(AuthTestSuite))
}
