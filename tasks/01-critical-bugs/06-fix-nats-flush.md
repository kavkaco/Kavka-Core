# Fix NATS Flush() After Every Publish

**File:** `infra/stream/stream_pub.go:34`
**Priority:** P2 (Medium)
**Impact:** Unnecessary latency on every event publish

## Problem

`nc.Flush()` is called after every `Publish()` call. `Flush()` performs a synchronous round-trip to the NATS server to ensure the message buffer is sent. This:

1. Adds RTT latency (~1-10ms) to every event publish
2. Defeats the purpose of async pub/sub
3. Reduces throughput under high event volume

## Fix

Remove `nc.Flush()` from the `Publish` method. The NATS client library buffers messages internally and flushes them efficiently.

```go
func (p *pub) Publish(event *eventsv1.StreamEvent) error {
    eventBuf, err := proto.Marshal(event)
    if err != nil {
        return err
    }
    return p.nc.Publish(subjEvent, eventBuf)
}
```

The flush in `stream_sub.go:64` (`nc.Flush()` after Subscribe) is correct — it ensures the subscription is registered before returning.

## Validation

- Events are still delivered to subscribers
- No increase in event loss under normal conditions
- Reduced publish latency in benchmarks
