FROM alpine:3.23
WORKDIR /app

COPY shared /app/shared
COPY services/restaurant-service/migrations /app/restaurant-svc-migrations
COPY bin/restaurant-service /app/bin/restaurant-service

EXPOSE 9091
CMD ["/app/bin/restaurant-service"]
