package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var structure = []string{
	"cmd/http",
	"cmd/grpc",
	"db/migrations",
	"docs",
	"proto",

	"internal/api/http",
	"internal/api/grpc/pb",
	"internal/config",
	"internal/domain",
	"internal/infra/db",
	"internal/infra/redis",
	"internal/server",
	"internal/usecase",
	"internal/util",

	"shared/pkg/logger",
	"shared/pkg/errors",
	"test",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run newservice.go <service-name>")
		os.Exit(1)
	}

	service := strings.TrimSpace(os.Args[1])
	if service == "" {
		fmt.Println("Invalid service name")
		os.Exit(1)
	}

	root := filepath.Join("services", service)

	// Prevent overwriting existing service
	if _, err := os.Stat(root); err == nil {
		fmt.Printf("Service '%s' already exists. Aborting.\n", service)
		os.Exit(1)
	}

	fmt.Printf("Creating new microservice: %s\n\n", service)

	// Create folder structure
	for _, dir := range structure {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0o755); err != nil {
			fmt.Printf("Failed to create %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("Created: %s\n", path)
	}

	// Create boilerplate entrypoints
	createFile(filepath.Join(root, "cmd", "http", "main.go"), httpMainTemplate(service))
	createFile(filepath.Join(root, "cmd", "grpc", "main.go"), grpcMainTemplate(service))
	createFile(filepath.Join(root, "Makefile"), makefileTemplate())
	createFile(filepath.Join(root, "README.md"), readmeTemplate(service))

	fmt.Println("\nDone! Service skeleton created.")
	fmt.Printf("Next steps:\n  cd services/%s\n", service)
}

func createFile(path, content string) {
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		fmt.Printf("File exists, skipping: %s\n", path)
		return
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Printf("Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("Created file: %s\n", path)
}

func httpMainTemplate(service string) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting %s HTTP service on :8080...")

	// TODO: Load config, setup logger, DI, routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server failed: %%v", err)
	}
}
`, service)
}

func grpcMainTemplate(service string) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting %s gRPC service on :9090...")

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %%v", err)
	}

	s := grpc.NewServer()

	// TODO: Register protobuf services
	// pb.RegisterYourServiceServer(s, &handler.YourService{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %%v", err)
	}
}
`, service)
}

func makefileTemplate() string {
	return `run-http:
	go run ./cmd/http

run-grpc:
	go run ./cmd/grpc

proto-gen:
	protoc -I proto \
		--go_out=paths=source_relative:./internal/api/grpc/pb \
		--go-grpc_out=paths=source_relative:./internal/api/grpc/pb \
		proto/*.proto

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

func readmeTemplate(service string) string {
	title := cases.Title(language.Und).String(service)

	return fmt.Sprintf(
		"# %s Service\n\n"+
			"This is the **%s microservice**, generated using a Clean Architecture skeleton.\n\n"+
			"## Structure\n\n"+
			"- `cmd/` – entrypoints (HTTP/gRPC)\n"+
			"- `internal/` – core business logic\n"+
			"- `proto/` – protobuf definitions\n"+
			"- `db/migrations/` – SQL migrations\n"+
			"- `docs/` – documentation\n"+
			"- `shared/` – reusable packages\n\n"+
			"## Commands\n\n"+
			"```bash\n"+
			"make proto-gen\n"+
			"make run-http\n"+
			"make run-grpc\n"+
			"make test\n"+
			"make migrate-up\n"+
			"```\n",
		title, service,
	)
}
