FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /api-gateway main.go

FROM alpine:3.20
COPY --from=builder /api-gateway /usr/local/bin/api-gateway
CMD ["api-gateway"]
