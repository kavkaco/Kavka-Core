package interceptor

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/kavkaco/Kavka-Core/internal/service/auth"
)

const accessTokenHeader = "X-Access-Token"

type UserID struct{}

var ErrEmptyUserID = errors.New("empty user id after passing auth interceptor")

type authInterceptor struct {
	authService *auth.AuthService
}

func (i *authInterceptor) WrapStreamingClient(connect.StreamingClientFunc) connect.StreamingClientFunc {
	panic("unimplemented")
}

func (i *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		accessToken := shc.RequestHeader().Get(accessTokenHeader)

		ctx, err := i.authRequired(ctx, accessToken)
		if err != nil {
			return err
		}

		return next(ctx, shc)
	})
}

func (i *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		accessToken := req.Header().Get(accessTokenHeader)

		ctx, err := i.authRequired(ctx, accessToken)
		if err != nil {
			return nil, err
		}

		return next(ctx, req)
	}
}

func (i *authInterceptor) authRequired(ctx context.Context, accessToken string) (context.Context, error) {
	if len(accessToken) == 0 {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no token provided"))
	}

	user, err := i.authService.Authenticate(ctx, accessToken)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	ctx = context.WithValue(ctx, UserID{}, user.UserID)

	return ctx, nil
}

func NewAuthInterceptor(authService *auth.AuthService) connect.Interceptor {
	return &authInterceptor{authService}
}
