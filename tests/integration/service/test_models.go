package service

import (
	"fmt"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/utils/random"
)

var channelChatDetailTestModel = model.ChannelChatDetail{
	Title:        "Channel1",
	Username:     "channel1",
	Owner:        userTestModels[0].UserID,
	Members:      []model.UserID{userTestModels[0].UserID},
	Admins:       []model.UserID{userTestModels[0].UserID},
	Description:  "Channel1:Description",
	RemovedUsers: []model.UserID{},
}

var groupChatDetailTestModel = model.GroupChatDetail{
	Title:        "Group1",
	Username:     "Group1",
	Owner:        userTestModels[0].UserID,
	Members:      []model.UserID{userTestModels[0].UserID},
	Admins:       []model.UserID{userTestModels[0].UserID},
	Description:  "Group1:Description",
	RemovedUsers: []model.UserID{},
}

var directChatDetailTestModel = model.DirectChatDetail{
	Sides: [2]model.UserID{userTestModels[0].UserID, userTestModels[1].UserID},
}

var userTestModels = [2]model.User{
	{
		UserID:    fmt.Sprintf("%d", random.GenerateUserID()),
		Name:      "User1:Name",
		LastName:  "User1:LastName",
		Email:     "user1@kavka.org",
		Username:  "user1",
		Biography: "User1:biography",
	},
	{
		UserID:    fmt.Sprintf("%d", random.GenerateUserID()),
		Name:      "User2:Name",
		LastName:  "User2:LastName",
		Email:     "user2@kavka.org",
		Username:  "user2",
		Biography: "User2:biography",
	},
}
