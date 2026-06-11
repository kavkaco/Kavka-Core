# Message Schema Redesign

**Priority:** P1 (High)
**Impact:** Unbounded array growth hits 16MB BSON limit; no individual message querying

## Problem

Current schema stores all messages for a chat in a single document:

```bson
// Collection: messages
{
  chat_id: ObjectId,
  messages: [
    { message_id, sender_id, content, created_at, ... },   // message 1
    { message_id, sender_id, content, created_at, ... },   // message 2
    // ... up to ~200K messages before hitting 16MB limit
  ]
}
```

Problems:
1. **16MB BSON limit** — at ~200 bytes/message, ~80K messages fills the doc
2. **No individual message querying** — every operation reads the entire array
3. **$push without cap** — unbounded growth with no eviction strategy
4. **$lookup on every fetch** — sender data joined from users collection

## New Schema

```bson
// Collection: messages_v2
{
  _id: ObjectId (message_id),
  chat_id: ObjectId,
  sender_id: string,
  sender_name: string,        // ← denormalized
  sender_last_name: string,   // ← denormalized
  sender_username: string,    // ← denormalized
  created_at: DateTime,
  edited: bool,
  seen: bool,
  type: "text" | "image" | "label",
  content: { text: "..." } | { image_url: "...", caption: "..." }
}

// Indexes:
// { chat_id: 1, created_at: -1 }   — for fetching messages by chat (sorted)
// { chat_id: 1, message_id: 1 }    — for individual message lookup
```

## Migration Strategy

1. Create `messages_v2` collection with new schema
2. Write a migration script that reads all old messages and inserts them as individual documents
3. Update all repository methods to use the new collection
4. Deploy, then drop `messages` collection after confirming data integrity

## Benefits

- No 16MB limit — each message is its own document
- Pagination with `$skip/$limit` or cursor-based (`_id > lastSeen`)
- Efficient `$lookup` eliminated (sender info denormalized)
- Individual message updates without reading entire array
- `$text` search index on message content
