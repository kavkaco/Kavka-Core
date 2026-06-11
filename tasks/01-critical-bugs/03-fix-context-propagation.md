# Fix Context Propagation

**Files:** Multiple
**Priority:** P1 (High)
**Impact:** Deadline/timeout propagation broken, potential goroutine leaks, tracing breaks

## Problem

`context.TODO()` is used in many places instead of propagating the parent context. This breaks:
- Request cancellation chains (if client disconnects, DB queries keep running)
- Distributed tracing (span context is lost)
- Deadline propagation (timeouts don't cascade)

## Locations

| File | Line(s) | Current | Should Be |
|---|---|---|---|
| `database/repo_mongo/auth_repository.go` | multiple | `context.TODO()` | propagate `ctx` param |
| `database/repo_mongo/user_repository.go` | 105, 118 | `context.TODO()` | use `ctx` param |
| `internal/service/chat/chat_service.go` | 135, 226, 231 | `context.TODO()` | use passed `ctx` |
| `database/mongo_adapter.go` | 38, 45 | `context.TODO()` | `context.Background()` (singleton init is fine) |

## Fix Strategy

1. **Repository layer**: All public methods already take `ctx`; internal methods like `FindOne` should use it instead of `context.TODO()`
2. **Service layer**: Replace `context.TODO()` in `CreateDirect`, `CreateGroup`, `CreateChannel` with the `ctx` parameter passed from the handler
3. **Event publishing goroutines**: These intentionally use `context.Background()` because they're fire-and-forget. This is correct — just ensure the context isn't tied to the request lifecycle.

## Validation

- All tests pass
- No compile errors
- Context cancellation properly cancels in-flight DB operations
