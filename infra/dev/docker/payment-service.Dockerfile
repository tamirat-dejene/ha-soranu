FROM alpine:3.23
WORKDIR /app

COPY shared /app/shared
COPY services/payment-service/migrations /app/payment-svc-migrations
COPY bin/payment-service /app/bin/payment-service

EXPOSE 9090 8081
CMD ["/app/bin/payment-service"]
