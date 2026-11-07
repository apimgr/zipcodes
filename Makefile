# Project identity is derived from the origin remote, with directory fallbacks.
PROJECTNAME := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)(\.git)?$$|\1|' || basename "$(CURDIR)")
PROJECTORG := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)/[^/]+(\.git)?$$|\1|' || basename "$(dir $(CURDIR))")
INTERNAL_NAME := $(PROJECTNAME)
VERSION := $(shell if [ -f release.txt ]; then cat release.txt; else echo "0.1.0"; fi)
BUILD_EPOCH := $(shell date -u +%s)
BUILD_DATE := $(shell date -u -d @$(BUILD_EPOCH) +%Y-%m-%dT%H:%M:%SZ)
COMMIT_ID := $(shell git rev-parse --short=7 HEAD 2>/dev/null || echo N/A)
OFFICIAL_SITE := $(shell if [ -f site.txt ]; then cat site.txt; else echo "$${OFFICIAL_SITE:-}"; fi)

BINDIR ?= binaries
RELDIR ?= releases
GO_CACHE ?= $(HOME)/go/pkg/mod
GO_BUILD ?= $(HOME)/.cache/go-build/$(PROJECTNAME)
PLATFORMS ?= linux/amd64,linux/arm64,darwin/amd64,darwin/arm64,windows/amd64,windows/arm64,freebsd/amd64,freebsd/arm64
REGISTRY ?= ghcr.io/$(PROJECTORG)/$(PROJECTNAME)
DOCKER_MEM ?= 4g
DOCKER_CPUS ?= 2

# Every Go tool runs in the maintained project Go image with persistent caches.
GO_DOCKER_RUN := docker run --rm \
	--name $(PROJECTNAME)-$$(tr -dc 'a-z0-9' </dev/urandom | head -c8) \
	--memory=$(DOCKER_MEM) --cpus=$(DOCKER_CPUS) \
	-v "$(CURDIR):/app" \
	-v "$(GO_CACHE):/usr/local/share/go/pkg/mod" \
	-v "$(GO_BUILD):/usr/local/share/go/cache" \
	-w /app -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false
GO_DOCKER := $(GO_DOCKER_RUN) casjaysdev/go:latest
SERVER_LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.CommitID=$(COMMIT_ID) -X main.BuildDate=$(BUILD_DATE)
CLI_LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: dev local build test release docker clean

dev:
	@mkdir -p "$(GO_CACHE)" "$(GO_BUILD)" "$${TMPDIR:-/tmp}/$(PROJECTORG)"
	@BUILD_DIR=$$(mktemp -d "$${TMPDIR:-/tmp}/$(PROJECTORG)/$(INTERNAL_NAME)-XXXXXX") && \
	  echo "Building development binaries in $$BUILD_DIR" && \
	  $(GO_DOCKER_RUN) -v "$$BUILD_DIR:/out" casjaysdev/go:latest sh -c 'go build -buildvcs=false -o /out/$(PROJECTNAME) ./src && if [ -d src/client ]; then go build -buildvcs=false -o /out/$(PROJECTNAME)-cli ./src/client; fi' && \
	  echo "Built binaries: $$BUILD_DIR"

local:
	@mkdir -p "$(BINDIR)" "$(GO_CACHE)" "$(GO_BUILD)"
	@$(GO_DOCKER) sh -c 'go build -buildvcs=false -trimpath -ldflags "$(SERVER_LDFLAGS)" -o $(BINDIR)/$(PROJECTNAME) ./src && if [ -d src/client ]; then go build -buildvcs=false -trimpath -ldflags "$(CLI_LDFLAGS)" -o $(BINDIR)/$(PROJECTNAME)-cli ./src/client; fi'

build:
	@mkdir -p "$(BINDIR)" "$(GO_CACHE)" "$(GO_BUILD)"
	@$(GO_DOCKER) go mod download
	@set -eu; for platform in $$(printf '%s' '$(PLATFORMS)' | tr ',' ' '); do \
	  os=$${platform%/*}; arch=$${platform#*/}; suffix=; \
	  if [ "$$os" = windows ]; then suffix=.exe; fi; \
	  echo "Building server $$os/$$arch"; \
	  $(GO_DOCKER) sh -c "GOOS=$$os GOARCH=$$arch go build -buildvcs=false -trimpath -ldflags '$(SERVER_LDFLAGS)' -o $(BINDIR)/$(PROJECTNAME)-$$os-$$arch$$suffix ./src"; \
	  if [ -d src/client ]; then \
	    echo "Building CLI $$os/$$arch"; \
	    $(GO_DOCKER) sh -c "GOOS=$$os GOARCH=$$arch go build -buildvcs=false -trimpath -ldflags '$(CLI_LDFLAGS)' -o $(BINDIR)/$(PROJECTNAME)-cli-$$os-$$arch$$suffix ./src/client"; \
	  fi; \
	done

# Tests and coverage artifacts are written under the project-specific system temp path.
test:
	@mkdir -p "$(GO_CACHE)" "$(GO_BUILD)" "$${TMPDIR:-/tmp}/$(PROJECTORG)"
	@$(GO_DOCKER) sh -c 'mkdir -p "$${TMPDIR:-/tmp}/$(PROJECTORG)" && coverage_dir=$$(mktemp -d "$${TMPDIR:-/tmp}/$(PROJECTORG)/$(INTERNAL_NAME)-XXXXXX") && trap "rm -rf $$coverage_dir" EXIT && go vet ./... && govulncheck ./... && go test -v -coverprofile="$$coverage_dir/coverage.out" ./... && go tool cover -func="$$coverage_dir/coverage.out" && coverage=$$(go tool cover -func="$$coverage_dir/coverage.out" | awk "/^total:/ {print \$$3}") && awk -v coverage="$$coverage" "BEGIN { if (coverage+0 < 60) exit 1 }" '

release: build
	@mkdir -p "$(RELDIR)"
	@printf '%s\n' '$(VERSION)' > "$(RELDIR)/version.txt"
	@cp "$(BINDIR)"/$(PROJECTNAME)-* "$(RELDIR)/"
	@tar --exclude='./.git' --exclude='./binaries' --exclude='./releases' -czf "$(RELDIR)/$(PROJECTNAME)-$(VERSION)-source.tar.gz" .
	@cd "$(RELDIR)" && sha256sum * > sha256.txt

# Build and publish the multi-architecture container image.
docker:
	@docker buildx build --platform linux/amd64,linux/arm64 --push \
	  --build-arg VERSION="$(VERSION)" \
	  --build-arg BUILD_DATE="$(BUILD_DATE)" \
	  --build-arg BUILD_EPOCH="$(BUILD_EPOCH)" \
	  --build-arg COMMIT_ID="$(COMMIT_ID)" \
	  -t "$(REGISTRY):$(VERSION)" -t "$(REGISTRY):latest" \
	  -f docker/Dockerfile .

# Remove build and release artifacts.
clean:
	@rm -rf $(BINDIR) $(RELDIR)
