# NATS JetStream Migration

**Priority:** P1 (High)
**Impact:** Event loss on subscriber disconnect, no replay capability

## Current State

NATS is used as a simple pub/sub with no persistence. Events are published on the `"events"` subject and delivered to all connected subscribers. If a subscriber disconnects, events are lost.

## Target State

Use NATS JetStream for:
1. **Durable streams** — events persist on disk
2. **Consumer groups** — work queues for async processing
3. **At-least-once delivery** — messages retried on failure
4. **Event replay** — new services can replay past events

## Implementation

### Step 1: Create Stream

```go
js, _ := nc.JetStream()
js.AddStream(&nats.StreamConfig{
    Name:      "kavka-events",
    Subjects:  []string{"events.>"},
    Retention: nats. LimitsPolicy,
    MaxAge:    7 * 24 * time.Hour, // 7 days retention
    Storage:   nats.FileStorage,
})
```

### Step 2: Publish with Dedup

```go
func (p *pub) Publish(event *eventsv1.StreamEvent) error {
    eventBuf, _ := proto.Marshal(event)
    
    // Use JetStream with dedup ID for idempotency
    _, err := p.js.Publish("events."+eventType, eventBuf,
        nats.MsgId(event.IdempotencyKey),
    )
    return err
}
```

### Step 3: Durable Consumers

```go
// Real-time notification consumer
js.Subscribe("events.>", handler,
    nats.Durable("notifications"),
    nats.DeliverAll(),
)

// Email worker consumer (queue group)
js.QueueSubscribe("events.user.registered", "email-workers", handler,
    nats.Durable("email-worker"),
    nats.ManualAck(),
)
```

## Event Categories

| Subject | Retention | Consumers |
|---|---|---|
| `events.user.*` | 7 days | Email, Analytics, Notification |
| `events.chat.*` | 7 days | Notification, Search Index |
| `events.message.*` | 7 days | Notification, Search Index |
| `commands.*` | Until processed | Worker queues |

## Migration

1. Add JetStream stream configuration
2. Update `StreamPublisher` to use JetStream
3. Update `StreamSubscriber` to use durable consumer
4. Add consumer groups for async workers (email, image processing)
5. Add retry and dead-letter queues
