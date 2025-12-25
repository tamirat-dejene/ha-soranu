# Basic Dockerfile for auth-service
FROM alpine:3.23
WORKDIR /app

COPY shared /app/shared
COPY services/auth-service/migrations /app/auth-svc-migrations
COPY bin/auth-service /app/bin/auth-service

EXPOSE 50051
CMD ["/app/bin/auth-service"]

# Advanced Dockerfile for auth-service with multi-stage build
# # ---- Build Stage ----
# FROM golang:1.25-alpine AS builder

# WORKDIR /app

# # Copy root module for caching
# COPY go.mod go.sum ./

# # RUN apk add --no-cache git
# RUN go mod download

# # Copy full repo
# COPY . .

# # Build auth-service
# WORKDIR /app/services/auth-service
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/auth-service ./cmd

# # ---- Runtime Stage ----
# FROM alpine:3.23

# WORKDIR /app

# COPY --from=builder /app/bin/auth-service /app/bin/auth-service

# # Copy migrations folder
# COPY --from=builder /app/services/auth-service/migrations /app/migrations

# EXPOSE 50051
# CMD ["/app/bin/auth-service"]
