package interceptor

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/go-redis/redis/v8"
	"google.golang.org/genproto/googleapis/rpc/code"
)

type RateLimiter struct {
	client      *redis.Client
	limit       int
	window      time.Duration
	prefix      string
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration, prefix string) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
		prefix: prefix,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, identifier string) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s:%s", rl.prefix, identifier)

	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rl.client.Expire(ctx, key, rl.window)
	}

	return count <= int64(rl.limit), nil
}

func (rl *RateLimiter) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		allowed, err := rl.Allow(ctx, req.Peer().Addr)
		if err != nil || !allowed {
			return nil, connect.NewError(connect.Code(code.Code_RESOURCE_EXHAUSTED), fmt.Errorf("rate limit exceeded"))
		}

		return next(ctx, req)
	}
}

func (rl *RateLimiter) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (rl *RateLimiter) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		allowed, err := rl.Allow(ctx, shc.Peer().Addr)
		if err != nil || !allowed {
			return connect.NewError(connect.Code(code.Code_RESOURCE_EXHAUSTED), fmt.Errorf("rate limit exceeded"))
		}

		return next(ctx, shc)
	})
}
