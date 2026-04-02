[![CI](https://github.com/AndriyKalashnykov/go-todo-web/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/AndriyKalashnykov/go-todo-web/actions/workflows/ci.yml)
[![Hits](https://hits.sh/github.com/AndriyKalashnykov/go-todo-web.svg?view=today-total&style=plastic)](https://hits.sh/github.com/AndriyKalashnykov/go-todo-web/)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/AndriyKalashnykov/go-todo-web)

# Go Todo Web

HTTP REST API for managing todo items, built with the [Echo](https://echo.labstack.com/) framework (v5) in Go. Supports CRUD operations and is containerized for Kubernetes deployment with multi-arch Docker images (amd64/arm64).

## Quick Start

```bash
make deps      # verify required tool dependencies
make build     # build the Go binary
make test      # run tests with coverage
make run       # start the application on port 8080
```

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| [GNU Make](https://www.gnu.org/software/make/) | 3.81+ | Build orchestration |
| [Go](https://go.dev/dl/) | 1.26+ | Language runtime and compiler |
| [Docker](https://www.docker.com/) | latest | Container image builds |
| [kubectl](https://kubernetes.io/docs/tasks/tools/) | latest | Kubernetes deployment |
| [golangci-lint](https://golangci-lint.run/) | 2.1.6+ | Go linting (auto-installed by `make deps`) |
| [gvm](https://github.com/moovweb/gvm) | latest | Go version management (optional) |

Install all required dependencies:

```bash
make deps
```

## Available Make Targets

Run `make help` to see all available targets.

### Build & Run

| Target | Description |
|--------|-------------|
| `make build` | Build the Go binary |
| `make run` | Run the application locally |
| `make clean` | Clean docker image and build artifacts |
| `make update` | Update dependency packages to latest versions |

### Testing & Code Quality

| Target | Description |
|--------|-------------|
| `make test` | Run tests with coverage and race detection |
| `make lint` | Run golangci-lint and hadolint |
| `make format` | Format Go source files |
| `make static-check` | Run all static analysis checks |
| `make coverage-check` | Verify test coverage meets minimum threshold |

### Docker

| Target | Description |
|--------|-------------|
| `make image-build` | Build docker image |
| `make image-test-fg` | Run container in foreground with test overrides |
| `make image-test-cli` | Run container in foreground with shell entrypoint |
| `make image-run-bg` | Run container in background |
| `make image-cli-bg` | Get into console of background container |
| `make image-logs` | Tail docker logs |
| `make image-stop` | Stop container running in background |
| `make image-push` | Push image to registry |

### Kubernetes

| Target | Description |
|--------|-------------|
| `make k8s-apply` | Deploy to kubernetes cluster |
| `make k8s-delete` | Delete from kubernetes cluster |

### CI

| Target | Description |
|--------|-------------|
| `make ci` | Run full local CI pipeline |
| `make ci-run` | Run GitHub Actions workflow locally using act |

### Utilities

| Target | Description |
|--------|-------------|
| `make help` | Show available make targets |
| `make deps` | Install and verify required tool dependencies |
| `make deps-check` | Show required Go versions and gvm status |
| `make deps-act` | Install act for local CI |
| `make deps-hadolint` | Install hadolint for Dockerfile linting |
| `make deps-prune` | Remove unused dependencies |
| `make deps-prune-check` | Verify no prunable dependencies (CI gate) |
| `make version` | Print current version (tag) |
| `make release` | Create and push a new tag (semver validated) |
| `make renovate-bootstrap` | Install nvm and npm for Renovate |
| `make renovate-validate` | Validate Renovate configuration |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/create` | Create a todo |
| GET | `/get/:id` | Get a todo by ID |
| GET | `/all` | List all todos |
| DELETE | `/delete/:id` | Delete a todo |
| PATCH | `/update/:id` | Update a todo |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Listen port | `8080` |
| `APP_CONTEXT` | Base context path | `/` |

### Kubernetes Downward API Variables

| Variable | Description |
|----------|-------------|
| `MY_NODE_NAME` | Name of k8s node |
| `MY_POD_NAME` | Name of k8s pod |
| `MY_POD_NAMESPACE` | Namespace of k8s pod |
| `MY_POD_IP` | K8s pod IP |
| `MY_POD_SERVICE_ACCOUNT` | Service account of k8s pod |

## CI/CD

GitHub Actions runs on every push to `main`, tags `v*`, and pull requests.

| Job | Triggers | Steps |
|-----|----------|-------|
| **static-check** | push, PR, tags | Lint (golangci-lint + hadolint) |
| **build** | push, PR, tags | Build Go binary |
| **test** | push, PR, tags | Test with coverage threshold |
| **build-oci-image** | tag push only | Build and push multi-arch OCI image to GHCR |

A cleanup workflow removes old workflow runs weekly (retains 7 days, minimum 5 runs).

[Renovate](https://docs.renovatebot.com/) keeps dependencies up to date with platform automerge enabled.

## Pulling Image from GitHub Container Registry

```bash
docker pull ghcr.io/andriykalashnykov/go-todo-web:latest
```
