# Fix Race Condition in StreamSubscriber

**File:** `infra/stream/stream_sub.go`
**Priority:** P1 (High)
**Impact:** Data race on `subscribedUsers` slice under concurrent access

## Problem

`subscribedUsers []StreamSubscribedUser` is accessed from:
- NATS subscription callback goroutine (reads)
- `UserSubscribe` method (writes)
- `UserUnsubscribe` method (writes)

None of these are synchronized with a mutex, causing a data race.

## Fix

Add `sync.RWMutex` to the `sub` struct:

```go
type sub struct {
    nc              *nats.Conn
    logger          *log.SubLogger
    mu              sync.RWMutex
    subscribedUsers []StreamSubscribedUser
}
```

Lock appropriately:
- `UserSubscribe`: `mu.Lock()`
- `UserUnsubscribe`: `mu.Lock()`
- NATS callback reads: `mu.RLock()`

## Validation

- Run with `-race` flag — no data races reported
- Concurrent subscribe/unsubscribe operations don't panic
- Events are still delivered to all subscribed users
