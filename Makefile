# Variables
PROJECTNAME = zipcodes
PROJECTORG = apimgr

# Version management - use VERSION env var if set, otherwise read from release.txt
VERSION ?= $(shell cat release.txt 2>/dev/null || grep -m1 "Version.*=" src/main.go | sed 's/.*"\(.*\)".*/\1/')
COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -w -s"

.PHONY: build releases test docker docker-dev clean

# Build for all platforms
build:
	@echo "Building $(PROJECTNAME) $(VERSION) for all platforms..."
	@mkdir -p binaries release

	@echo "  → Linux AMD64"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-linux-amd64 ./src

	@echo "  → Linux ARM64"
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-linux-arm64 ./src

	@echo "  → Windows AMD64"
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-windows-amd64.exe ./src

	@echo "  → Windows ARM64"
	@GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-windows-arm64.exe ./src

	@echo "  → macOS AMD64"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-darwin-amd64 ./src

	@echo "  → macOS ARM64"
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-darwin-arm64 ./src

	@echo "  → FreeBSD AMD64"
	@GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-freebsd-amd64 ./src

	@echo "  → FreeBSD ARM64"
	@GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-freebsd-arm64 ./src

	@echo "  → Host"
	@CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME) ./src

	@echo ""
	@echo "✓ Build complete! Version: $(VERSION)"
	@echo "  Binaries: ./binaries/"

# Release to GitHub
release: build
	@echo "Preparing releases $(VERSION)..."
	@mkdir -p release

	@# Copy binaries to releases directory
	@cp binaries/$(PROJECTNAME)-* releases/ 2>/dev/null || true

	@# Strip musl binaries if they exist
	@for bin in releases/$(PROJECTNAME)-linux-*; do \
		if ldd $$bin 2>&1 | grep -q musl; then \
			echo "  → Stripping $$bin"; \
			strip $$bin 2>/dev/null || true; \
		fi; \
	done

	@# Update version in release.txt for next build
	@echo "$(VERSION)" | awk -F. '{$$NF = $$NF + 1;} 1' OFS=. > release.txt
	@NEXT_VERSION=$$(cat release.txt); \
	echo "  → Updated release.txt to $$NEXT_VERSION"

	@# Create GitHub release
	@echo "  → Creating GitHub release..."
	@mkdir -p releases
	@echo "Copying platform binaries to releases..."
	@cp binaries/$(PROJECTNAME)-* releases/ 2>/dev/null || { echo "Error: Build first"; exit 1; }
	@echo "Creating source archives (no VCS files)..."
	@git archive --format=tar.gz --prefix=$(PROJECTNAME)-$(VERSION)/ HEAD -o releases/$(PROJECTNAME)-$(VERSION)-src.tar.gz
	@git archive --format=zip --prefix=$(PROJECTNAME)-$(VERSION)/ HEAD -o releases/$(PROJECTNAME)-$(VERSION)-src.zip
	@gh releases delete $(VERSION) -y 2>/dev/null || true
	@git tag -d $(VERSION) 2>/dev/null || true
	@git tag $(VERSION)
	@git push origin :refs/tags/$(VERSION) 2>/dev/null || true
	@git push origin $(VERSION)
	@gh releases create $(VERSION) \
		--title "$(PROJECTNAME) $(VERSION)" \
		--generate-notes \
		releases/*

	@echo ""
	@echo "✓ Release $(VERSION) created!"
	@echo "  Next version will be: $$(cat release.txt)"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "✓ Tests complete!"

# Build and push Docker images (multi-arch)
docker: build
	@echo "Building Docker images for $(PROJECTNAME) $(VERSION)..."

	@# Build for both architectures
	@docker buildx build --platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION) \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):latest \
		--push .

	@echo ""
	@echo "✓ Docker images pushed!"
	@echo "  ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION)"
	@echo "  ghcr.io/$(PROJECTORG)/$(PROJECTNAME):latest"

# Build Docker image for development (local only, not pushed)
docker-dev:
	@echo "Building development Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(PROJECTNAME):dev \
		.
	@echo "✓ Docker development image built: $(PROJECTNAME):dev"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf binaries/ releases/
	@rm -f coverage.out
	@go clean
	@echo "✓ Clean complete!"
