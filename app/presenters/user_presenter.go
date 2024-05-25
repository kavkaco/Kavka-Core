package presenters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/internal/model"
)

type UserDto struct {
	Message   string       `json:"message"`
	Code      int          `json:"code"`
	User      *model.User  `json:"user"`
	UserChats []model.Chat `json:"chats"`
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

	ctx.JSON(http.StatusOK, UserDto{
		Message:   "user found",
		Code:      200,
		User:      user,
		UserChats: marshaledChatsJson,
	})

	return nil
}
