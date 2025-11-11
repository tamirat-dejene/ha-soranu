# Use a minimal base image
FROM alpine:3.22

# Add necessary runtime dependencies (if your Go binary uses net/http, DNS, etc.)
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

# Set working directory
WORKDIR /app

# Copy in the Go binary and related directories
COPY ./bin/api-gateway .
COPY shared ./shared
COPY build ./build

# Set non-root user for security
RUN adduser -D -g '' appuser
USER appuser

# Expose HTTP and gRPC ports
EXPOSE 3030 50051

# Command to run the service
ENTRYPOINT ["./api-gateway"]
