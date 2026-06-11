# Bulk Operations

**Priority:** P3 (Low)
**Impact:** Reduce round-trips for batch operations

## Opportunities

### 1. Chat List Update

When creating a chat, `AddToUsersChatsList` is called separately for each user:

```go
// Current — 2 round-trips:
s.chatRepo.AddToUsersChatsList(ctx, userID, chatID)
s.chatRepo.AddToUsersChatsList(ctx, recipientID, chatID)

// Optimized — 1 round-trip with BulkWrite:
func (repo *chatRepository) AddToUsersChatsListBatch(ctx context.Context, chatID model.ChatID, userIDs []model.UserID) error {
    models := make([]mongo.WriteModel, 0, len(userIDs))
    for _, userID := range userIDs {
        models = append(models, mongo.NewUpdateOneModel().
            SetFilter(bson.M{"user_id": userID}).
            SetUpdate(bson.M{"$addToSet": bson.M{"chats_list_ids": chatID}}))
    }
    _, err := repo.usersCollection.BulkWrite(ctx, models)
    return err
}
```

### 2. Message Insert

If messages arrive in batches (e.g., syncing from another device), use `InsertMany` instead of multiple `InsertOne` calls.
