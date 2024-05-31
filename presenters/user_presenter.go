package presenters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/internal/model"
)

type UserDto struct {
	Message   string         `json:"message"`
	Code      int            `json:"code"`
	User      *UserDetailDto `json:"user"`
	UserChats []model.Chat   `json:"chats"`
}

type UserDetailDto struct {
	UserID    model.UserID `json:"userID"`
	Name      string       `json:"name"`
	LastName  string       `json:"lastName"`
	Email     string       `json:"email"`
	Username  string       `json:"username"`
	Biography string       `json:"biography"`
}

func AccessDenied(ctx *gin.Context) {
	code := http.StatusForbidden

	ctx.JSON(code, CodeMessageDto{
		Code:    code,
		Message: "forbidden",
	})
}

func UserResponse(ctx *gin.Context, user *model.User, userChats []model.Chat) error {
	// Marshal all of the chats into json using by ChatAsJson function
	marshaledChatsJson := []model.Chat{}

	for _, v := range userChats {
		chatJson, err := ChatAsJSON(v, user.UserID)
		if err != nil {
			return err
		}

		marshaledChatsJson = append(marshaledChatsJson, chatJson.(model.Chat))
	}

	// Response JSON
	var normalizedUser UserDetailDto

	normalizedUser.UserID = user.UserID
	normalizedUser.Email = user.Email
	normalizedUser.Username = user.Username
	normalizedUser.Name = user.Name
	normalizedUser.LastName = user.LastName
	normalizedUser.Biography = user.Biography

	ctx.JSON(http.StatusOK, UserDto{
		Message:   "user found",
		Code:      200,
		User:      &normalizedUser,
		UserChats: marshaledChatsJson,
	})

	return nil
}
