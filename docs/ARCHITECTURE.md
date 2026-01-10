# System Architecture

## Overview

Ha-Soranu is a microservices-based food delivery platform designed for scalability and asynchronous processing. It leverages gRPC for internal low-latency communication and Apache Kafka for event-driven workflows, ensuring loose coupling between services.

![Architecture Diagram](./hasoranu.png)

## Core Components

### 1. API Gateway (`services/api-gateway`)

- **Type**: REST API (HTTP/JSON)
- **Role**: Entry point for all client applications (Web, Mobile).
- **Responsibilities**:
  - Routing requests to appropriate backend microservices.
  - Protocol translation (HTTP to gRPC).
  - Request validation and rate limiting.
  - Aggregating responses from multiple services.

### 2. Auth Service (`services/auth-service`)

- **Type**: gRPC Service
- **Role**: Identity and Access Management.
- **Responsibilities**:
  - User registration and login.
  - JWT token issuance and validation.
  - Managing user profiles and addresses.
  - Storing sensitive data securely.

### 3. Restaurant Service (`services/restaurant-service`)

- **Type**: gRPC Service
- **Role**: Catalog and Order Management.
- **Responsibilities**:
  - Managing restaurant profiles and menus.
  - Handling order creation and lifecycle (Placed, Cooking, Ready, Delivered).
  - Validating item availability and prices.
  - Emitting `OrderCreated` events.

### 4. Payment Service (`services/payment-service`)

- **Type**: HTTP/Worker Service
- **Role**: Financial Transaction Processing.
- **Responsibilities**:
  - Handling payment intents and captures.
  - Integrating with third-party payment gateways (e.g., Stripe, PayPal).
  - Managing refunds.
  - Emitting `PaymentProcessed` events to trigger order confirmation.

### 5. Notification Service (`services/notification-service`)

- **Type**: gRPC/Worker Service
- **Role**: User Communication.
- **Responsibilities**:
  - Listening to system events (e.g., `OrderStatusChanged`, `PaymentFailed`).
  - Sending emails, SMS, or push notifications to users.
  - Notifying restaurants of new orders.

## Communication Patterns

### Synchronous (gRPC)

Used for critical, real-time operations where an immediate response is required.

- **API Gateway → Auth Service**: Validate tokens, get user details.
- **API Gateway → Restaurant Service**: Browse menu, create order.

### Asynchronous (Kafka)

Used for side effects and decoupling long-running processes.

- **Restaurant Service → Kafka**: Publishes `OrderCreated`.
- **Payment Service → Kafka**: Publishes `PaymentSuccess`/`PaymentFailed`.
- **Notification Service**: Consumes events to send updates.

## Data Consistency

The system follows the **Database-per-Service** pattern. Each service manages its own database schema and migrations, ensuring strict data encapsulation. Cross-service data consistency is achieved through eventual consistency using Kafka events.
