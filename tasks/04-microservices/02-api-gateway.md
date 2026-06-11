# API Gateway

**Priority:** P2 (Medium)
**Impact:** Single entry point, centralized auth, routing

## Architecture

```
Client → [API Gateway :8080] → Auth Service :8081
                              → Chat Service :8082
                              → Message Service :8083
                              → Search Service :8084
                              → User Service :8085
```

## Responsibilities

1. **TLS termination**
2. **Rate limiting** per client IP/user
3. **Auth interceptor** — validate access token, inject user_id
4. **Request routing** — proxy to appropriate service
5. **CORS** handling
6. **Request/response logging**

## Implementation

The existing `delivery/grpc/grpc_server.go` is already the gateway — it becomes a reverse proxy:

```go
// Instead of injecting services directly:
// Before:
chatHandler := grpc_handlers.NewChatGrpcHandler(services.ChatService)

// After:
chatHandler := grpc_handlers.NewChatGrpcClient("http://chat-service:8082")
```

For the first phase, keep the monolith pattern but add rate limiting using the existing Redis instance.
