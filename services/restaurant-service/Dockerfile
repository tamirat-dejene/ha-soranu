FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /restaurant-service main.go

FROM alpine:3.20
COPY --from=builder /restaurant-service /usr/local/bin/restaurant-service
CMD ["restaurant-service"]
