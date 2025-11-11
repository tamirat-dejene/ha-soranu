FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /notification-service main.go

FROM alpine:3.20
COPY --from=builder /notification-service /usr/local/bin/notification-service
CMD ["notification-service"]
