# Kavka-Core Improvement Tasks

## Overview

This directory contains the complete task breakdown for transforming Kavka-Core from a modular monolith into a production-grade, event-driven microservice architecture capable of handling heavy loads.

## Phase Structure

| Phase | Focus | Status |
|---|---|---|
| **01-critical-bugs** | P0/P1 bugs that break functionality or cause data corruption | ⏳ In Progress |
| **02-database** | Schema redesign, missing indexes, pagination | ⏳ In Progress |
| **03-architecture** | Event-driven design with NATS JetStream, sagas, cache, task queue | 📝 Planned |
| **04-microservices** | Service extraction, API gateway, inter-service communication | 📝 Planned |
| **05-performance** | Connection pooling, rate limiting, read/write separation, bulk ops | 📝 Planned |

## Implementation Order

1. Fix critical bugs first (data corruption, broken queries)
2. Fix context propagation and concurrency safety
3. Schema migration for messages (unbounded array → separate documents)
4. Add pagination to all list endpoints
5. NATS pub/sub → JetStream with durable consumers
6. Add caching layer (Redis)
7. Add async task queue (NATS JetStream consumer groups)
8. Extract microservices one by one
9. Add rate limiting, connection pool tuning

## Current State

- **Language:** Go 1.22
- **Architecture:** Clean Architecture modular monolith
- **Storage:** MongoDB (primary), Redis (tokens), MinIO (objects)
- **Messaging:** NATS pub/sub (basic, no persistence)
- **Total:** 83 Go source files, ~7,300 LOC
