package service

import (
	"context"
	"fmt"
	"testing"

	repository_mongo "github.com/kavkaco/Kavka-Core/database/repo_mongo"
	"github.com/kavkaco/Kavka-Core/internal/model"
	service "github.com/kavkaco/Kavka-Core/internal/service/user"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	service service.UserService

	userID model.UserID
}

func (s *UserTestSuite) SetupSuite() {
	ctx := context.TODO()

	userRepo := repository_mongo.NewUserMongoRepository(db)
	s.service = service.NewUserService(userRepo)

	user, err := userRepo.Create(ctx, &model.User{
		UserID:    fmt.Sprintf("%d", random.GenerateUserID()),
		Name:      "User5:Name",
		LastName:  "User5:LastName",
		Email:     "user5@kavka.org",
		Username:  "user5",
		Biography: "User5:biography",
	})
	require.NoError(s.T(), err)

	s.userID = user.UserID
}

func (s *UserTestSuite) TestA_UpdateProfile() {
	ctx := context.TODO()

	testCases := []struct {
		userID    string
		name      string
		lastName  string
		username  string
		biography string
		Valid     bool
		Error     error
	}{
		{
			userID:    "",
			name:      "",
			lastName:  "",
			username:  "",
			biography: "User5:Biography changed",
			Valid:     false,
		},
		{
			userID:    s.userID,
			name:      "l",
			lastName:  "Us",
			username:  "l",
			biography: "Ul",
			Valid:     false,
		},
		{
			userID:    "invalid",
			name:      "User5:NameChanged",
			lastName:  "User5:LastNameChanged",
			username:  "user5_changed",
			biography: "User5:Biography changed",
			Error:     service.ErrNotFound,
			Valid:     false,
		},
		{
			userID:    s.userID,
			name:      "User5:NameChanged",
			lastName:  "User5:LastNameChanged",
			username:  "user5_changed",
			biography: "User5:Biography changed",
			Valid:     true,
		},
	}

	for _, tc := range testCases {
		varror := s.service.UpdateProfile(ctx, tc.userID, tc.name, tc.lastName, tc.username, tc.biography)
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

func TestUserSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(UserTestSuite))
}
