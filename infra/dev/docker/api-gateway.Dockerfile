FROM alpine:3.23
WORKDIR /app

COPY shared /app/shared
COPY bin/api-gateway /app/bin/api-gateway

EXPOSE 8080
CMD ["/app/bin/api-gateway"]