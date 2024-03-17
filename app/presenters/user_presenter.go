package presenters

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kavkaco/Kavka-Core/internal/domain/chat"
	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"github.com/kavkaco/Kavka-Core/pkg/session"
)

// the data and information that would be send to user after authentication;
// the messages of chats must not be sent!

type UserInfoDto struct {
	Message   string      `json:"message"`
	Code      int         `json:"code"`
	UserInfo  *user.User  `json:"user"`
	UserChats []chat.Chat `json:"chats"`
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

func ResponseUserInfo(ctx *gin.Context, userInfo *user.User, userChats []chat.Chat) {
	ctx.JSON(http.StatusOK, UserInfoDto{
		Message:   "Success",
		Code:      200,
		UserInfo:  userInfo,
		UserChats: userChats,
	})
}
