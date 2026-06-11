# Microservice Split Plan

**Priority:** P2 (Medium)
**Impact:** Enables independent scaling, deployment, and team ownership

## Service Boundaries

| Microservice | Port | Dependencies | Replicas |
|---|---|---|---|
| **Auth Service** | 8081 | MongoDB (auth), Redis (tokens) | 2-3 |
| **Chat Service** | 8082 | MongoDB (chats, users) | 3-5 |
| **Message Service** | 8083 | MongoDB (messages), NATS | 5-10 |
| **Search Service** | 8084 | MongoDB (users, chats, messages) | 1-2 |
| **User Service** | 8085 | MongoDB (users), MinIO | 2-3 |
| **Notification Service** | 8086 | NATS (subscriber only) | 2-3 |
| **API Gateway** | 8080 | All services | 2-3 |

## Migration Order

1. **Extract Auth Service** (no dependencies on other services)
2. **Extract User Service** (depends only on MongoDB)
3. **Extract Chat Service** (depends on User + Message repos)
4. **Extract Message Service** (depends on Chat + User)
5. **Extract Search Service** (read-only, depends on all)
6. **Extract Notification Service** (NATS subscriber only)
7. **Add API Gateway** that replaces the monolith's gRPC server

## Shared Library

Create a shared Go module for common types and interfaces:

```bash
kavka-common/
├── model/           # Shared domain types
├── repository/      # Interface contracts
├── proto/           # Generated protobuf
├── stream/          # NATS publisher/subscriber
├── cache/           # Cache interface
├── queue/           # Queue interface
└── middleware/      # Auth interceptor
```

## Communication

- **Synchronous:** gRPC + Connect-Web (existing protobuf)
- **Asynchronous:** NATS JetStream (existing infra)
- **Service-to-service auth:** Internal JWT tokens
