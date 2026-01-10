# API Reference

Ha-Soranu exposes its functionality through a REST API for external clients and uses gRPC for high-performance internal communication.

## REST API (External)

The API Gateway acts as the single entry point for all external traffic. It exposes a RESTful interface compliant with OpenAPI 3.0.

### Specification
The complete OpenAPI specification (Swagger) can be found at:
- **Location**: `services/api-gateway/docs/swagger.yaml`
- **Usage**: You can import this file into [Swagger Editor](https://editor.swagger.io/) or Postman to explore endpoints and generate clients.

### Key Endpoints

#### Authentication
- `POST /api/v1/user/auth/register`: Register a new user.
- `POST /api/v1/user/auth/login`: Login and receive JWT.

#### Restaurants
- `GET /api/v1/restaurant/restaurants`: List restaurants.
- `GET /api/v1/restaurant/restaurants/{id}`: Get restaurant details and menu.

#### Orders
- `POST /api/v1/order/orders`: Place a new order.
- `GET /api/v1/order/orders/{id}`: Get order status.

---

## gRPC API (Internal)

Services communicate internally using Protocol Buffers (Protobuf) over gRPC.

### Definition Files
All `.proto` files are located in the `protos/` directory.

| Service | Proto File | Package | Description |
|---|---|---|---|
| **Auth** | [`protos/auth.proto`](../protos/auth.proto) | `auth` | RPCs for login, signup, and token validation. |
| **User** | [`protos/user.proto`](../protos/user.proto) | `user` | User profile data structures and RPCs. |
| **Restaurant**| [`protos/restaurant.proto`](../protos/restaurant.proto) | `restaurant` | Restaurant search, menu retrieval, order management. |
| **Notification** | [`protos/notification.proto`](../protos/notification.proto) | `notification` | Unary and streaming RPCs for notifications. |
| **Order** | [`protos/order.proto`](../protos/order.proto) | `order` | Order data structures. |

### Code Generation
Go stubs are generated using the `Makefile` command:
```bash
make proto
```
This generates code into the `shared/protos` directory.

---

## Event API (Kafka)

Asynchronous communication is handled via Apache Kafka. Messages are serialized using Protobuf for type safety and performance.

### Event Envelope
All events are wrapped in a standard envelope defined in `protos/event_envelope.proto` to ensure consistent metadata (ID, timestamp, type).

### Key Events

| Event Type | Producer | Consumers | Payload Proto | Description |
|---|---|---|---|---|
| `order.created` | Restaurant Service | Payment, Notification | `order.Order` | Emitted when a user places an order. |
| `payment.processed` | Payment Service | Restaurant, Notification | `payment.PaymentEvent` | Emitted after payment attempt (success/failure). |
| `order.status_changed` | Restaurant Service | Notification | `order.OrderStatus` | Emitted when order moves to cooking, ready, etc. |
