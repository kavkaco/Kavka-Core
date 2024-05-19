package auth_manager_test

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	token_manager "github.com/kavkaco/Kavka-Core/pkg/token_manager"
// 	token_manager_mock "github.com/kavkaco/Kavka-Core/pkg/token_manager/mocks"

// 	"go.uber.org/mock/gomock"
// )

// func TestGenerateToken(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockTokenManager := token_manager_mock.NewMockTokenManager(ctrl)

// 	createdAt := time.Now().UTC()
// 	expr := time.Minute * 2
// 	claims := token_manager.TokenClaims{
// 		UserStaticID: "user-static-id",
// 		TokenType:    token_manager.AccessToken,
// 		CreatedAt:    createdAt,
// 	}
// 	mockTokenManager.EXPECT().GenerateToken(context.Background(), token_manager.AccessToken, &claims, expr).AnyTimes()
// }
