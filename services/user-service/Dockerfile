FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /user-service main.go

FROM alpine:3.20
COPY --from=builder /user-service /usr/local/bin/user-service
CMD ["user-service"]
