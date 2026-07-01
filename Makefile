.DEFAULT_GOAL := help

# ---------------------------------------------------------------------------
# Constants & tool versions
# ---------------------------------------------------------------------------
OWNER              := andriykalashnykov
PROJECT            := go-todo-web
VERSION            := $(shell cat version.txt 2>/dev/null || echo "dev")
OPV                := $(OWNER)/$(PROJECT):$(VERSION)
WEBPORT            := 8080:8080
CURRENTTAG         := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
NEWTAG             ?= $(shell bash -c 'read -p "Please provide a new tag (current tag - $(CURRENTTAG)): " newtag; echo $$newtag')

# === Tool Versions (pinned) ===
GVM_SHA               := dd652539fa4b771840846f8319fad303c7d0a8d2 # v1.0.22
GOLANGCI_LINT_VERSION := v2.1.6
HADOLINT_VERSION      := 2.12.0
ACT_VERSION           := 0.2.87
NVM_VERSION           := 0.40.4
NODE_VERSION          := 22

# you may need to change to "sudo docker" if not a member of 'docker' group
DOCKERCMD          := docker

BUILD_TIME         := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
# unique id from last git commit
MY_GITREF          := $(shell git rev-parse --short HEAD)

SEMVER_RE          := ^v[0-9]+\.[0-9]+\.[0-9]+$$

# === Go version management ===
GO_VERSIONS := $(shell find . -name 'go.mod' -not -path './vendor/*' -exec grep -oP '^go \K[0-9.]+' {} \; | sort -uV)
GO_VERSION  := $(shell grep -oP '^go \K[0-9.]+' go.mod)

HAS_GVM := $(shell [ -s "$$HOME/.gvm/scripts/gvm" ] && echo true || echo false)
define go-exec
$(if $(filter true,$(HAS_GVM)),bash -c '. $$GVM_ROOT/scripts/gvm && gvm use go$(GO_VERSION) >/dev/null && $(1)',bash -c '$(1)')
endef

# ---------------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------------

