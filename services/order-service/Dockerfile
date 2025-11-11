FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /order-service main.go

FROM alpine:3.20
COPY --from=builder /order-service /usr/local/bin/order-service
CMD ["order-service"]
