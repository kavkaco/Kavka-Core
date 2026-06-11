# Fix Proto Transformer Global Mutable State

**File:** `internal/model/proto_model_transformer/`
**Priority:** P1 (High)
**Impact:** Race condition and incorrect results under concurrent gRPC requests

## Problem

The transformer files use package-level variables like:
```go
var transformedChats []*chatv1.Chat
var transformedUsers []*userv1.User
```

Under concurrent requests, these shared slices get overwritten or read simultaneously, causing:
- Incorrect data returned to clients
- Panic (concurrent slice writes)
- Memory corruption

## Fix

1. Remove all package-level variables
2. Make every function allocate and return fresh slices:

```go
// Before:
var transformedChats []*chatv1.Chat
for _, c := range chats {
    transformedChats = append(transformedChats, ChatToProto(c))
}
return transformedChats

// After:
func ChatToProto(chat model.ChatDTO) (*chatv1.Chat, error) { ... }
func ChatsToProto(chats []model.ChatDTO) []*chatv1.Chat {
    result := make([]*chatv1.Chat, 0, len(chats))
    for _, c := range chats {
        result = append(result, ChatToProto(c))
    }
    return result
}
```

## Validation

- Run with `-race` flag — no data races reported
- Concurrent requests return correct, non-interleaved results
- All tests pass
