FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /delivery-service main.go

FROM alpine:3.20
COPY --from=builder /delivery-service /usr/local/bin/delivery-service
CMD ["delivery-service"]
