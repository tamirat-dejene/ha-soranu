# Deployment & Infrastructure

## Overview
Ha-Soranu is designed to run on **Kubernetes**, utilizing **Docker** containers for services.

## Infrastructure Components

### Docker
Each microservice has a corresponding `Dockerfile` in `infra/dev/docker/`.
- Multi-stage builds are used to create lightweight, production-ready images (using `alpine` or `scratch` where possible).
- Common base images help maintain consistency.

### Kubernetes
The application is orchestrated using Kubernetes manifests located in `infra/dev/k8s/`.
- **Deployments**: Define the desired state for pods (replicas, images, env vars).
- **Services**: Expose deployments internally (ClusterIP) or externally.
- **ConfigMaps/Secrets**: Manage configuration and sensitive data.

### Tilt (Local Development)
We use [Tilt](https://tilt.dev/) for local development orchestration. It watches for file changes, rebuilds containers incrementally, and updates the local Kubernetes cluster in real-time.

## Configuration

Services are configured primarily via **Environment Variables**.

### Common Variables
| Variable | Description |
|---|---|
| `PORT` | The port the service listens on (gRPC or HTTP). |
| `POSTGRES_HOST` | Database hostname. |
| `POSTGRES_PORT` | Database port. |
| `POSTGRES_USER` | Database username. |
| `POSTGRES_PASSWORD` | Database password. |
| `REDIS_ADDR` | Redis connection address. |
| `KAFKA_BROKERS` | Comma-separated list of Kafka brokers. |

## Production Deployment Guide

1. **Cluster Setup**: Ensure you have a running Kubernetes cluster (EKS, GKE, DigitalOcean, etc.).
2. **Dependencies**:
   - Deploy **Postgres** and **Redis** (managed services recommended).
   - Deploy **Kafka** (managed or via Strimzi operator).
3. **Secrets**:
   - Create Kubernetes Secrets for database passwords, JWT keys, and API keys.
   - *Do not commit `secrets.yaml` to version control!*
4. **Build Images**:
   - Build and tag Docker images for each service.
   - Push to a container registry (Docker Hub, BCR, ECR).
5. **Apply Manifests**:
   - Update image tags in deployment manifests.
   - Run `kubectl apply -f infra/prod/k8s/`.

## CI/CD Pipeline (Recommended)
- **CI**: Run tests and linting on every commit. Build Docker images on merge.
- **CD**: Use GitOps (e.g., ArgoCD) to sync the `infra` directory with the cluster.
