# Payment Service

Handles payment intent creation, Stripe webhooks, and order status updates. Exposes a small HTTP API and starts a gRPC server (reserved for future RPCs). Automatically runs DB migrations on startup.

## Architecture
- Language: Go
- Storage: PostgreSQL (payments table)
- Payments: Stripe Payment Intents
- Messaging: Kafka (publishes order status updates)
- HTTP server: payment intents, retrieval, webhooks
- gRPC server: placeholder for future payment RPCs

Key files:
- [services/payment-service/cmd/main.go](cmd/main.go): wiring, servers, migrations
- [services/payment-service/internal/api/http/server.go](internal/api/http/server.go): HTTP routes
- [services/payment-service/internal/usecase/service.go](internal/usecase/service.go): business logic, Stripe + events
- [services/payment-service/internal/repository/postgres.go](internal/repository/postgres.go): Postgres persistence
- [services/payment-service/migrations](migrations): Goose migrations (auto-applied on startup)

## Ports
- HTTP: `PAYMENT_HTTP_PORT` (default 8081)
- gRPC: `PAYMENT_SRV_PORT` (default 9090)

## HTTP API
Base URL: `http://localhost:<PAYMENT_HTTP_PORT>`

- GET `/healthz`
  - Returns `200 OK` with body `ok`.

- POST `/payments/intent`
  - Creates a Stripe Payment Intent and a local `payment` record.
  - Body:
    ```json
    {
      "order_id": "<uuid>",
      "amount": 1200,
      "currency": "usd"
    }
    ```
  - Response:
    ```json
    {
      "payment_id": "<uuid>",
      "client_secret": "<stripe_client_secret>"
    }
    ```

- GET `/payments/{id}`
  - Returns the persisted payment by ID.

- POST `/payments/webhook`
  - Stripe webhook endpoint. Verifies `Stripe-Signature` using `STRIPE_WEBHOOK_SECRET`.
  - Handles events:
    - `payment_intent.succeeded` → updates payment to `succeeded` and publishes order status `COMPLETED`.
    - `payment_intent.payment_failed` → updates payment to `failed` and publishes order status `CANCELLED`.

## Domain & Events
- Payment statuses: `pending | succeeded | failed | canceled`.
- On webhook processing, publishes `OrderStatusUpdated` (see shared `orderpb`). Kafka broker is configured via `KAFKA_BROKER_URL`.

## Environment Variables
Defaults come from [services/payment-service/env.go](env.go).

- Service
  - `SRV_ENV`: service environment (default: `development`).
  - `PAYMENT_SRV_NAME`: service name (default: `payment-service`).
  - `PAYMENT_SRV_PORT`: gRPC port (default: `9090`).
  - `PAYMENT_HTTP_PORT`: HTTP port (default: `8081`).

- PostgreSQL
  - `POSTGRES_HOST` (default: `postgres-db`)
  - `POSTGRES_PORT` (default: `5432`)
  - `POSTGRES_USER` (default: `postgres`)
  - `POSTGRES_PASSWORD` (default: `password`)
  - `POSTGRES_DB` (default: `payment-servicedb`)

- Redis (not currently used in handlers, reserved for caching)
  - `REDIS_HOST` (default: `localhost`)
  - `REDIS_PORT` (default: `6379`)
  - `REDIS_PASSWORD` (default: empty)
  - `REDIS_DB` (default: `0`)

- Stripe
  - `STRIPE_SECRET_KEY` (required for creating payment intents)
  - `STRIPE_WEBHOOK_SECRET` (required for webhook signature verification)

- Kafka
  - `KAFKA_BROKER_URL` (default: `localhost:9092`)

## Database Migrations
Migrations are applied automatically on startup using Goose (see
[services/payment-service/migrations/migrate.go](migrations/migrate.go)).
On first run, the service will create the database if missing and run `Up` migrations from `migrations` (container path `/app/payment-svc-migrations`).

To run manually (optional), install goose and run against your DSN.

## Local Development
### Prerequisites
- Go 1.21+
- PostgreSQL running and reachable
- Kafka broker (for event publishing)
- Stripe account with a test key and webhook signing secret

### Run the service
Set the necessary environment and run the main package:

```bash
export SRV_ENV=development
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=payment_servicedb
export KAFKA_BROKER_URL=localhost:9092
export STRIPE_SECRET_KEY=sk_test_xxx
export STRIPE_WEBHOOK_SECRET=whsec_xxx
export PAYMENT_HTTP_PORT=8081
export PAYMENT_SRV_PORT=9090

go run ./services/payment-service/cmd
```

### Quick checks
- Health:
  ```bash
  curl -i http://localhost:8081/healthz
  ```
- Create payment intent:
  ```bash
  curl -s -X POST http://localhost:8081/payments/intent \
    -H 'Content-Type: application/json' \
    -d '{"order_id":"11111111-1111-1111-1111-111111111111","amount":1200,"currency":"usd"}'
  ```
- Get by id:
  ```bash
  curl -s http://localhost:8081/payments/<payment_id>
  ```

### Stripe webhook (local)
Use the Stripe CLI to forward webhooks:

```bash
stripe listen --forward-to localhost:8081/payments/webhook
```

## Docker & Kubernetes
- Dockerfile: [infra/dev/docker/payment-service.Dockerfile](../../infra/dev/docker/payment-service.Dockerfile)
- K8s manifest: [infra/dev/k8s/payment-service-deployment.yaml](../../infra/dev/k8s/payment-service-deployment.yaml)

Example Docker run (environment must provide Postgres, Kafka, Stripe secrets):
```bash
docker build -f infra/dev/docker/payment-service.Dockerfile -t payment-service:dev .

docker run --rm -p 8081:8081 -p 9090:9090 \
  -e POSTGRES_HOST=host.docker.internal \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=payment_servicedb \
  -e KAFKA_BROKER_URL=host.docker.internal:9092 \
  -e STRIPE_SECRET_KEY=sk_test_xxx \
  -e STRIPE_WEBHOOK_SECRET=whsec_xxx \
  --name payment-svc payment-service:dev
```

## Notes
- The gRPC server currently exposes no RPCs; it boots to reserve the port and enable future expansion.
- If `STRIPE_SECRET_KEY` or `STRIPE_WEBHOOK_SECRET` are missing, intent creation or webhook validation will fail.
- Database schema is created via migrations; ensure DB connectivity on boot.