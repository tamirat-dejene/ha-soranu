# Ha-Soranu 🍱

**Ha-Soranu** is a scalable, microservices-based online food delivery platform built with Go. It enables users to browse restaurants, view menus, and place orders through a modern, distributed architecture.

## Architecture

The system follows a microservices architecture using **gRPC** for inter-service communication and an **API Gateway** to expose RESTful endpoints to clients.

![Ha-Soranu architecture diagram](./docs/hasoranu.png)

*Figure: High-level architecture — API Gateway, Auth, Restaurant, Payment, and Notification services, Postgres, Redis, Kafka.*

### Core Services

| Service | Type | Description |
|---------|------|-------------|
| **[api-gateway](./services/api-gateway)** | REST API | Entry point for all client requests. Handles routing, transformation (HTTP ↔ gRPC), and initial validation. |
| **[auth-service](./services/auth-service)** | gRPC | Manages user identity, authentication (JWT, OAuth), and profiles. |
| **[restaurant-service](./services/restaurant-service)** | gRPC | Manages restaurant profiles, menus, and order processing. Publishes order events to Kafka. |
| **[notification-service](./services/notification-service)** | gRPC | Handles real-time notifications for users and restaurants. Consumes order events from Kafka. |
| **[payment-service](./services/payment-service)** | gRPC | Processes payments and refunds. Integrates with payment providers, records transactions, and emits payment events.

## Tech Stack

- **Language**: Go (Golang)
- **Communication**: gRPC (Inter-service), REST (Client-facing)
- **Frameworks**: 
  - [Gin](https://github.com/gin-gonic/gin) (HTTP Gateway)
  - [gRPC-Go](https://github.com/grpc/grpc-go) (RPC)
- **Databases**: PostgreSQL (Per-service databases)
- **Caching**: Redis (Token storage, session management)
- **Messaging**: Apache Kafka (Event-driven architecture with binary Protobuf serialization)
- **Infrastructure**: 
  - Docker & Kubernetes (Containerization & Orchestration)
  - [Tilt](https://tilt.dev) (Local Development Environment)
- **Tooling**: Make, Goose (Migrations)

## Getting Started

The easiest way to run the entire platform locally is using **Tilt**.

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://www.docker.com/)
- [Kubernetes Cluster](https://kubernetes.io/) (Minikube)
- [Tilt](https://docs.tilt.dev/install.html)
- [Make](https://www.gnu.org/software/make/)

### Running Locally

1. **Start your Kubernetes cluster** (if not running):
   ```bash
   minikube start
   ```

2. **Run Tilt**:
   ```bash
   tilt up
   ```
   This will build all services, deploy them to Kubernetes, and set up port forwarding.

3. **Access the Application**:
   - **API Gateway**: `http://localhost:8080`
   - **Tilt UI**: `http://localhost:10350` (to monitor logs and status)

## Project Structure

```bash
ha-soranu/
├── services/               # Microservices source code
│   ├── api-gateway/        # REST API Gateway
│   ├── auth-service/       # Authentication & User Service
│   ├── restaurant-service/ # Restaurant & Menu Service
│   ├── notification-service/ # Notification Service
│   └── payment-service/    # Payment & Transaction Service
├── protos/                 # Protocol Buffer definitions (gRPC contracts)
├── shared/                 # Shared libraries (DB packages, Logger, Middleware, Events)
├── infra/                  # Infrastructure configurations (K8s, Docker)
├── bin/                    # Compiled binaries (ignored by git)
├── Tiltfile                # Tilt configuration for local dev
├── Makefile                # Global build and run commands
└── go.mod                  # Go module definition (Workspace mode)
```


## Key Features

- **Authentication System**:
  - Secure Email/Password login.
  - **Google OAuth 2.0** integration.
  - JWT-based session management with Access & Refresh tokens.
  
- **Restaurant Management**:
  - Restaurant registration and profile management.
  - Menu CRUD operations (Items, prices, descriptions).
  - Geospatial discovery (Streaming API).
  
- **Order Processing**:
  - Order placement and status tracking.
  - Real-time order updates.
  - **Event-Driven**: Asynchronous order processing using Kafka with binary Protobuf serialization.
  
- **Notification System**:
  - Real-time notifications for users and restaurants.
  - Order placement notifications for restaurants.
  - Order status update notifications for customers.
  - **Event-Driven**: Kafka consumer processes order events to create notifications.

- **Payment System**:
  - Payment intent creation, capture, and refund flows.
  - Integration with payment providers (e.g., Stripe-like gateways).
  - Secure tokenized processing, idempotent operations, and audit logs.
  - **Event-Driven**: Emits payment events and reacts to order lifecycle updates.

## Development

### Working with Protos

Protocol buffers are the contract between services. If you modify files in `protos/`, you must regenerate the Go code:

```bash
make proto
```

### Database Migrations

Each service manages its own database schema using **Goose**.

```bash
# Example: Create a new migration for auth-service
cd services/auth-service
goose -dir migrations postgres "user=postgres dbname=authdb sslmode=disable" create add_new_table sql
```