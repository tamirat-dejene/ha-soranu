package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var structure = []string{
	"cmd",
	"cmd/grpc",
	"db/migrations",
	"docs",
	"proto",

	"internal/api/http/handler",
	"internal/api/http/middleware",
	"internal/api/http/route",

	"internal/api/grpc",
	"internal/api/grpc/handler",
	"internal/api/grpc/pb",

	"internal/config",

	"internal/domain/entity",
	"internal/domain/repository",
	"internal/domain/cache",
	"internal/domain/errs",

	"internal/infra/db/mocks",
	"internal/infra/db/model_gen",
	"internal/infra/db/repository",
	"internal/infra/redis/mocks",
	"internal/infra/redis/repository",

	"internal/server",
	"internal/usecase",
	"internal/util/httpresponse",

	"pkg/logger",
	"pkg/errors",
	"test",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run newservice.go <service-name>")
		os.Exit(1)
	}

	service := os.Args[1]
	service = strings.TrimSpace(service)
	if service == "" {
		fmt.Println("Invalid service name")
		os.Exit(1)
	}

	root := filepath.Join("services", service)
	fmt.Printf("Creating clean architecture structure for %s (under services/)...\n\n", service)

	for _, dir := range structure {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf(" - Failed to create %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf(" - Created: %s\n", path)
	}

	// Boilerplate files
	createFile(filepath.Join(root, "main.go"), mainTemplate(service))
	createFile(filepath.Join(root, "Makefile"), makefileTemplate())
	createFile(filepath.Join(root, ".gitignore"), gitignoreTemplate())
	createFile(filepath.Join(root, "Dockerfile"), dockerfileTemplate(service))
	createFile(filepath.Join(root, "docker-compose.yaml"), dockerComposeTemplate(service))
	createFile(filepath.Join(root, "README.md"), readmeTemplate(service))
	createFile(filepath.Join(root, "go.mod"), gomodTemplate(service))

	fmt.Println("\n - Done! Service structure created successfully.")
	fmt.Println("Next steps:")
	fmt.Println("  cd services/" + service)
	fmt.Println("  go mod tidy")
	fmt.Println("  make run")
}

func createFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf(" - Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf(" - Created: %s\n", path)
}

func mainTemplate(service string) string {
	return fmt.Sprintf(`package main

import (
"context"
"fmt"
"log"
"net"

"google.golang.org/grpc"
)

func main() {
fmt.Println(" - Starting %s service...")

// TODO: Initialize config, DI, repositories

// Start gRPC server (example)
go func() {
lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("failed to listen: %%v", err)
		}
s := grpc.NewServer()
// TODO: register gRPC services here, e.g. pb.RegisterYourServiceServer(s, &handler.YourService{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %%v", err)
		}
}()

// Placeholder: keep process alive
select {}
}
`, service)
}

func makefileTemplate() string {
	return `run:
go run ./cmd

run-grpc:
go run ./cmd/grpc

proto-gen:
# Requires protoc and the Go plugins (protoc-gen-go, protoc-gen-go-grpc)
protoc -I proto --go_out=paths=source_relative:./internal/api/grpc/pb --go-grpc_out=paths=source_relative:./internal/api/grpc/pb proto/*.proto

test:
go test ./... -v

lint:
golangci-lint run

migrate-up:
go run ./cmd/migrate_cmd.go up

migrate-down:
go run ./cmd/migrate_cmd.go down
`
}

func gitignoreTemplate() string {
	return `bin/
vendor/
.env
.env.*
*.log
*.out
.idea/
.vscode/
`
}

func dockerfileTemplate(service string) string {
	return fmt.Sprintf(`FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /%s main.go

FROM alpine:3.20
COPY --from=builder /%s /usr/local/bin/%s
CMD ["%s"]
`, service, service, service, service)
}

func dockerComposeTemplate(service string) string {
	return fmt.Sprintf(`version: "3.8"

services:
  %s:
build: .
container_name: %s
ports:
  - "8080:8080"
  - "9090:9090"
env_file:
  - .env
depends_on:
  - postgres
  - redis

  postgres:
image: postgres:15
environment:
  POSTGRES_USER: user
  POSTGRES_PASSWORD: password
  POSTGRES_DB: %s_db
ports:
  - "5432:5432"

  redis:
image: redis:7
ports:
  - "6379:6379"
`, service, service, service)
}

func readmeTemplate(service string) string {
	return fmt.Sprintf("# %s Service\n\nThis is the **%s microservice**, built using Clean Architecture principles in Go.\n\n## Structure\n\n- `cmd/` — CLI entrypoints\n- `cmd/grpc/` — example gRPC entrypoint\n- `proto/` — protobuf definitions\n- `internal/api/grpc/pb/` — generated gRPC Go code\n- `internal/` — core logic (domain, usecase, infra, api)\n- `pkg/` — reusable helper packages\n- `db/` — migrations\n- `docs/` — documentation\n\n## Commands\n\n```bash\nmake proto-gen\nmake run\nmake run-grpc\nmake test\nmake migrate-up\n```\n", strings.Title(service), service)
}

func gomodTemplate(service string) string {
	return fmt.Sprintf(`module github.com/tamirat-dejene/ha-soranu/%s
go 1.25
`, service)
}
