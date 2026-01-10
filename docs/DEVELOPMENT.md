# Development Guide

This guide will help you set up your local development environment for Ha-Soranu.

## Prerequisites

Ensure you have the following installed:

- **Go 1.25+**: [Download](https://go.dev/dl/)
- **Docker Desktop**: [Download](https://www.docker.com/products/docker-desktop)
- **Kubernetes Cluster**: Minikube, Kind, or Docker Desktop's built-in K8s.
- **Tilt**: [Download](https://docs.tilt.dev/install.html) - Orchestrates the dev environment.
- **Protoc**: Protocol buffers compiler.
- **Go Plugins for Protoc**:
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
- **Goose**: Database migration tool.
    ```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

## Quick Start (with Tilt)

Tilt is the recommended way to run the full stack locally. It builds your services, deploys them to Kubernetes, and syncs changes live.

1. **Start Kubernetes**: Ensure your local cluster is running (e.g., `minikube start`).
2. **Launch Tilt**:
   ```bash
   tilt up
   ```
3. **Open Dashboard**: Press `Space` to open the Tilt web UI. You should see all services (`api-gateway`, `auth-service`, etc.) spinning up.

## Running Tests

### End-to-End / Load Tests
Located in the `tests/` directory.

1. **Setup**:
   ```bash
   cd tests
   make setup
   source venv/bin/activate
   ```
2. **Seed Data**:
   ```bash
   make seed
   ```
3. **Run Tests**:
   ```bash
   make run    # Interactive (with Locust UI)
   # OR
   ./run_tests.sh all # Headless
   ```

## Development Workflow

1. **Modify Code**: Edit a file in `services/`.
2. **Automatic Sync**: Tilt detects the change, rebuilds the binary (if needed), and live-updates the container in seconds.
3. **Modify Protos**:
   - Edit `.proto` files in `protos/`.
   - Run `make proto` from the root directory.
   - Update your Go code to implement new interfaces.
