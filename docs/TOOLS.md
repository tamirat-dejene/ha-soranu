# Tools & Utilities

## CLI Tools

### `make`
The project uses a root `Makefile` to handle common tasks, primarily Protobuf generation.

- **`make proto`**: Scans the `protos/` directory and generates Go gRPC stubs into `shared/protos/`. It relies on `protoc` and the configured Go plugins.

### `protoc`
The Protocol Buffers compiler. Used to generate code from `.proto` definitions.
- **Version**: Ensure you are using a recent version of `protoc` (v3+).

### `goose`
Database migration tool.
- Used to create and run SQL migrations.
- Each service manages its own migrations in a `migrations/` subdirectory.

### `tilt`
Development environment orchestrator.
- **`Tiltfile`**: The configuration file at the root of the repo. It defines how to build Docker images (live update) and apply Kubernetes manifests.

### `kubectl`
Kubernetes command-line tool.
- Used to inspect the state of the cluster, view logs (if not using Tilt), and manage resources.

## Scripts

### `tests/run_tests.sh`
A helper script to run load tests and integration scenarios headlessly.

### `tests/Makefile`
Manages the Python virtual environment and dependencies for the test suite (Locust).
