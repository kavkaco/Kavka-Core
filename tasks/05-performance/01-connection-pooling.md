# MongoDB Connection Pool Tuning

**File:** `database/mongo_adapter.go`
**Priority:** P2 (Medium)

## Current
```go
client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
```
Default `maxPoolSize=100`.

## Tuned for Heavy Load
```go
client, err := mongo.Connect(ctx, options.Client().
    ApplyURI(uri).
    SetMaxPoolSize(500).
    SetMinPoolSize(50).
    SetMaxConnIdleTime(30*time.Second).
    SetMaxConnecting(50), // limit concurrent connection creation
)
```

## Additional Settings

| Setting | Value | Reason |
|---|---|---|
| `maxPoolSize` | 500 | Handle high concurrency |
| `minPoolSize` | 50 | Avoid cold starts |
| `maxConnIdleTime` | 30s | Reclaim idle connections |
| `maxConnecting` | 50 | Throttle connection storms |
| `socketTimeout` | 10s | Bound query execution time |
| `connectTimeout` | 5s | Quick fail on network issues |
| `heartbeatInterval` | 10s | Detect replica set changes |

## Validation
- Benchmark under load shows no connection starvation
- No `ErrPoolTimeout` in logs
- Connection count stabilizes at `minPoolSize` during idle
