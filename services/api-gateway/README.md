# API Gateway

## Overview

The API Gateway is an HTTP/REST service that serves as the entry point for the Ha-Soranu microservices application. It provides a unified RESTful API interface to clients while internally communicating with backend gRPC services. The gateway handles request/response transformation between HTTP/JSON and gRPC protocols, making the microservices easily accessible to web and mobile clients.

## Architecture

### Technology Stack

- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin) - High-performance HTTP router
- **Communication Protocol**: gRPC for backend service communication
- **Language**: Go 1.x
- **Logging**: Uber Zap (structured logging)

### Project Structure

```
api-gateway/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── dto/             # Data Transfer Objects
│   │   │   ├── auth_dto.go  # Authentication DTOs
│   │   │   ├── user_dto.go  # User management DTOs
│   │   │   └── err_dto.go   # Error response DTOs
│   │   └── handler/         # HTTP request handlers
│   │       ├── auth_handler.go
│   │       └── user_handler.go
│   ├── client/
│   │   └── auth_user_client.go  # gRPC client wrapper
│   ├── domain/
│   │   └── auth_domain.go   # Domain models
│   ├── errs/
│   │   └── errors.go        # Error definitions
│   └── server/
│       ├── server.go        # Server setup and routing
│       └── middleware.go    # Gin logger middleware
├── env.go                   # Environment configuration
└── README.md
```

### Key Components

#### 1. **Server** (`internal/server/server.go`)
- Initializes Gin router with recovery and logging middleware
- Configures CORS to allow cross-origin requests from all origins
- Sets up route groups and endpoints
- Manages server lifecycle

#### 2. **Handlers** (`internal/api/handler/`)
- **AuthHandler**: Manages authentication-related HTTP endpoints
  - User registration
  - Email/password login
  - Google OAuth login
  - Logout
  - Token refresh
- **UserHandler**: Manages user profile operations
  - Get user details
  - Manage phone numbers (add, update, remove)
  - Manage addresses (get, add, remove)

#### 3. **DTOs** (`internal/api/dto/`)
- Transform between HTTP JSON and gRPC protocol buffer messages
- Contains request and response data structures for:
  - Authentication operations
  - User operations
  - Error responses

#### 4. **gRPC Client** (`internal/client/auth_user_client.go`)
- `UAServiceClient` (User-Auth Service Client) connects to auth-service
- Maintains gRPC connections to `AuthService` and `UserService`
- Uses insecure credentials (suitable for internal microservice communication)

#### 5. **Middleware** (`internal/server/middleware.go`)
- Custom Gin logger that integrates with Zap structured logging
- Logs request details: method, path, status code, duration, client IP

## API Endpoints

### Base URL
```
http://localhost:8080
```

### Health & Info Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Welcome message |
| GET | `/health` | Health check endpoint |

### Authentication Endpoints (`/api/v1/auth`)

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/v1/auth/register` | Register new user | `{ "email", "password", "username", "phone_number" }` |
| POST | `/api/v1/auth/login` | Login with email/password | `{ "email", "password" }` |
| POST | `/api/v1/auth/google` | Login with Google OAuth | `{ "id_token" }` |
| POST | `/api/v1/auth/logout` | Logout user | `{ "refresh_token" }` |
| POST | `/api/v1/auth/refresh` | Refresh access token | `{ "refresh_token" }` |

### User Management Endpoints (`/api/v1/user`)

| Method | Endpoint | Description | Query Params |
|--------|----------|-------------|--------------|
| GET | `/api/v1/user/` | Get user details | `user_id` |
| GET | `/api/v1/user/phone-number` | Get phone number | `user_id` |
| POST | `/api/v1/user/phone-number` | Add phone number | Body: `{ "user_id", "phone_number" }` |
| PUT | `/api/v1/user/phone-number` | Update phone number | Body: `{ "user_id", "phone_number" }` |
| DELETE | `/api/v1/user/phone-number` | Remove phone number | Body: `{ "user_id" }` |
| GET | `/api/v1/user/addresses` | Get all addresses | `user_id` |
| POST | `/api/v1/user/addresses` | Add new address | Body: Address details |
| DELETE | `/api/v1/user/addresses` | Remove address | Body: `{ "user_id", "address_id" }` |

## Configuration

### Environment Variables

The API Gateway uses the following environment variables (defined in [`env.go`](file:///home/tamirat-dejene/Documents/dis-sys/ha-soranu/services/api-gateway/env.go)):

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `SRV_ENV` | Service environment (development/production) | `development` |
| `AUTH_SRV_NAME` | Hostname/service name of auth-service | `auth-service` |
| `AUTH_SRV_PORT` | Port of auth-service gRPC server | `9090` |
| `API_GATEWAY_PORT` | HTTP port for the API Gateway | `8080` |

### Configuration Loading

Configuration is loaded via the `GetEnv()` function which reads environment variables and provides sensible defaults. Helper functions:
- `getString(key, defaultValue)`: Get string environment variable
- `getInt(key, defaultValue)`: Get integer environment variable with validation

## Request/Response Flow

1. **Client → API Gateway**: Client sends HTTP/JSON request to a REST endpoint
2. **Gateway → Handler**: Gin router forwards request to appropriate handler
3. **Handler → DTO**: Handler binds JSON to DTO and validates
4. **DTO → Proto**: DTO converts to gRPC protobuf message
5. **Gateway → Backend**: gRPC client sends request to auth-service
6. **Backend → Gateway**: auth-service responds with gRPC message
7. **Proto → DTO**: Response converted from protobuf to DTO
8. **Gateway → Client**: Handler returns HTTP/JSON response

## Error Handling

### Error Response Format

All errors are returned in a consistent JSON format:

```json
{
  "error": "Error message description"
}
```

### Error Types

- **400 Bad Request**: Invalid request body or parameters
- **500 Internal Server Error**: Backend service errors or internal failures

The `ErrorResponseFromGRPCError()` function in DTOs translates gRPC errors to HTTP error responses.

## Dependencies

### Go Modules

Key dependencies include:

```go
require (
    github.com/gin-contrib/cors      // CORS middleware
    github.com/gin-gonic/gin         // HTTP web framework
    go.uber.org/zap                  // Structured logging
    google.golang.org/grpc           // gRPC client
)
```

### Internal Dependencies

- `shared/pkg/logger`: Centralized logger initialization
- `shared/protos/authpb`: Generated gRPC code for AuthService
- `shared/protos/userpb`: Generated gRPC code for UserService

## Running the Service

### Prerequisites

1. Go 1.x installed
2. Auth-service running on configured host:port
3. Environment variables set (optional, defaults provided)

### Local Development

```bash
# From project root
cd services/api-gateway

