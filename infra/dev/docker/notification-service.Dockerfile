FROM alpine:3.23
WORKDIR /app

COPY ./bin/notification-service /app/bin/notification-service
COPY ./services/notification-service/migrations /app/notification-svc-migrations

EXPOSE 50053
CMD ["/app/bin/notification-service"]