#help: @ Show available make targets
help:
	@grep -E '^#[a-zA-Z0-9_-]+:.*@' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*@"}; {sub(/^#/, "", $$1); printf "  \033[1;36m%-20s\033[0m %s\n", $$1, $$2}'

#deps: @ Install and verify required tool dependencies
deps:
	@if [ -z "$$CI" ] && [ ! -s "$$HOME/.gvm/scripts/gvm" ]; then \
		echo "Installing gvm (Go Version Manager)..."; \
		curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/$(GVM_SHA)/binscripts/gvm-installer | bash -s $(GVM_SHA); \
		echo ""; \
		echo "gvm installed. Please restart your shell or run:"; \
		echo "  source $$HOME/.gvm/scripts/gvm"; \
		echo "Then re-run 'make deps' to install Go $(GO_VERSION) via gvm."; \
		exit 0; \
	fi
	@if [ "$(HAS_GVM)" = "true" ]; then \
		for v in $(GO_VERSIONS); do \
			bash -c '. $$GVM_ROOT/scripts/gvm && gvm list' 2>/dev/null | grep -q "go$$v" || { \
				echo "Installing Go $$v via gvm..."; \
				bash -c '. $$GVM_ROOT/scripts/gvm && gvm install go'"$$v"' -B'; \
			}; \
		done; \
	else \
		command -v go >/dev/null 2>&1 || { echo "Error: Go required. Install gvm from https://github.com/moovweb/gvm or Go from https://go.dev/dl/"; exit 1; }; \
	fi
	@$(call go-exec,command -v golangci-lint) >/dev/null 2>&1 || { echo "Installing golangci-lint..."; \
		$(call go-exec,go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)); }

#deps-check: @ Show required Go versions and gvm status
deps-check:
	@echo "Go versions required: $(GO_VERSIONS)"
	@echo "Primary Go version:   $(GO_VERSION)"
	@command -v gvm >/dev/null 2>&1 && { \
		bash -c '. $$GVM_ROOT/scripts/gvm && gvm list'; \
	} || echo "gvm not installed - install from https://github.com/moovweb/gvm"

#deps-act: @ Install act for local CI
deps-act: deps
	@command -v act >/dev/null 2>&1 || { echo "Installing act $(ACT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash -s -- -b /usr/local/bin v$(ACT_VERSION); \
	}

#deps-hadolint: @ Install hadolint for Dockerfile linting
deps-hadolint:
	@command -v hadolint >/dev/null 2>&1 || { echo "Installing hadolint $(HADOLINT_VERSION)..."; \
		curl -sSfL -o /tmp/hadolint https://github.com/hadolint/hadolint/releases/download/v$(HADOLINT_VERSION)/hadolint-Linux-x86_64 && \
		install -m 755 /tmp/hadolint /usr/local/bin/hadolint && \
		rm -f /tmp/hadolint; \
	}

#format: @ Format Go source files
format: deps
	@$(call go-exec,gofmt -w .)

#test: @ Run tests with coverage and race detection
test: deps
	@$(call go-exec,go test -race --cover -parallel=1 -v -coverprofile=coverage.out ./...)
	@$(call go-exec,go tool cover -func=coverage.out | sort -rnk3)

#lint: @ Run golangci-lint and hadolint
lint: deps deps-hadolint
	@$(call go-exec,golangci-lint run ./...)
	@hadolint Dockerfile

#static-check: @ Run all static analysis checks
static-check: lint

#coverage-check: @ Verify test coverage meets minimum threshold
coverage-check: test
	@$(call go-exec,go tool cover -func=coverage.out) | grep total: | \
		awk '{gsub(/%/,"",$$NF); if ($$NF+0 < 80) {printf "FAIL: coverage %s%% below 80%% threshold\n", $$NF; exit 1} else {printf "OK: coverage %s%% meets 80%% threshold\n", $$NF}}'

#build: @ Build the Go binary
build: deps
	@$(call go-exec,CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -a -o manager main.go)

#run: @ Run the application locally
run: deps
	@$(call go-exec,go run main.go)

#image-build: @ Build docker image
image-build: build
	@echo MY_GITREF is $(MY_GITREF)
	@$(DOCKERCMD) buildx build --load --build-arg MY_VERSION=$(VERSION) --build-arg MY_BUILDTIME=$(BUILD_TIME) -f Dockerfile -t $(OPV) .

#clean: @ Clean docker image and build artifacts
clean:
	@$(DOCKERCMD) image rm $(OPV) || true
	@rm -f manager coverage.out

#update: @ Update dependency packages to latest versions
update: deps
	@$(call go-exec,go get -u ./... && go mod tidy)

#image-test-fg: @ Run container in foreground with test overrides
image-test-fg: image-build
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
image-run-bg: image-build
	@$(DOCKERCMD) run -d -p $(WEBPORT) --rm --name $(PROJECT) $(OPV)

#image-cli-bg: @ Get into console of background container
image-cli-bg: image-build
	@$(DOCKERCMD) exec -it $(PROJECT) /bin/sh

#image-logs: @ Tail docker logs
image-logs:
	@$(DOCKERCMD) logs -f $(PROJECT)

#image-stop: @ Stop container running in background
image-stop:
	@$(DOCKERCMD) stop $(PROJECT)

#image-push: @ Push image to registry
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
	@if git rev-parse -q --verify "refs/tags/$(NT)" >/dev/null 2>&1; then \
		echo "ERROR: tag $(NT) already exists locally. Pick a new version or delete it: git tag -d $(NT)"; \
		exit 1; \
	fi
	@if git ls-remote --exit-code --tags origin "refs/tags/$(NT)" >/dev/null 2>&1; then \
		echo "ERROR: tag $(NT) already exists on origin. Pick a new version."; \
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

#ci: @ Run full local CI pipeline
ci: deps static-check coverage-check build
	@echo "Local CI pipeline passed."

#ci-run: @ Run GitHub Actions workflow locally using act
ci-run: deps-act
	@act push --container-architecture linux/amd64 \
		--artifact-server-path /tmp/act-artifacts

#deps-prune: @ Remove unused dependencies
deps-prune: deps
	@$(call go-exec,go mod tidy)

#deps-prune-check: @ Verify no prunable dependencies (CI gate)
deps-prune-check: deps
	@$(call go-exec,go mod tidy)
	@if ! git diff --exit-code go.mod go.sum >/dev/null 2>&1; then \
		echo "ERROR: go.mod/go.sum not tidy. Run 'make deps-prune'."; \
		git checkout go.mod go.sum; \
		exit 1; \
	fi
	@echo "No prunable dependencies found."

#renovate-bootstrap: @ Install nvm and npm for Renovate
renovate-bootstrap:
	@command -v node >/dev/null 2>&1 || { \
		echo "Installing nvm $(NVM_VERSION)..."; \
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v$(NVM_VERSION)/install.sh | bash; \
		export NVM_DIR="$$HOME/.nvm"; \
		[ -s "$$NVM_DIR/nvm.sh" ] && . "$$NVM_DIR/nvm.sh"; \
		nvm install $(NODE_VERSION); \
	}

#renovate-validate: @ Validate Renovate configuration
renovate-validate: renovate-bootstrap
	@npx --yes renovate --platform=local

.PHONY: help deps deps-check deps-act deps-hadolint format test lint \
	static-check coverage-check build run image-build clean update \
	image-test-fg image-test-cli image-run-bg image-cli-bg \
	image-logs image-stop image-push \
	k8s-apply k8s-delete release version ci ci-run \
	deps-prune deps-prune-check \
	renovate-bootstrap renovate-validate
