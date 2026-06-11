# Add Pagination to All List Endpoints

**Priority:** P1 (High)
**Impact:** OOM server when chats have 100K+ messages

## Problem

These endpoints return ALL results without pagination:
1. `FetchMessages` — returns every message in a chat
2. `GetUserChats` — returns all user's chats with full join data
3. `Search` — returns up to 10 results but no cursor/offset

## Implementation

### FetchMessages

Add `limit` and `offset` parameters:

```go
type FetchMessagesParams struct {
    ChatID  model.ChatID
    Limit   int64   // default: 50, max: 200
    Offset  int64   // for cursor-based: messageID
}

// In repository:
pipeline := bson.A{
    {"$match": bson.M{"chat_id": chatID}},
    {"$sort": bson.M{"created_at": -1}},
    {"$skip": offset},
    {"$limit": limit},
}
```

### GetUserChats

Add pagination support to the `GetUserChats` repository method. Since `chatIDs` come from the user's `chats_list_ids`, paginate at the application level:

```go
func (s *ChatService) GetUserChats(ctx context.Context, userID model.UserID, page, limit int) ([]model.ChatDTO, *vali.ValiErr) {
    user, _ := s.userRepo.FindByUserID(ctx, userID)
    totalChats := user.ChatsListIDs
    
    // Calculate offset
    start := (page - 1) * limit
    end := min(start + limit, len(totalChats))
    pageIDs := totalChats[start:end]
    
    userChats, _ := s.chatRepo.GetUserChats(ctx, userID, pageIDs)
    return userChats, nil
}
```

### Search

Add `page`/`limit` to `SearchRepository.Search` method signature:
```go
Search(ctx context.Context, input string, limit, offset int64) (*model.SearchResultDTO, error)
```

## Validation

- Fetching messages with large datasets returns only the requested page
- Memory usage is bounded regardless of chat size
- All pagination tests pass
- Backward compatible (default limit = 50)
