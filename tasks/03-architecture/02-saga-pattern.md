# Saga Pattern for Distributed Transactions

**Priority:** P2 (Medium)
**Impact:** Data inconsistency when multi-step operations fail mid-way

## Problem

Creating a chat currently involves:
1. `chatRepo.Create()` → INSERT into `chats`
2. `messageRepo.Create()` → INSERT into `messages` (message store)
3. `chatRepo.AddToUsersChatsList(userID)` → UPDATE `users`
4. `chatRepo.AddToUsersChatsList(recipientID)` → UPDATE `users`
5. Event publish → NATS

If step 2 fails, step 1 has already created an orphan chat document. If step 4 fails, one user gets the chat but the other doesn't.

## Solution: Choreography-Based Saga

Each service publishes events after successful operations, and subscribes to compensating events:

```go
// Chat Service
func (s *ChatService) CreateDirect(ctx, userID, recipientID) (*model.ChatDTO, error) {
    // 1. Create chat document
    chat, err := s.chatRepo.Create(ctx, chatModel)
    
    // 2. Publish "chat.created" event (includes saga_id)
    s.eventBus.Publish("events.chat.created", SagaEvent{
        SagaID: sagaID,
        ChatID: chat.ChatID,
        Steps:  []string{"chat_created"},
    })
    
    return chat, nil
}

// Message Store Service (separate consumer)
func HandleChatCreated(event SagaEvent) {
    err := s.messageRepo.Create(ctx, event.ChatID)
    if err != nil {
        // Publish compensating event
        s.eventBus.Publish("events.chat.compensate", event)
        return
    }
    s.eventBus.Publish("events.chat.store_created", event.NextStep())
}
```

## Compensation Handlers

| Step | Compensating Action | Event |
|---|---|---|
| Chat created | `chatRepo.Destroy()` | `events.chat.compensate` |
| Message store created | `messageRepo.Destroy()` | `events.chat.store_compensate` |
| User chat list updated | `userRepo.RemoveFromChatsList()` | `events.chat.list_compensate` |
| Event published | No-op (idempotent) | N/A |

## Implementation Priority

1. Add `SagaID` to event structures
2. Extract chat creation into event-driven workflow
3. Add compensating event handlers
4. Add dead-letter handling for failed sagas
5. Add saga state tracking (in-progress, completed, failed)
