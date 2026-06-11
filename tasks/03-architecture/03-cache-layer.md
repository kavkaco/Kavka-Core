# Add Caching Layer

**Files:** `infra/cache/` (currently empty)
**Priority:** P2 (Medium)
**Impact:** MongoDB hit on every request for static data

## Current State

`infra/cache/` contains only a `.gitkeep` file. Redis is available but only used for token storage via `go-auth-manager`.

## Cache Strategy

### What to Cache

| Data | TTL | Invalidation |
|---|---|---|
| User profiles (by user_id) | 5 min | On profile update |
| User profiles (by username) | 5 min | On profile update |
| User profiles (by email) | 5 min | On email update |
| Chat list for user | 1 min | On chat create/join/leave |
| Direct chat lookup | 10 min | On direct chat creation |

### What NOT to Cache

- Messages (too volatile, too much data)
- Auth tokens (already in Redis via `go-auth-manager`)
- Search results (query-dependent)

## Implementation

### Cache Interface

```go
package cache

type Cache interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}
```

### Redis Implementation

```go
type redisCache struct {
    client *redis.Client
}

func (c *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}
```

### Service Integration

```go
type UserService struct {
    repo  repository.UserRepository
    cache cache.Cache
}

func (s *UserService) FindByUserID(ctx, userID) (*model.User, error) {
    // Try cache first
    var user model.User
    err := s.cache.Get(ctx, "user:"+userID, &user)
    if err == nil {
        return &user, nil
    }
    
    // Fall back to DB
    user, err = s.repo.FindByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Populate cache
    s.cache.Set(ctx, "user:"+userID, user, 5*time.Minute)
    return user, nil
}
```

## Implementation Priority

1. Define `Cache` interface in `infra/cache/`
2. Implement Redis-backed cache
3. Integrate with `UserService` for profile lookups
4. Integrate with `ChatService` for chat lookups
5. Add cache invalidation on writes
