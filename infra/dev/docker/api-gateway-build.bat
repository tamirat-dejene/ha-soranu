#!/bin/bash
set -e

IMAGE_NAME="ha-soranu/api-gateway:dev"
DOCKERFILE_PATH="./infra/dev/docker/api-gateway.Dockerfile"
CONTEXT_PATH="./../../"
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

docker buildx build -t "$IMAGE_NAME" \
  -f "$DOCKERFILE_PATH" \
  --build-arg CGO_ENABLED="$CGO_ENABLED" \
  --build-arg GOOS="$GOOS" \
  --build-arg GOARCH="$GOARCH" \
  --output type=local,dest=./build-output \
  "$CONTEXT_PATH"
  
echo "Docker image '$IMAGE_NAME' built successfully and output saved to './build-output'."