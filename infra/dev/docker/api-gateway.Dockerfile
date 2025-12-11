# Basic Dockerfile for auth-service
FROM alpine:3.23
WORKDIR /app

COPY shared /app/shared
COPY bin/api-gateway /app/bin/api-gateway

EXPOSE 8080
CMD ["/app/bin/api-gateway"]

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

# # Build api-gateway
# WORKDIR /app/services/api-gateway
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api-gateway ./cmd

# # ---- Runtime Stage ----
# FROM alpine:3.23

# WORKDIR /app

# COPY --from=builder /app/bin/api-gateway /app/bin/api-gateway

# EXPOSE 8080
# CMD ["/app/bin/api-gateway"]
