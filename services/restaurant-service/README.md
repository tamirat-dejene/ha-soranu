# Restaurant-Service Service

This is the **restaurant-service microservice**, built using Clean Architecture principles in Go.

## Structure

- `cmd/` — CLI entrypoints
- `cmd/grpc/` — example gRPC entrypoint
- `proto/` — protobuf definitions
- `internal/api/grpc/pb/` — generated gRPC Go code
- `internal/` — core logic (domain, usecase, infra, api)
- `pkg/` — reusable helper packages
- `db/` — migrations
- `docs/` — documentation

## Commands

```bash
make proto-gen
make run
make run-grpc
make test
make migrate-up
```
