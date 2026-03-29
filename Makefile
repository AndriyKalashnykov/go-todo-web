.DEFAULT_GOAL := help

# ---------------------------------------------------------------------------
# Constants & tool versions
# ---------------------------------------------------------------------------
OWNER              := andriykalashnykov
PROJECT            := go-todo-web
VERSION            := v0.0.1
OPV                := $(OWNER)/$(PROJECT):$(VERSION)
WEBPORT            := 8080:8080
CURRENTTAG         := $(shell git describe --tags --abbrev=0)
NEWTAG             ?= $(shell bash -c 'read -p "Please provide a new tag (current tag - ${CURRENTTAG}): " newtag; echo $$newtag')

GOLANGCI_LINT_VERSION := v2.1.6
NVM_VERSION           := 0.40.4

# you may need to change to "sudo docker" if not a member of 'docker' group
DOCKERCMD          := "docker"

BUILD_TIME         := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
# unique id from last git commit
MY_GITREF          := $(shell git rev-parse --short HEAD)

SEMVER_RE          := ^v[0-9]+\.[0-9]+\.[0-9]+$$

# ---------------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------------

#help: @ Show available make targets
help:
	@grep -E '^#[a-zA-Z0-9_-]+:.*@' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*@"}; {sub(/^#/, "", $$1); printf "  \033[1;36m%-20s\033[0m %s\n", $$1, $$2}'

#deps: @ Verify required tool dependencies
deps:
	@command -v go >/dev/null 2>&1        || { echo "go is required but not installed"; exit 1; }
	@command -v docker >/dev/null 2>&1    || { echo "docker is required but not installed"; exit 1; }
	@command -v kubectl >/dev/null 2>&1   || { echo "kubectl is required but not installed"; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is required but not installed"; exit 1; }

#test: @ Run tests with coverage
test:
	@go test --cover -parallel=1 -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | sort -rnk3

#lint: @ Run golangci-lint
lint:
	@golangci-lint run ./...

#build: @ Build the Go binary
build:
	@CGO_ENABLED=0 go build -ldflags "-X main.Version=${MY_VERSION} -X main.BuildTime=${MY_BUILDTIME}" -a -o manager main.go

#run: @ Run the application locally
run:
	@go run main.go

#image: @ Build docker image
image:
	@echo MY_GITREF is $(MY_GITREF)
	@$(DOCKERCMD) buildx build --load --build-arg MY_VERSION=$(VERSION) --build-arg MY_BUILDTIME=$(BUILD_TIME) -f Dockerfile -t $(OPV) .

#clean: @ Clean docker image
clean:
	@$(DOCKERCMD) image rm $(OPV) | true

#update: @ Update dependency packages to latest versions
update:
	@go get -u ./...; go mod tidy

#image-test-fg: @ Run container in foreground with test overrides
image-test-fg: image
	@$(DOCKERCMD) run -it -p $(WEBPORT) \
	-e APP_CONTEXT=/myhello/ \
	-e MY_NODE_NAME=node1 \
	-e MY_POD_NAME=pod1 \
	-e MY_POD_NAMESPACE=ns1 \
	-e MY_POD_IP=podip1 \
	-e MY_POD_SERVICE_ACCOUNT=podsa1 \
	--rm $(OPV)

#image-test-cli: @ Run container in foreground with shell entrypoint
image-test-cli:
	@$(DOCKERCMD) run -it --rm --entrypoint "/bin/sh" $(OPV)

#image-run-bg: @ Run container in background
image-run-bg: image
	@$(DOCKERCMD) run -d -p $(WEBPORT) --rm --name $(PROJECT) $(OPV)

#image-cli-bg: @ Get into console of background container
image-cli-bg: image
	@$(DOCKERCMD) exec -it $(PROJECT) /bin/sh

#image-logs: @ Tail docker logs
image-logs:
	@$(DOCKERCMD) logs -f $(PROJECT)

#image-stop: @ Stop container running in background
image-stop:
	@$(DOCKERCMD) stop $(PROJECT)

#image-push: @ Push image to Docker Hub
image-push:
	@$(DOCKERCMD) push $(OPV)

#k8s-apply: @ Deploy to kubernetes cluster
k8s-apply:
	@sed -e 's/v0.0.1/$(VERSION)/' go-todo-web.yaml | kubectl apply -f -

#k8s-delete: @ Delete from kubernetes cluster
k8s-delete:
	@kubectl delete -f go-todo-web.yaml

#release: @ Create and push a new tag (semver validated)
release:
	$(eval NT=$(NEWTAG))
	@if ! echo "$(NT)" | grep -qE '$(SEMVER_RE)'; then \
		echo "Error: '$(NT)' is not a valid semver tag (expected vX.Y.Z)"; \
		exit 1; \
	fi
	@echo -n "Are you sure to create and push ${NT} tag? [y/N] " && read ans && [ $${ans:-N} = y ]
	@echo ${NT} > ./version.txt
	@git add -A
	@git commit -a -s -m "Cut ${NT} release"
	@git tag ${NT}
	@git push origin ${NT}
	@git push
	@echo "Done."

#version: @ Print current version (tag)
version:
	@echo $(shell git describe --tags --abbrev=0)

#ci: @ Run lint, test, and build (CI pipeline)
ci: lint test build

#renovate-bootstrap: @ Install nvm and npm for Renovate
renovate-bootstrap:
	@command -v node >/dev/null 2>&1 || { \
		echo "Installing nvm $(NVM_VERSION)..."; \
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v$(NVM_VERSION)/install.sh | bash; \
		export NVM_DIR="$$HOME/.nvm"; \
		[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh"; \
		nvm install --lts; \
	}

#renovate-validate: @ Validate Renovate configuration
renovate-validate: renovate-bootstrap
	@npx --yes renovate --platform=local

.PHONY: help deps test lint build run image clean update \
	image-test-fg image-test-cli image-run-bg image-cli-bg \
	image-logs image-stop image-push \
	k8s-apply k8s-delete release version ci \
	renovate-bootstrap renovate-validate
