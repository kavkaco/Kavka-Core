package presenters

import (
	"Kavka/internal/domain/user"
	"Kavka/pkg/session"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserInfoDto struct {
	Message  string     `json:"message"`
	Code     int        `json:"code"`
	UserInfo *user.User `json:"user"`
}

func SendTokensHeader(ctx *gin.Context, tokens session.LoginTokens) {
	ctx.Header("refresh", tokens.RefreshToken)
	ctx.Header("authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
}

func ResponseUserInfo(ctx *gin.Context, userInfo *user.User) {
	ctx.JSON(http.StatusOK, UserInfoDto{
		Message:  "Success",
		Code:     200,
		UserInfo: userInfo,
	})
}