# Run the service
go run cmd/main.go
```

### Using Make (from project root)

```bash
make run-gateway
```

### Docker/Kubernetes

The service is containerized and deployed via Tilt for local Kubernetes development. See the root `Tiltfile` for deployment configuration.

## Logging

The API Gateway uses structured logging via Uber Zap:

- **Startup logs**: Service initialization and configuration
- **Request logs**: Each HTTP request (method, path, status, duration, IP)
- **Error logs**: Detailed error information with context
- **Connection logs**: gRPC client connection events

Logs are written to stdout in JSON format (production) or console format (development).

## CORS Configuration

CORS is configured to allow:
- **All origins**: `AllowAllOrigins: true`
- **All methods**: GET, POST, PUT, DELETE, OPTIONS, etc.
- **All headers**: Custom headers permitted

> [!WARNING]
> The current CORS configuration allows all origins for development convenience. In production, restrict `AllowOrigins` to specific trusted domains.

## Security Considerations

> [!CAUTION]
> **Development Mode**: The following security measures should be implemented for production:

1. **Authentication Middleware**: Currently, no authentication is enforced at the gateway level. Implement JWT validation middleware for protected endpoints.

2. **Rate Limiting**: Add rate limiting to prevent abuse.

3. **HTTPS/TLS**: Use TLS certificates for production deployments.

4. **CORS**: Restrict allowed origins to trusted domains.

5. **Input Validation**: Add comprehensive request validation beyond basic JSON binding.

6. **gRPC Security**: Use TLS credentials instead of `insecure.NewCredentials()` for gRPC connections.

## Future Enhancements

- [ ] JWT validation middleware at gateway level
- [ ] Request rate limiting
- [ ] API versioning strategy
- [ ] OpenAPI/Swagger documentation generation
- [ ] Circuit breaker pattern for gRPC calls
- [ ] Request/Response caching
- [ ] Metrics and monitoring (Prometheus)
- [ ] Distributed tracing (OpenTelemetry)

## Related Services

- **auth-service**: Backend gRPC service for authentication and user management
- **user-service**: (Future) Separate user management service

## Proto Definitions

The API Gateway interfaces with backend services using Protocol Buffers defined in:
- [`protos/auth.proto`](file:///home/tamirat-dejene/Documents/dis-sys/ha-soranu/protos/auth.proto) - AuthService RPCs
- [`protos/user.proto`](file:///home/tamirat-dejene/Documents/dis-sys/ha-soranu/protos/user.proto) - UserService RPCs

## Troubleshooting

### Common Issues

**Problem**: Cannot connect to auth-service  
**Solution**: Verify auth-service is running and `AUTH_SRV_NAME:AUTH_SRV_PORT` is correct

**Problem**: CORS errors in browser  
**Solution**: Verify CORS middleware is configured (should allow all origins by default)

**Problem**: 400 Bad Request errors  
**Solution**: Check request body format matches expected DTO structure

**Problem**: 500 Internal Server Error  
**Solution**: Check API Gateway logs and auth-service logs for detailed error information

## Contributing

When adding new endpoints:
1. Define DTOs in `internal/api/dto/`
2. Create/update handler in `internal/api/handler/`
3. Add route in `server.SetupRoutes()`
4. Update this README with endpoint documentation

## License

Part of the Ha-Soranu microservices platform.