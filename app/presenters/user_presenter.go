package presenters

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/internal/model/chat"
	"github.com/kavkaco/Kavka-Core/internal/model/user"
	"github.com/kavkaco/Kavka-Core/pkg/session"
)

// the data and information that would be send to user after authentication.
// the messages of chats must not be sent

type UserInfoDto struct {
	Message   string       `json:"message"`
	Code      int          `json:"code"`
	UserInfo  *user.User   `json:"user"`
	UserChats []chat.ChatC `json:"chats"`
}

func AccessDenied(ctx *gin.Context) {
	code := http.StatusForbidden

	ctx.JSON(code, SimpleMessageDto{
		Code:    code,
		Message: "Forbidden",
	})
}

func SendTokensHeader(ctx *gin.Context, tokens session.LoginTokens) {
	ctx.Header("refresh", tokens.RefreshToken)
	ctx.Header("authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
}

func ResponseUserInfo(ctx *gin.Context, userInfo *user.User, userChats []chat.ChatC) error {
	// Marshal all of the chats into json using by ChatAsJson function
	marshaledChatsJson := []chat.ChatC{}

	for _, v := range userChats {
		chatJson, err := ChatAsJSON(v, userInfo.StaticID)
		if err != nil {
			return err
		}

		marshaledChatsJson = append(marshaledChatsJson, chatJson.(chat.ChatC))
	}

	// Response JSON

	ctx.JSON(http.StatusOK, UserInfoDto{
		Message:   "Success",
		Code:      200,
		UserInfo:  userInfo,
		UserChats: marshaledChatsJson,
	})

	return nil
}
