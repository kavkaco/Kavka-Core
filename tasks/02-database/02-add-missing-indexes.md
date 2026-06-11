# Add Missing MongoDB Indexes

**File:** `database/mongo_adapter.go`
**Priority:** P1 (High)
**Impact:** Full collection scans on every query

## Missing Indexes

| Collection | Index | Reason |
|---|---|---|
| `messages` | `{ chat_id: 1, created_at: -1 }` | Primary query pattern for fetching messages |
| `messages` | `{ chat_id: 1, message_id: 1 }` | Used by `FetchMessage` and `UpdateMessageContent` |
| `user_auth` | `{ user_id: 1 }` | Used by `GetUserAuth`, `ChangePassword`, etc. |
| `users` | `{ user_id: 1 }` (unique) | Primary lookup key — currently no index! |
| `chats` | `{ chat_detail.members: 1 }` | Queries for user's chats (alternative to user's `chats_list_ids`) |

## Implementation

Add to `ConfigureCollections` in `database/mongo_adapter.go`:

```go
// Messages indexes
_, err = db.Collection(MessagesCollection).Indexes().CreateMany(ctx, []mongo.IndexModel{
    {
        Keys: bson.D{
            {Key: "chat_id", Value: 1},
            {Key: "created_at", Value: -1},
        },
    },
    {
        Keys: bson.D{
            {Key: "chat_id", Value: 1},
            {Key: "message_id", Value: 1},
        },
    },
})

// User auth indexes
_, err = db.Collection(AuthCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "user_id", Value: 1}},
})

// Users - user_id index (for FindByUserID lookups)
_, err = db.Collection(UsersCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "user_id", Value: 1}},
    Options: options.Index().SetUnique(true),
})
```

## Validation

- `explain()` shows index scans instead of collection scans for all queries
- Integration tests pass
- No performance regression
