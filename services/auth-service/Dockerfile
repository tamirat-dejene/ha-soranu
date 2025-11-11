FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /auth-service main.go

FROM alpine:3.20
COPY --from=builder /auth-service /usr/local/bin/auth-service
CMD ["auth-service"]
