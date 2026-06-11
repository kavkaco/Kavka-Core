# Rate Limiting

**Priority:** P2 (Medium)
**Impact:** Prevent abuse, ensure fair resource usage

## Implementation

Use existing Redis instance (already used for token storage) for distributed rate limiting.

### Per-User Token Bucket

```go
type RateLimiter struct {
    client *redis.Client
}

func (rl *RateLimiter) Allow(ctx context.Context, userID string, limit int, window time.Duration) (bool, error) {
    key := "ratelimit:" + userID
    count, err := rl.client.Incr(ctx, key).Result()
    if err != nil {
        return false, err
    }
    if count == 1 {
        rl.client.Expire(ctx, key, window)
    }
    return count <= int64(limit), nil
}
```

### Rate Limits

| Endpoint | Limit | Window |
|---|---|---|
| `/auth.v1.AuthService/Login` | 10 | 1 minute |
| `/auth.v1.AuthService/Register` | 3 | 1 hour |
| `/message.v1.MessageService/SendTextMessage` | 60 | 1 minute |
| `/chat.v1.ChatService/CreateGroup` | 10 | 1 hour |
| All other endpoints | 100 | 1 minute |

### Middleware Integration

Add rate limiter as a second interceptor in the gRPC chain:

```go
rateLimiter := NewRateLimiter(redisClient)

router.Handle(authGrpcRoute, authGrpcRouter) // no rate limit for auth
router.Handle(chatGrpcRoute, rateLimiter.Wrap(chatGrpcRouter))
```
