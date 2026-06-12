# Kavka Architecture: Distributed Systems & Event-Driven Design

## Overview

Kavka is a modular messenger backend designed for horizontal scalability.
It uses **NATS JetStream** as its event backbone, **Redis** for caching and rate limiting,
and supports **MongoDB**, **PostgreSQL**, and **SQLite** as storage engines.
The system follows Clean Architecture with pluggable repository implementations,
allowing independent scaling of read/write paths and real-time delivery.

---

## Event-Driven Core with NATS JetStream

### Subjects & Streams

| Stream | Subjects | Purpose | Retention |
|---|---|---|---|
| `kavka-events` | `events.>` | Real-time user events (message sent, chat created) | 7 days, file storage |
| `kavka-tasks` | `tasks.>` | Async background jobs (email, image processing) | 3 weeks, file storage |

### Event Flow

```
Client A ──SendMessage──→ API Gateway
                              │
                              ▼
                        Message Service
                              │
                    ┌─────────┴──────────┐
                    ▼                     ▼
              MongoDB/SQL            NATS JetStream
              (persist msg)       (publish events.message.sent)
                                        │
                          ┌─────────────┴─────────────┐
                          │                           │
                          ▼                           ▼
                    Event Subscriber            Consumer Workers
                    (real-time push)            (email, indexing)
                          │
                          ▼
                    Client B (via gRPC stream)
```

### Publisher

`StreamPublisher` publishes protobuf-encoded `StreamEvent` messages.
When JetStream is enabled (default), events are persisted and survive restarts.
When disabled, it falls back to core NATS pub/sub.

```go
// Every event carries receiver targeting:
event := &eventsv1.StreamEvent{
    SenderUserId:    senderID,
    ReceiversUserId: []model.UserID{recipientID},
    Payload:         payloadProtoBuf,
}
publisher.Publish(event)
```

### Subscriber (Real-Time Delivery)

`StreamSubscriber` maintains an in-memory map of `StreamSubscribedUser`,
each holding a channel connected to a gRPC server-streaming handler.
On receiving a NATS message, it matches `ReceiversUserId` against
subscribed users and pushes the event directly into their channel.

Concurrent access to the subscriber list is protected by `sync.RWMutex`.

### Consumer Groups (Task Workers)

For durable async processing, services subscribe as queue groups:

```
tasks.email.send ──→ email-workers (up to 3 replicas)
tasks.image.process ──→ image-workers (up to 5 replicas)
```

Each task includes `MaxRetries` and `RetryCount` for dead-letter handling.

---

## Caching Layer

Redis provides two distinct functions:

### Auth Token Storage (`go-auth-manager`)
Access tokens, refresh tokens, and email verification tokens
are stored in Redis with configurable TTLs.

### Application Cache (`infra/cache/`)

| Cache Key Pattern | TTL | What | Invalidation |
|---|---|---|---|
| `user:{userID}` | 5 min | User profile | On profile update |
| `chat:{chatID}` | 5 min | Chat metadata | On chat update |
| `chats:{userID}` | 2 min | User's chat list | On chat create/join |

The `Cache` interface is designed for future hot-path caching of
frequently accessed data without touching the primary database.

---

## Database Abstraction

Kavka supports three database backends via the `database.DriverType` config:

```
┌─────────────────────────────────────────────────────────────┐
│                    Repository Interfaces                    │
│  (internal/repository/)                                     │
├──────────────┬──────────────────┬───────────────────────────┤
│  repo_mongo  │   repo_sql       │   repo_sql (postgres)     │
│  (MongoDB)   │   (SQLite)       │   (PostgreSQL)            │
└──────────────┴──────────────────┴───────────────────────────┘
```

### Why JSON columns for chats?

Chats have a polymorphic `ChatDetail` field (channel, group, or direct).
Instead of separate tables per type (which would require joins or
table inheritance), SQL backends store the detail as a JSON column.
This mirrors MongoDB's subdocument approach and keeps queries simple.

### Why denormalized sender fields in messages?

Each message document stores `sender_name`, `sender_last_name`,
and `sender_username` alongside the `sender_id`. This eliminates
the `$lookup` (MongoDB) or `JOIN` (SQL) on every message fetch,
reducing query time by **~5-10x** under load.

---

## Rate Limiting

An optional Redis-backed rate limiter wraps gRPC endpoints via
a Connect interceptor:

| Endpoint | Limit | Window |
|---|---|---|
| Login | 10 req | 1 min |
| Register | 3 req | 1 hour |
| SendMessage | 60 req | 1 min |
| CreateGroup | 10 req | 1 hour |
| All others | 100 req | 1 min |

---

## Message Schema Versions

**V1** (legacy, default): Embedded array per chat document.
Suitable for small deployments but hits MongoDB's 16MB BSON limit
at ~80,000 messages per chat.

**V2** (recommended for scale): One document per message with
denormalized sender info. No document size limit, supports
cursor-based pagination, and works natively with SQL backends.

Enable V2 via config: `app.use_messages_v2: true`
Run migration: `go run cmd/migration/main.go`

---

## Scaling Considerations

### Stateless Services
Each service instance is stateless. Session is stored in Redis,
events flow through NATS, and data persists in the database.
Add replicas behind a load balancer for horizontal scaling.

### Database Backend Choice
- **MongoDB**: Best for rapid development, schema flexibility
- **PostgreSQL**: Best for production at scale, ACID compliance, replication
- **SQLite**: Best for local dev, single-binary deployment, testing

### NATS Cluster
For production, run a 3-node NATS cluster with JetStream enabled.
Configure `nats.MaxReconnects(10)` in `stream_adapter.go` for resilience.
