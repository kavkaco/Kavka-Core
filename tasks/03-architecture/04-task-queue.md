# Asynchronous Task Queue

**Files:** `infra/queue/` (currently empty)
**Priority:** P2 (Medium)
**Impact:** Email sending blocks request handler, no retry on failure

## Current State

`infra/queue/` contains only a `.gitkeep` file. Email is sent synchronously in the auth handler. No async task processing exists.

## Target: Queue Interface

```go
package queue

type Task struct {
    ID        string
    Type      string
    Payload   []byte
    Retries   int
    MaxRetries int
}

type Queue interface {
    Enqueue(ctx context.Context, task Task) error
    Dequeue(ctx context.Context, queue string) (*Task, error)
    Ack(ctx context.Context, taskID string) error
    Nack(ctx context.Context, taskID string) error
}
```

## NATS JetStream Implementation

```go
type natsQueue struct {
    js nats.JetStreamContext
}

func (q *natsQueue) Enqueue(ctx context.Context, task Task) error {
    data, _ := json.Marshal(task)
    _, err := q.js.Publish("commands."+task.Type, data)
    return err
}
```

## Task Types

| Task | Queue | Consumer | Retry Policy |
|---|---|---|---|
| `email.verify` | `commands.email.*` | Email Worker | 3 retries, 1min backoff |
| `email.reset_password` | `commands.email.*` | Email Worker | 3 retries, 1min backoff |
| `image.process` | `commands.image.*` | Image Worker | 2 retries, 5min backoff |
| `notification.deliver` | `commands.notify.*` | Notify Worker | 5 retries, 30s backoff |

## Implementation Priority

1. Define `Queue` interface in `infra/queue/`
2. Implement NATS JetStream-backed queue
3. Create Email Worker service
4. Replace sync email sending with async enqueue
5. Add retry + dead-letter handling
6. Add worker pool with configurable concurrency
