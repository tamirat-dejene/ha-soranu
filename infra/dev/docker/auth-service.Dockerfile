FROM alpine:3.23
WORKDIR /app

COPY shared /app/shared
COPY bin/auth-service /app/bin/auth-service

EXPOSE 9090
CMD ["/app/bin/auth-service"]
