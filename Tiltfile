# ============================================
# Tiltfile — Local Dev Environment
# ============================================

# Extensions
load('ext://restart_process', 'docker_build_with_restart')

# --- Local Builds ---
# Build binaries locally to ./bin/ directory
local_resource(
    'build-api-gateway',
    cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api-gateway ./services/api-gateway/cmd',
    deps=['services/api-gateway', 'shared', 'go.mod', 'go.sum'],
    ignore=['**/tmp', '**/.git'],
    labels="BUILD_ONLY",
)

local_resource(
    'build-auth-service',
    cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/auth-service ./services/auth-service/cmd',
    deps=['services/auth-service', 'shared', 'go.mod', 'go.sum'],
    ignore=['**/tmp', '**/.git'],
    labels="BUILD_ONLY",
)

local_resource(
    'build-restaurant-service',
    cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/restaurant-service ./services/restaurant-service/cmd',
    deps=['services/restaurant-service', 'shared', 'go.mod', 'go.sum'],
    ignore=['**/tmp', '**/.git'],
    labels="BUILD_ONLY",
)

local_resource(
    'build-notification-service',
    cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/notification-service ./services/notification-service/cmd',
    deps=['services/notification-service', 'shared', 'go.mod', 'go.sum'],
    ignore=['**/tmp', '**/.git'],
    labels="BUILD_ONLY",
)

# --- Docker Builds with Live Update ---
docker_build_with_restart(
    'ha-soranu/auth-service',
    context='.',
    dockerfile='./infra/dev/docker/auth-service.Dockerfile',
    entrypoint=['./bin/auth-service'],
    only=["./bin/auth-service", "./shared/", "./services/auth-service/migrations/"],
    live_update=[
        sync("./bin/auth-service", "/app/bin/auth-service"),
        sync("./shared", "/app/shared"),
    ],
)

docker_build_with_restart(
    'ha-soranu/api-gateway',
    context='.',
    dockerfile='./infra/dev/docker/api-gateway.Dockerfile',
    entrypoint=['./bin/api-gateway'],
    only=["./bin/api-gateway", "./shared/"],
    live_update=[
        sync("./bin/api-gateway", "/app/bin/api-gateway"),
        sync("./shared", "/app/shared"),
    ],
)

docker_build_with_restart(
    'ha-soranu/restaurant-service',
    context='.',
    dockerfile='./infra/dev/docker/restaurant-service.Dockerfile',
    entrypoint=['./bin/restaurant-service'],
    only=["./bin/restaurant-service", "./shared/", "./services/restaurant-service/migrations/"],
    live_update=[
        sync("./bin/restaurant-service", "/app/bin/restaurant-service"),
        sync("./shared", "/app/shared"),
    ],
)

docker_build_with_restart(
    'ha-soranu/notification-service',
    context='.',
    dockerfile='./infra/dev/docker/notification-service.Dockerfile',
    entrypoint=['./bin/notification-service'],
    only=["./bin/notification-service", "./shared/", "./services/notification-service/migrations/"],
    live_update=[
        sync("./bin/notification-service", "/app/bin/notification-service"),
        sync("./shared", "/app/shared"),
    ],
)

# --- Kubernetes Resources ---
k8s_yaml([
    'infra/dev/k8s/ha-soranu-namespace.yaml',
    'infra/dev/k8s/api-gateway-deployment.yaml',
    'infra/dev/k8s/auth-service-deployment.yaml',
    'infra/dev/k8s/restaurant-service-deployment.yaml',
    'infra/dev/k8s/notification-service-deployment.yaml',
    'infra/dev/k8s/config-map.yaml',
    'infra/dev/k8s/secrets.yaml',
])
 
# --- Port Forwards ---
# Expose services to localhost
k8s_resource('api-gateway', port_forwards=['8080:8080'], labels="MONITORED")
k8s_resource('auth-service', port_forwards=['50051:50051'], labels="MONITORED")
k8s_resource('restaurant-service', port_forwards=['50052:50052'] , labels="MONITORED")
k8s_resource('notification-service', port_forwards=['50053:50053'], labels="MONITORED")

# --- End of File ---
print("Tiltfile loaded successfully — monitoring api-gateway, auth-service, restaurant-service, and notification-service.")