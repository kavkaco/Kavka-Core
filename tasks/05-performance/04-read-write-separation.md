# Read/Write Separation

**Priority:** P3 (Low)
**Impact:** Offload reads to secondaries, increase write throughput

## Strategy

Use MongoDB replica set with primary for writes and secondary for reads.

### Connection Setup

```go
// Primary (writes)
primaryClient, _ := mongo.Connect(ctx, options.Client().
    ApplyURI("mongodb://primary:27017").
    SetWriteConcern(writeconcern.New(writeconcern.WMajority())))

// Secondary (reads)
secondaryClient, _ := mongo.Connect(ctx, options.Client().
    ApplyURI("mongodb://secondary:27017,secondary2:27017").
    SetReadPreference(readpref.SecondaryPreferred()))
```

### Repository Changes

Inject both clients into repositories:

```go
type chatRepository struct {
    primary   *mongo.Collection  // writes
    secondary *mongo.Collection  // reads
}
```

### Query Routing

| Operation | Collection | Target |
|---|---|---|
| INSERT / UPDATE / DELETE | All | Primary |
| `Find` (single) | Users, Chats | SecondaryPreferred |
| `Aggregate` (list) | Messages, Chats | SecondaryPreferred |
| `FindOne` (auth) | user_auth | Primary (read-your-writes) |
