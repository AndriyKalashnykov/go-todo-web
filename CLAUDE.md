# Go Todo Web

## Project Overview

Go Todo Web is an HTTP REST API for managing todo items, built with the [Echo](https://echo.labstack.com/) framework (v5) in Go. It supports CRUD operations and is containerized for Kubernetes deployment.

## Tech Stack

- **Language**: Go 1.26.1
- **Framework**: Echo v5
- **Container**: Docker (multi-arch: amd64, arm64)
- **Registry**: GitHub Container Registry (ghcr.io)
- **CI**: GitHub Actions
- **Linting**: golangci-lint v2, hadolint
- **Dependency Management**: Renovate

## Project Structure

```
main.go              # Application entrypoint (NewApp setup + server start)
main_test.go         # App setup and route registration tests
handlers/            # HTTP request handlers (CRUD operations)
handlers/*_test.go   # HTTP integration tests via Echo router
models/              # Data models (Todo) and in-memory store
models/*_test.go     # Unit tests for all CRUD operations
Makefile             # Build orchestration and CI pipeline
Dockerfile           # Multi-stage Docker build
go.mod               # Go module definition
go.sum               # Go module checksums
renovate.json        # Renovate dependency update configuration
go-todo-web.yaml     # Kubernetes deployment manifest
version.txt          # Current version tag
.github/workflows/   # GitHub Actions CI/CD workflows
```

## Common Commands

```bash
make help            # Show all available targets
make deps            # Install and verify required tool dependencies
make deps-check      # Show required Go versions and gvm status
make build           # Build the Go binary
make run             # Run the application locally (port 8080)
make test            # Run tests with coverage and race detection
make lint            # Run golangci-lint and hadolint
make format          # Format Go source files
make static-check    # Run all static analysis checks
make coverage-check  # Verify test coverage meets minimum threshold
make ci              # Run full local CI pipeline
make ci-run          # Run GitHub Actions workflow locally using act
```

## API Endpoints

| Method | Path           | Description       |
|--------|----------------|-------------------|
| POST   | `/create`      | Create a todo     |
| GET    | `/get/:id`     | Get a todo by ID  |
| GET    | `/all`         | List all todos    |
| DELETE | `/delete/:id`  | Delete a todo     |
| PATCH  | `/update/:id`  | Update a todo     |

## Docker

```bash
make image-build     # Build docker image
make image-test-fg   # Run container in foreground with test env vars
make image-run-bg    # Run container in background
make image-stop      # Stop background container
make image-push      # Push image to registry
```

## Kubernetes

```bash
make k8s-apply       # Deploy to kubernetes cluster
make k8s-delete      # Delete from kubernetes cluster
```

## CI/CD

GitHub Actions runs on every push to `main`, tags `v*`, and pull requests.

| Job | Triggers | Steps |
|-----|----------|-------|
| **static-check** | push, PR, tags | Lint (golangci-lint + hadolint) |
| **build** | push, PR, tags | Build Go binary |
| **test** | push, PR, tags | Test with coverage threshold |
| **build-oci-image** | tag push only | Build and push multi-arch OCI image to GHCR |

A cleanup workflow (`cleanup-runs.yml`) removes old workflow runs weekly (retains 7 days, minimum 5 runs).

[Renovate](https://docs.renovatebot.com/) keeps dependencies up to date with platform automerge enabled.

## Skills

Use the following skills when working on related files:

| File(s) | Skill |
|---------|-------|
| `Makefile` | `/makefile` |
| `renovate.json` | `/renovate` |
| `README.md` | `/readme` |
| `.github/workflows/*.yml` | `/ci-workflow` |

When spawning subagents, always pass conventions from the respective skill into the agent's prompt.
