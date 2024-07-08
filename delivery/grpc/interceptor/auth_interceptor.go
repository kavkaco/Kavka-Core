package interceptor

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
)

const accessTokenHeader = "X-Access-Token"

type UserIDKey struct{}

func NewAuthInterceptor(authService auth.AuthService) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			accessToken := req.Header().Get(accessTokenHeader)

			if len(accessToken) == 0 {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no token provided"))
			}

			user, err := authService.Authenticate(ctx, accessToken)
			if err != nil || user == nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, nil)
			}

			ctx = context.WithValue(ctx, UserIDKey{}, user.UserID)

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
