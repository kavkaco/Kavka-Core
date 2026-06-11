# Implement Missing Repository Methods

**Priority:** P1 (High)
**Impact:** Broken functionality, panics

## Empty Implementations

### 1. GetChatMembers

**File:** `database/repo_mongo/chat_repository.go:56-58`

Currently returns `[]model.Member{}` (empty). The interface also misses `ctx` parameter.

**Fix:**
1. Add `ctx context.Context` to the interface: `GetChatMembers(ctx context.Context, chatID model.ChatID) ([]model.Member, error)`
2. Implement with MongoDB aggregation:

```go
func (repo *chatRepository) GetChatMembers(ctx context.Context, chatID model.ChatID) ([]model.Member, error) {
    pipeline := bson.A{
        {"$match": bson.M{"_id": chatID}},
        {"$project": bson.M{
            "members": "$chat_detail.members",
        }},
        {"$lookup": bson.M{
            "from":         "users",
            "localField":   "members",
            "foreignField": "user_id",
            "as":           "member_users",
        }},
        {"$project": bson.M{
            "member_users.user_id":   1,
            "member_users.name":      1,
            "member_users.last_name": 1,
        }},
    }
    // ...
}
```

### 2. SearchInChat

**File:** `database/repo_mongo/search_repository.go:64-65`

Currently panics. Implement full-text search within a specific chat's messages.

### 3. UpdateTextMessage

**File:** `internal/service/message/message_service.go:141-143`

Currently panics. Implement with repository call.

## Validation

- `GetChatMembers` returns all members with names
- `SearchInChat` returns matching messages
- `UpdateTextMessage` updates message content and sets `edited = true`
- All integration tests pass
