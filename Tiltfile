# ============================================
# Tiltfile — Local Dev Environment
# ============================================

# Extensions
load('ext://restart_process', 'docker_build_with_restart')

# --- Local Builds ---
# Build binaries locally to ./bin/ directory
# local_resource(
#     'build-api-gateway',
#     cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api-gateway ./services/api-gateway/cmd',
#     deps=['services/api-gateway', 'shared', 'go.mod', 'go.sum'],
#     ignore=['**/tmp', '**/.git'],
#     labels="BUILD_ONLY",
# )

# local_resource(
#     'build-auth-service',
#     cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/auth-service ./services/auth-service/cmd',
#     deps=['services/auth-service', 'shared', 'go.mod', 'go.sum'],
#     ignore=['**/tmp', '**/.git'],
#     labels="BUILD_ONLY",
# )

# --- Docker Builds with Live Update ---
docker_build_with_restart(
    'ha-soranu/auth-service',
    context='.',  # repo root
    dockerfile='./infra/dev/docker/auth-service.Dockerfile',
    entrypoint=['/app/bin/auth-service'],
    live_update=[
        sync('./services/auth-service', '/app/services/auth-service'),
        sync('./shared', '/app/shared'),
    ],
)

docker_build_with_restart(
    'ha-soranu/api-gateway',
    context='.',
    dockerfile='./infra/dev/docker/api-gateway.Dockerfile',
    entrypoint=['/app/bin/api-gateway'],
    live_update=[
        sync('./services/api-gateway', '/app/services/api-gateway'),
        sync('./shared', '/app/shared'),
    ],
)


# --- Kubernetes Resources ---
k8s_yaml([
    'infra/dev/k8s/api-gateway-deployment.yaml',
    'infra/dev/k8s/auth-service-deployment.yaml',
    'infra/dev/k8s/postgres-deployment.yaml',
    'infra/dev/k8s/redis-deployment.yaml',
    'infra/dev/k8s/config-map.yaml',
    'infra/dev/k8s/secrets.yaml',
])

# --- Port Forwards ---
# Expose services to localhost
k8s_resource('api-gateway', port_forwards=['8080:8080'])
k8s_resource('auth-service', port_forwards=['9090:9090'])
k8s_resource('postgres-db', port_forwards=['5432:5432'])
k8s_resource('redis', port_forwards=['6379:6379'])

# --- End of File ---
print("Tiltfile loaded successfully — monitoring api-gateway, and auth-service.")
