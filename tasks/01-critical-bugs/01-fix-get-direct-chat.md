# Fix GetDirectChat Bug

**File:** `database/repo_mongo/chat_repository.go:232-248`
**Priority:** P0 (Critical)
**Impact:** `GetDirectChat` never returns results for direct chats

## Problem

The query at line 244 uses:

```go
"chat_detail.chat_type": bson.M{"$ne": "direct"},
```

Two issues:
1. **Wrong field path**: `chat_detail.chat_type` should be `chat_type` (the field is at document root, not inside `chat_detail`)
2. **Negation logic**: `$ne: "direct"` **excludes** direct chats instead of matching them. The intent was to filter **only** direct chats, so it should match `$eq: "direct"`

## Fix

Change line 244 from:
```go
"chat_detail.chat_type": bson.M{"$ne": "direct"},
```
to:
```go
"chat_type": "direct",
```

## Validation

- `GetDirectChat` returns the direct chat document for a pair of users
- All tests in `tests/integration/repository/` pass
- Manual testing via gRPC call returns direct chat correctly
