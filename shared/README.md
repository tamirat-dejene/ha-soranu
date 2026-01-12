# Shared Library

## Overview

The `shared` module contains cross-cutting libraries reused by all services: database access, caching, logging, messaging (Kafka), auth helpers, and generated protobuf code. Everything is Go 1.x compatible and kept framework-agnostic so services can import without extra wrappers.

## Structure

- `db/pg/` – thin pgx pool wrapper with `PostgresClient` interface for querying and transactions.
- `pkg/logger/` – Zap-based logger with env-aware config.
- `pkg/events/` – domain event names and Kafka publisher for order-related events.
- `pkg/messaging/kafka/` – transport abstractions (`Producer`, `Consumer`, `Message`) plus Sarama implementations.
- `pkg/caching/` – cache abstraction with Redis client implementation.
- `pkg/auth/` – JWT validation helpers for RSA access/refresh tokens.
- `pkg/utils/` – reserved for generic utilities (currently empty).
- `protos/` – generated Go protobuf files for shared proto definitions.

## Key Components

### Database (PostgreSQL)

- `PostgresClient` interface provides `Query`, `QueryRow`, `Exec`, `BeginTx`, and `Close` with pgx under the hood.
- Connection helper: create DSN (`postgres://user:pass@host:port/db?sslmode=disable`) and call `postgres.NewPostgresClient`.
- Transactions use `BeginTx` to obtain a `Tx` with the same query surface.

### Logging

- Init with `logger.InitLogger(env)` where `env` is `development` or `production`.
- Convenience functions: `logger.Info/Error/Debug/Warn/Fatal` and `logger.With` for contextual loggers.

### Events

- Event names: `order.placed`, `order.shipped`, `order.cancelled`, `order.status_updated`.
- Publisher: `events.EventPublisher` wraps protobuf marshaling into an envelope (`EventEnvelope`) and publishes to Kafka topics matching the event type.

### Messaging (Kafka)

- Transport-agnostic interfaces in `pkg/messaging/kafka`.
- Sarama-backed producer/consumer in `pkg/messaging/kafka/sarama` with sensible defaults (Kafka 4.1 client version, newest offsets, sync producer).
- Usage: construct `sarama.NewProducer([]string{broker})` or `sarama.NewConsumer(brokers, groupID)` and pass to service use cases.

### Caching

- `CacheClient` interface abstracts cache operations (get/set, incr/decr, expire, ping).
- Redis implementation in `pkg/caching/redis` using go-redis.

### Auth Helpers

- RSA JWT validators in `pkg/auth/jwtvalidator`: parse PEM strings and validate access/refresh tokens into typed claims.

## Protobufs

- Source protos live in the repository root `protos/`. Generated Go code is staged under `shared/protos/` for reuse.
- Regenerate after proto changes:

```bash
make proto
```

## Usage Patterns (snippets)

> Adjust package import paths to the module root `github.com/tamirat-dejene/ha-soranu`.

- Postgres client:

```go
pg, _ := postgres.NewPostgresClient("postgres://user:pass@host:5432/db?sslmode=disable")
defer pg.Close()
row := pg.QueryRow(ctx, "SELECT 1")
```

- Logger:

```go
logger.InitLogger("development")
logger.Info("service starting")
```

- Kafka consumer (Sarama):

```go
consumer, _ := sarama.NewConsumer([]string{"localhost:9092"}, "group-id")
defer consumer.Close()
_ = consumer.Subscribe(ctx, []string{"order.placed"}, handler)
```

- Redis cache:

```go
cache, _ := redis.NewRedisClient("localhost", 6379, "", 0)
defer cache.Close()
_ = cache.Set(ctx, "key", "value", time.Minute)
```

## Conventions

- Keep packages dependency-light and framework-neutral.
- Prefer interfaces in `pkg/...` so services can mock during tests.
- Avoid embedding service-specific logic here; keep it reusable and generic.
