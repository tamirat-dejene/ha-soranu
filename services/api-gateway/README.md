# API Gateway

HTTP entrypoint for the microservices stack. It fronts the auth and user services, handles authentication via the auth service gRPC API, and exposes a small set of HTTP endpoints built on Gin. The code follows a lightweight clean-architecture layout and uses generated protobuf clients from the shared `authpb` package.

## Responsibilities

- Expose REST endpoints for authentication and basic user operations
- Forward auth flows to `auth-service` over gRPC and propagate JWTs
- Enforce authentication on protected routes via middleware (`Authorization: Bearer <token>`)
- Provide a single HTTP surface for clients and a place to hang cross-cutting middleware

## Tech Stack

- Go (module `github.com/tamirat-dejene/ha-soranu`, Go version per `go.mod`)
- Gin for HTTP routing
- gRPC (insecure transport for internal service-to-service calls)
- Viper for configuration
- Zap-backed logger (via shared `pkg/logger`)

## Project Layout

- `main.go` — process entrypoint; loads config, connects to auth service, starts Gin
- `internal/server` — HTTP server wiring and route registration
- `internal/api/handler` — HTTP handlers (auth)
- `internal/api/middleware` — auth middleware backed by the auth service
- `shared/` — shared configuration, logger, and generated protobuf clients
- `docs/swagger.yaml` — simple OpenAPI sketch of exposed endpoints

## Configuration

Set environment variables (usually via `app.env` next to `main.go`, Kubernetes ConfigMap, or Docker `-e` flags):

| Variable | Purpose | Default if missing |
| --- | --- | --- |
| `SERVER_PORT` | HTTP port the gateway listens on | `8080` |
| `AUTH_SERVICE_ADDR` | gRPC address of `auth-service` | `localhost:9090` |
| `USER_SERVICE_ADDR` | gRPC address of `user-service` (reserved for future routes) | none |
| `JWT_SECRET` | Shared JWT signing key (used by downstream services) | none |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Present for compatibility with shared config; currently unused by the gateway | none |

Kubernetes dev manifests also expect `GATEWAY_HTTP_ADDR`, `GATEWAY_GRPC_ADDR`, and `KAFKA_BROKERS` keys in the `app-config` ConfigMap (`infra/dev/k8s/app-config.yaml`).

## Quickstart (local)

1) Install prerequisites: Go toolchain, `protoc` with `protoc-gen-go` and `protoc-gen-go-grpc`, and `make`.
2) From repo root, move into the service: `cd services/api-gateway`.
3) Provide config (example `app.env`):

```env
SERVER_PORT=8080
AUTH_SERVICE_ADDR=localhost:9090
USER_SERVICE_ADDR=localhost:9091
JWT_SECRET=dev-secret
```

4) Run the gateway:

- Fast path: `go run .` (or `go run ./main.go`)
- With Makefile target (expects a `cmd` entrypoint): `make run`

5) Hit an endpoint:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
	-H "Content-Type: application/json" \
	-d '{"name":"Jane","email":"jane@example.com","password":"secret123"}'
```

## HTTP Surface

- `POST /api/v1/auth/register` — create an account via auth service; returns `user_id`
- `POST /api/v1/auth/login` — authenticate; returns a JWT token
- `GET /api/v1/user/profile` — protected example; requires `Authorization: Bearer <token>` header and echoes `user_id`

Swagger sketch lives in `docs/swagger.yaml`; update it if you add routes.

## gRPC Client Generation

`make proto-gen` regenerates Go stubs for any `.proto` files placed under `services/api-gateway/proto/` into `internal/api/grpc/pb`. It requires `protoc` and the Go plugins in `PATH`.

## Testing and Quality

- Unit tests: `make test`
- Lint (if `golangci-lint` is installed): `make lint`

## Container & Deployment Notes

- Development Dockerfile: `infra/dev/docker/api-gateway.Dockerfile` (expects a prebuilt binary at `bin/api-gateway` and copies shared packages).
- Kubernetes dev manifest: `infra/dev/k8s/api-gateway-deployment.yaml` (LoadBalancer on port `8080`, env from `app-config` ConfigMap).
- When running a container locally, supply required env vars, e.g.:

```bash
docker build -f infra/dev/docker/api-gateway.Dockerfile -t api-gateway .
docker run --rm -p 8080:8080 -e AUTH_SERVICE_ADDR=host.docker.internal:9090 api-gateway
```

## Troubleshooting

- Cannot connect to auth service: verify `AUTH_SERVICE_ADDR` is reachable from the gateway container/host.
- 401 on protected routes: ensure `Authorization: Bearer <token>` header is present and the token is issued by `auth-service`.