# Kavka Setup & Deployment Guide

## Prerequisites

- Go 1.22+
- Docker & Docker Compose (for development)
- Access to a MongoDB, PostgreSQL, or SQLite instance
- NATS server (with JetStream enabled for production)

---

## Quick Start (Development)

### 1. Clone & Dependencies

```bash
git clone https://github.com/kavkaco/Kavka-Core
cd Kavka-Core

# Download Go dependencies
go mod tidy

# Generate JWT secret key
openssl genpkey -algorithm RSA -out config/jwt_secret_key.pem -pkeyopt rsa_keygen_bits:2048
```

### 2. Start Infrastructure (Docker)

```bash
docker compose up -d
```

This starts MongoDB, Redis, MinIO, and NATS with default development
credentials defined in `docker-compose.yml`.

### 3. Configure

Edit `config/config.development.yml`:

```yaml
sql:
    driver: "mongodb"        # Options: mongodb, postgres, sqlite
    # ... connection details per driver (see below)

app:
    use_messages_v2: false   # Set true after running migration
```

### 4. Run

```bash
KAVKA_ENV=development go run ./cmd/server/
```

Server starts on `0.0.0.0:8000` with h2c (HTTP/2 cleartext) + Connect-Web gRPC.

---

## Database Configuration

### MongoDB

```yaml
sql:
    driver: "mongodb"

mongo:
    host: "127.0.0.1"
    port: 27017
    username: "mongo"
    password: "mongo"
    db_name: "kavka"
```

The `mongo` block is required even when not using MongoDB ---
the config loader expects it. For SQL drivers, set `sql.driver` accordingly.

### PostgreSQL

```yaml
sql:
    driver: "postgres"
    host: "127.0.0.1"
    port: 5432
    username: "postgres"
    password: "postgres"
    db_name: "kavka"
    ssl_mode: "disable"      # Use "require" in production
```

### SQLite

```yaml
sql:
    driver: "sqlite"
    sqlite_path: "./data/kavka.db"
```

SQLite requires no external server. The database file is created
automatically on first run.

---

## Running with Docker Compose (Development)

The included `docker-compose.yml` provides everything needed for local development:

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f app

# Stop everything
docker compose down
```

Default services and ports:

| Service | Port | Credentials |
|---|---|---|
| `app` (Go server) | 8000 | — |
| `mongo` | 27017 | mongo:mongo |
| `redis` | 6379 | redis: |
| `nats` | 4222 | — |
| `minio` | 9000, 9001 | minio:minio123 |

---

## Production Deployment

### Architecture Recommendation

```
                         ┌─────────────┐
                         │  Load Balancer │
                         │   (HTTPS/TLS)  │
                         └───────┬───────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
              ┌─────▼─────┐           ┌───────▼──────┐
              │  App       │  ... N    │   App        │
              │  Instance  │           │   Instance   │
              └─────┬─────┘           └───────┬──────┘
                    │                         │
         ┌──────────┴─────────────────────────┴──────────┐
         │                    NATS Cluster               │
         │              (3 nodes, JetStream)              │
         └──────────┬─────────────────────────┬──────────┘
                    │                         │
         ┌──────────▼──────┐        ┌─────────▼─────────┐
         │  PostgreSQL      │        │   Redis Sentinel   │
         │  (Primary+Replica)│       │   (or Cluster)    │
         └─────────────────┘        └───────────────────┘
```

### Production Checklist

1. **Database**: Use PostgreSQL with streaming replication.
   Set `ssl_mode: "require"` and configure `max_pool_size: 500`.

2. **NATS**: Deploy 3+ nodes with JetStream file storage.
   Configure `nats.MaxReconnects(10)` for resilience.

3. **Redis**: Deploy with Sentinel or Redis Cluster for HA.
   Rate limiting and token storage tolerate eventual consistency.

4. **TLS**: Terminate TLS at the load balancer. The app speaks
   h2c internally. Set `http.cors.allow_origins` to your domain.

5. **Secrets**: Generate a unique `jwt_secret_key.pem` per deployment.
   Never commit it to version control.

6. **Message Schema V2**: For production at scale, enable V2:

   ```bash
   # Migrate data first
   go run cmd/migration/main.go

   # Then enable in config
   app:
       use_messages_v2: true
   ```

7. **Environment Variables**:

   ```bash
   KAVKA_ENV=production
   # Or set via config filename: config.production.yml
   ```

### Kubernetes Deployment

Minimal Kubernetes manifests are provided in `k8s/` (if present).
Key settings for the deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kavka-core
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: app
        image: kavkaco/kavka-core:latest
        ports:
        - containerPort: 8000
        env:
        - name: KAVKA_ENV
          value: "production"
        volumeMounts:
        - name: config
          mountPath: /app/config
        - name: jwt-secret
          mountPath: /app/config/jwt_secret_key.pem
            subPath: jwt_secret_key.pem
```

### Health Checks

The server exposes pprof endpoints in development mode.
For production, add health check endpoints:

- `GET /healthz` — Database connectivity, NATS connection
- `GET /readyz` — Full readiness (all dependencies reachable)

### Monitoring

- **Metrics**: Export Go runtime metrics via pprof endpoints
- **Logging**: Structured JSON logs via zerolog (configurable targets)
- **Tracing**: OpenTelemetry-ready (add exporter in main.go)

---

## Development vs Production: Key Differences

| Aspect | Development | Production |
|---|---|---|
| Database | MongoDB (default) or SQLite | PostgreSQL with replication |
| NATS | Single node, in-memory | Cluster with JetStream file storage |
| TLS | h2c (no TLS) | Terminated at load balancer |
| Secrets | Dev PEM in repo | Vault / K8s Secrets |
| Logging | Console + file (debug) | JSON, structured (info) |
| Profiling | pprof enabled | pprof disabled |
| Message V2 | Disabled (V1) | Enabled after migration |
| CORS | `*` allowed | Specific origins |
| Rate Limiting | Disabled | Enabled |
