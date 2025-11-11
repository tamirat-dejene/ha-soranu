FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /payment-service main.go

FROM alpine:3.20
COPY --from=builder /payment-service /usr/local/bin/payment-service
CMD ["payment-service"]
