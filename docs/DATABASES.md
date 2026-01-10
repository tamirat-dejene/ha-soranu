# Database & Persistence

Ha-Soranu uses a **Database-per-Service** pattern to ensure loose coupling and independent scaling.

## Technologies

### PostgreSQL
Used as the primary persistent store for structured data.
- **Benefits**: Relational data integrity, ACID transactions, robust ecosystem.
- **Connections**: Services connect using the `pgx` driver.

### Redis / Valkey
Used for ephemeral data and high-speed caching.
- **Use Cases**:
    - **Caching**: Storing restaurant menus or user sessions to reduce DB load.
    - **Tokens**: Storing revocable JWTs or refresh tokens.

## Schema Management
We use [Pressly Goose](https://github.com/pressly/goose) for database migrations. This allows us to version control our database changes.

### Migration Structure
Migrations are located in each service's directory under `migrations/`:
```
services/
  ├── auth-service/
  │   └── migrations/
  │       ├── 20231026120000_create_users.sql
  │       └── ...
  ├── restaurant-service/
  │   └── migrations/
  │       └── ...
```

### Running Migrations

#### Automatically (Dev/Prod)
Services are configured to embed migration files and run them automatically on purely "up" migrations at startup in many environments to ensure the schema is always current.

#### Manually (Development)
You can run migrations manually using the `goose` CLI if you need to roll back or test specific states.

1. **Install Goose**:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. **Run Command** (example for auth-service):
   ```bash
   cd services/auth-service
   # export DB credentials first!
   export GOOSE_DRIVER=postgres
   export GOOSE_DBSTRING="user=postgres password=password dbname=authdb sslmode=disable"
   goose -dir migrations status
   goose -dir migrations up
   ```

## Service Data Owners

| Service | Database | Key Tables |
|---|---|---|
| **Auth** | `authdb` | `users`, `addresses` |
| **Restaurant** | `restaurantdb` | `restaurants`, `menus`, `items` |
| **Order** | `orderdb` | `orders`, `order_items` |
| **Payment**| `paymentdb` | `payments`, `refunds` |
