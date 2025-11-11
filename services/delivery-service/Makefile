run:
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
