# App

Containerized apps deployed to k3s on Raspberry Pi 5 (`elxr-rpi5`).


## Structure

```
app/
├── Makefile           # Build, load, deploy commands
├── api/               # Go REST API (counter, fake LLM with SSE, swagger)
└── deploy/
    ├── base/          # Shared k8s manifests (kustomize)
    └── overlays/rpi/  # Raspberry Pi overlay
```

## Prerequisites

- Docker with buildx (for ARM64 cross-compilation)
- SSH access to `root@elxr-rpi5`
- k3s running on the Pi

## Usage

```bash
make build          # Build ARM64 Docker images
make deploy         # Build + load images to Pi + apply k8s manifests + restart pods
make clean          # Remove local images
```

## Apps

| App | Port | Description |
|-----|------|-------------|
| api | 30080 | Go REST API with counter, fake LLM streaming, and Swagger UI |

Access after deploy: `http://elxr-rpi5:30080`

## Adding a new app

1. Create `app/<name>/` with a `Dockerfile`
2. Add k8s manifests to `deploy/base/`
3. Add the app name to `APPS` in `Makefile`
