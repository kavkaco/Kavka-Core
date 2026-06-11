package proto_model_transformer

import (
	"github.com/kavkaco/Kavka-Core/internal/model"
	modelv1 "github.com/kavkaco/Kavka-ProtoBuf/gen/go/protobuf/model/user/v1"
)

func UserToProto(user model.User) *modelv1.User {
	return &modelv1.User{
		UserId:    user.UserID,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Biography: user.Biography,
	}
}

func UsersToProto(users []model.User) []*modelv1.User {
	transformedUsers := make([]*modelv1.User, 0, len(users))

	for _, v := range users {
		transformedUsers = append(transformedUsers, UserToProto(v))
	}

	return transformedUsers
}
