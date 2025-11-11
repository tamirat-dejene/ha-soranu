FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /review-service main.go

FROM alpine:3.20
COPY --from=builder /review-service /usr/local/bin/review-service
CMD ["review-service"]
