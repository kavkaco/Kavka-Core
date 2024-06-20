package grpc_model

import (
	"github.com/kavkaco/Kavka-Core/internal/model"
	modelv1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/model/user/v1"
)

func TransformUserToGrpcModel(user *model.User) *modelv1.User {
	return &modelv1.User{
		UserId:    user.UserID,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Biography: user.Biography,
	}
}

func TransformGrpcModelToUser(user *modelv1.User) *model.User {
	return &model.User{
		UserID:    user.UserId,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Biography: user.Biography,
	}
}
