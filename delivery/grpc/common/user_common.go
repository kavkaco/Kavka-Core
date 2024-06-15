package grpc_common

import (
	"github.com/kavkaco/Kavka-Core/internal/model"
	commonv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/proto/common/v1"
)

func TransformUserToGrpcModel(user *model.User) *commonv1.User {
	return &commonv1.User{
		UserId:    user.UserID,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Biography: user.Biography,
	}
}

func TransformGrpcModelToUser(user *commonv1.User) *model.User {
	return &model.User{
		UserID:    user.UserId,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Biography: user.Biography,
	}
}
