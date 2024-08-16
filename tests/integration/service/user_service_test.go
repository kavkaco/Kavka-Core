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

	name := "User5:NameChanged"
	lastName := "User5:LastNameChanged"
	username := "user5_changed"
	biography := "User5:Biography changed"

	varror := s.service.UpdateProfile(ctx, s.userID, name, lastName, username, biography)
	require.Nil(s.T(), varror)
}

func (s *UserTestSuite) TestB_InvalidInputUpdateProfile() {
	ctx := context.TODO()

	varror := s.service.UpdateProfile(ctx, s.userID, "", "", "", "")
	require.NotNil(s.T(), varror)
}
func (s *UserTestSuite) TestC_InvalidUserIDUpdateProfile() {
	ctx := context.TODO()

	name := "User5:NameChanged"
	lastName := "User5:LastNameChanged"
	username := "user5_changed"
	biography := "User5:Biography changed"

	varror := s.service.UpdateProfile(ctx, "invalid", name, lastName, username, biography)
	require.Equal(s.T(), varror.Error, service.ErrNotFound)
}
func TestUserSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(UserTestSuite))
}
