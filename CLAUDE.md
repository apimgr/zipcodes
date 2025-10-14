# 🚀 Universal Go API Server - Generic Specification v2.0

**Template for building self-contained Go API servers**

This specification provides a complete template for building production-ready Go API servers with embedded data, admin interfaces, and OS-specific directory management.

---

## 📖 Table of Contents

1. [Quick Start](#quick-start)
2. [Project Setup](#project-setup)
3. [Architecture Overview](#architecture-overview)
4. [Directory Layout](#directory-layout)
5. [Data Management](#data-management)
6. [Authentication](#authentication)
7. [Routes & Endpoints](#routes--endpoints)
8. [Configuration](#configuration)
9. [Build System](#build-system)
10. [Testing](#testing)
11. [Deployment](#deployment)
12. [Optional Features](#optional-features)

---

## 🎯 Quick Start

### Adapt This Template

1. Replace `{projectname}` with your project name (lowercase, no spaces)
2. Replace `{organization}` with your GitHub organization/username
3. Replace `{purpose}` with a brief description of your project
4. Customize data sources for your specific needs
5. Modify API endpoints for your domain

### Example Projects Using This Template

- **Airports API**: Public airport location information with GeoIP
- **Weather Stations**: Meteorological data with geographic search
- **Store Locator**: Retail locations with proximity search
- **Device Registry**: IoT devices with status monitoring

---

## 🏗️ Project Setup

### Initial Configuration

```yaml
Project Details:
  Name: {projectname}
  Module: github.com/{organization}/{projectname}
  Language: Go 1.21+
  Purpose: {purpose}

Binary Name: {projectname}

Default Ports:
  HTTP: Random (64000-64999) on first run
  HTTPS: HTTP port + 363 (if enabled)
```

### Required go.mod

```go
module github.com/{organization}/{projectname}

go 1.23

require (
	github.com/go-chi/chi/v5 v5.2.3
	modernc.org/sqlite v1.39.0
)
```

### Project Structure

```
./
├── src/                        # All source code
│   ├── main.go                 # Entry point
│   ├── server/                 # HTTP server & handlers
│   ├── database/               # Database layer
│   ├── auth/                   # Authentication
│   ├── paths/                  # OS-specific paths
│   ├── scheduler/              # Task scheduler
│   ├── {domain}/               # Your domain logic
│   │   ├── data.go             # Data loading/indexing
│   │   ├── service.go          # Business logic
│   │   └── handlers.go         # HTTP handlers
│   ├── static/                 # Embedded assets
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/              # HTML templates
│       └── *.html
├── test/
│   ├── integration/
│   └── unit/
├── scripts/
│   ├── install.sh
│   ├── linux.sh
│   ├── macos.sh
│   └── windows.ps1
├── binaries/                   # Built binaries (gitignored)
├── release/                    # Release artifacts
├── Makefile                    # Build automation
├── Dockerfile                  # Container definition
├── docker-compose.yml          # Compose configuration
├── Jenkinsfile                 # CI/CD pipeline
├── README.md                   # User documentation
├── SPEC.md                     # This file
├── LICENSE.md                  # License
├── release.txt                 # Version tracking
├── .gitignore
└── .dockerignore
```

---

## 🏗️ Architecture Overview

### System Design

```
┌─────────────────────────────────────────┐
│         Single Go Binary                │
│  ┌─────────────────────────────────┐   │
│  │  Embedded Assets (go:embed)     │   │
│  │  • Data files (JSON/CSV/etc)    │   │
│  │  • HTML templates               │   │
│  │  • CSS/JS/Images                │   │
│  │  • Static assets                │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  In-Memory Data Structures      │   │
│  │  • Indexes for fast lookup      │   │
│  │  • Cached computations          │   │
│  │  • Runtime state                │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  HTTP Server (Chi Router)       │   │
│  │  • Public routes (no auth)      │   │
│  │  • Admin routes (auth required) │   │
│  │  • API v1 endpoints             │   │
│  │  • Documentation routes         │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  SQLite Database                │   │
│  │  • Admin credentials (hashed)   │   │
│  │  • Server settings              │   │
│  │  • Audit logs                   │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  Task Scheduler                 │   │
│  │  • Periodic updates             │   │
│  │  • Maintenance tasks            │   │
│  │  • Data refresh                 │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### Technology Stack

```yaml
Core:
  Language: Go 1.21+
  HTTP Router: Chi v5
  Database: SQLite (modernc.org/sqlite - pure Go, no CGO)
  Templates: Go html/template
  Embedding: Go embed.FS

Authentication:
  Hashing: SHA-256
  Tokens: Bearer tokens
  Web Auth: Basic Auth
  Sessions: Secure HttpOnly cookies

Optional:
  GeoIP: oschwald/geoip2-golang (if geolocation needed)
  Metrics: Prometheus (disabled by default)
  Caching: In-memory or external
```

---

## 📁 Directory Layout

### OS-Specific Paths

```yaml
Linux/BSD (with root privileges):
  Config:  /etc/{projectname}/
  Data:    /var/lib/{projectname}/
  Logs:    /var/log/{projectname}/
  Runtime: /run/{projectname}/

Linux/BSD (without root):
  Config:  ~/.config/{projectname}/
  Data:    ~/.local/share/{projectname}/
  Logs:    ~/.local/state/{projectname}/
  Runtime: ~/.local/run/{projectname}/

macOS (with privileges):
  Config:  /Library/Application Support/{ProjectName}/
  Data:    /Library/Application Support/{ProjectName}/data/
  Logs:    /Library/Logs/{ProjectName}/
  Runtime: /var/run/{projectname}/

macOS (without privileges):
  Config:  ~/Library/Application Support/{ProjectName}/
  Data:    ~/Library/Application Support/{ProjectName}/data/
  Logs:    ~/Library/Logs/{ProjectName}/
  Runtime: ~/Library/Application Support/{ProjectName}/run/

Windows (system):
  Config:  C:\ProgramData\{ProjectName}\config\
  Data:    C:\ProgramData\{ProjectName}\data\
  Logs:    C:\ProgramData\{ProjectName}\logs\
  Runtime: C:\ProgramData\{ProjectName}\run\

Windows (user):
  Config:  %APPDATA%\{ProjectName}\
  Data:    %LOCALAPPDATA%\{ProjectName}\
  Logs:    %LOCALAPPDATA%\{ProjectName}\logs\
  Runtime: %LOCALAPPDATA%\{ProjectName}\run\

Docker:
  Config:  /config
  Data:    /data
  Logs:    /data/logs
  Runtime: /tmp
```

### Directory Detection Code

```go
// src/paths/paths.go
package paths

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func GetDefaultDirs(projectName string) (configDir, dataDir, logsDir string) {
	isRoot := false
	if runtime.GOOS == "windows" {
		isRoot = os.Getenv("USERDOMAIN") == os.Getenv("COMPUTERNAME")
	} else {
		isRoot = os.Geteuid() == 0
	}

	if isRoot {
		switch runtime.GOOS {
		case "windows":
			programData := os.Getenv("ProgramData")
			if programData == "" {
				programData = "C:\\ProgramData"
			}
			configDir = filepath.Join(programData, projectName, "config")
			dataDir = filepath.Join(programData, projectName, "data")
			logsDir = filepath.Join(programData, projectName, "logs")
		default: // Linux, BSD, macOS
			configDir = filepath.Join("/etc", projectName)
			dataDir = filepath.Join("/var/lib", projectName)
			logsDir = filepath.Join("/var/log", projectName)
		}
	} else {
		// User-specific paths (see full implementation in template)
		var homeDir string
		currentUser, err := user.Current()
		if err == nil {
			homeDir = currentUser.HomeDir
		}
		// ... continue with XDG/macOS/Windows user paths
	}

	return configDir, dataDir, logsDir
}
```

### Environment Variables & Flags

```yaml
Directory Overrides (in priority order):
  1. Command-line flags
  2. Environment variables
  3. OS-specific defaults

Command-Line Flags:
  --config DIR        # Configuration directory
  --data DIR          # Data directory
  --logs DIR          # Logs directory
  --port PORT         # HTTP port (default: random 64000-64999)
  --address ADDR      # Listen address (default: 0.0.0.0)
  --dev               # Development mode
  --version           # Show version
  --help              # Show help

Environment Variables:
  CONFIG_DIR          # Configuration directory
  DATA_DIR            # Data directory
  LOGS_DIR            # Logs directory
  PORT                # Server port
  ADDRESS             # Listen address
  DATABASE_URL        # Full connection string
  ADMIN_USER          # Admin username (first run only)
  ADMIN_PASSWORD      # Admin password (first run only)
  ADMIN_TOKEN         # Admin API token (first run only)
```

---

## 💾 Data Management

### Embedded Data (Small Datasets < 20MB)

```go
//go:embed data/{yourdata}.json
var embeddedData embed.FS

func LoadData() (YourDataType, error) {
	data, err := embeddedData.ReadFile("data/{yourdata}.json")
	if err != nil {
		return nil, err
	}

	var result YourDataType
	err = json.Unmarshal(data, &result)
	return result, err
}
```

### Downloaded Data (Large Datasets > 20MB)

```go
const dataURL = "https://example.com/data/{yourdata}.json"

func (s *Service) DownloadData() error {
	dataPath := filepath.Join(s.dataDir, "{yourdata}.json")

	// Check if already exists
	if fileExists(dataPath) {
		return nil
	}

	// Download
	resp, err := http.Get(dataURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Save to disk
	out, err := os.Create(dataPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
```

### In-Memory Indexing

```go
type Service struct {
	data    map[string]*Item
	indexes struct {
		byName     map[string][]*Item
		byCategory map[string][]*Item
		byLocation map[string][]*Item
	}
}

func (s *Service) BuildIndexes() {
	s.indexes.byName = make(map[string][]*Item)
	s.indexes.byCategory = make(map[string][]*Item)

	for _, item := range s.data {
		// Index by name (case-insensitive)
		nameLower := strings.ToLower(item.Name)
		s.indexes.byName[nameLower] = append(s.indexes.byName[nameLower], item)

		// Index by category
		s.indexes.byCategory[item.Category] = append(s.indexes.byCategory[item.Category], item)
	}
}
```

---

## 🔐 Authentication

### Admin Authentication System

```yaml
Admin Account:
  Username: administrator (fixed)
  Password: Auto-generated on first run (or from ADMIN_PASSWORD env var)
  Token: Auto-generated SHA-256 hash
  Storage: SQLite database (hashed)

Credentials File:
  Location: {CONFIG_DIR}/admin_credentials
  Permissions: 0600 (read/write for owner only)
  Format:
    Username: administrator
    Password: {generated_password}
    API Token: {generated_token}

    Keep these credentials secure!
    They will not be shown again.

Authentication Methods:
  Web UI: Basic Auth (username + password)
  API: Bearer Token (Authorization: Bearer {token})
```

### Database Schema

```sql
CREATE TABLE admin_credentials (
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  token_hash TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  type TEXT NOT NULL, -- string, number, boolean, json
  description TEXT,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE audit_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  action TEXT NOT NULL,
  resource TEXT NOT NULL,
  old_value TEXT,
  new_value TEXT,
  success BOOLEAN NOT NULL
);
```

---

## 🛣️ Routes & Endpoints

### Route Structure

```yaml
Public Routes (No Authentication):
  GET  /                          - Homepage
  GET  /healthz                   - Health check
  GET  /static/*                  - Static assets
  GET  /api/v1/health             - API health check

API Routes (Public):
  GET  /api/v1/{resource}         - List resources
  GET  /api/v1/{resource}/:id     - Get resource by ID
  GET  /api/v1/{resource}/search  - Search resources
  GET  /api/v1/{resource}/nearby  - Geographic search (if applicable)

Documentation Routes (Public):
  GET  /docs                      - Swagger UI
  GET  /api/v1/docs               - API documentation
  GET  /api/v1/openapi.json       - OpenAPI spec
  GET  /graphql                   - GraphQL Playground (optional)
  POST /api/v1/graphql            - GraphQL endpoint (optional)

Admin Routes (Authentication Required):
  GET  /admin                     - Admin dashboard
  GET  /admin/settings            - Server settings
  POST /admin/settings            - Update settings
  GET  /admin/database            - Database management
  POST /admin/database/test       - Test database connection
  GET  /admin/logs                - Log viewer
  GET  /admin/audit               - Audit log

Admin API (Bearer Token Required):
  GET  /api/v1/admin              - Admin info
  GET  /api/v1/admin/settings     - Get settings
  PUT  /api/v1/admin/settings     - Update settings
  POST /api/v1/admin/reload       - Reload configuration
  GET  /api/v1/admin/stats        - Server statistics
```

### Response Format

```json
{
  "success": true,
  "data": { },
  "timestamp": "2024-01-01T12:00:00Z"
}

{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

## ⚙️ Configuration

### Configuration Storage

```yaml
Priority Order:
  1. Command-line flags (highest priority)
  2. Environment variables
  3. Database settings
  4. Built-in defaults (lowest priority)

Database-Driven Configuration:
  - All settings stored in SQLite database
  - No configuration files (no .env, no config.json)
  - Changes via admin UI update database
  - Automatic reload (no restart required for most settings)

Settings Categories:
  server.*     - Server configuration
  database.*   - Database settings
  security.*   - Security policies
  features.*   - Feature toggles
  monitoring.* - Monitoring/metrics
```

### Common Settings

```yaml
Server Settings:
  server.title: "{ProjectName}"              # Display name
  server.tagline: "{Your tagline}"           # Subtitle
  server.description: "{Full description}"   # About text
  server.http_port: 8080                     # HTTP port
  server.https_port: 8443                    # HTTPS port (if enabled)
  server.https_enabled: false                # Enable HTTPS
  server.timezone: "UTC"                     # Server timezone
  server.date_format: "US"                   # US, EU, ISO
  server.time_format: "12-hour"              # 12-hour, 24-hour

Security Settings:
  security.session_timeout: 30               # Minutes
  security.max_login_attempts: 5             # Before lockout
  security.lockout_duration: 15              # Minutes
  security.password_min_length: 8            # Characters
  security.require_special_char: false       # In passwords

Feature Settings:
  features.registration_enabled: false       # Public registration
  features.api_enabled: true                 # Enable API
  features.graphql_enabled: false            # Enable GraphQL
  features.metrics_enabled: false            # Prometheus metrics
```

---

## 🔨 Build System

### Makefile

```makefile
# Variables
PROJECTNAME = {projectname}
PROJECTORG = {organization}
VERSION = $(shell cat release.txt 2>/dev/null || echo "0.0.1")
COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -w -s"

.PHONY: build release test docker clean deps

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	@go mod download
	@go mod tidy

# Build for all platforms using Docker Alpine
build: deps
	@echo "Building $(PROJECTNAME) v$(VERSION) for all platforms using Docker Alpine..."
	@mkdir -p binaries release
	@docker run --rm -v $$(pwd):/workspace -w /workspace golang:1.23-alpine sh -c ' \
		apk add --no-cache git make > /dev/null 2>&1 && \
		echo "  → Linux AMD64" && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-linux-amd64 ./src && \
		echo "  → Linux ARM64" && \
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-linux-arm64 ./src && \
		echo "  → Windows AMD64" && \
		GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-windows-amd64.exe ./src && \
		echo "  → Windows ARM64" && \
		GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-windows-arm64.exe ./src && \
		echo "  → macOS AMD64" && \
		GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-macos-amd64 ./src && \
		echo "  → macOS ARM64" && \
		GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-macos-arm64 ./src && \
		echo "  → FreeBSD AMD64" && \
		GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-bsd-amd64 ./src && \
		echo "  → FreeBSD ARM64" && \
		GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME)-bsd-arm64 ./src && \
		echo "  → Host" && \
		CGO_ENABLED=0 go build $(LDFLAGS) -o binaries/$(PROJECTNAME) ./src \
	'
	@cp binaries/$(PROJECTNAME)-* release/ 2>/dev/null || true

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

# Build Docker image
docker:
	@echo "Building Docker image..."
	@docker build -t $(PROJECTNAME):latest \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf binaries/ release/
	@rm -f coverage.out
	@go clean
```

### Dockerfile (Multi-Stage)

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE} -w -s" \
    -o {projectname} \
    ./src

# Runtime stage
FROM alpine:latest

# Add ca-certificates
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -S {projectname} && adduser -S {projectname} -G {projectname}

# Copy binary
COPY --from=builder /build/{projectname} /usr/local/bin/{projectname}

# Create directories
RUN mkdir -p /config /data /logs && \
    chown -R {projectname}:{projectname} /config /data /logs

# Switch to non-root
USER {projectname}

# Set working directory
WORKDIR /data

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["{projectname}", "--status"]

EXPOSE 8080

ENTRYPOINT ["{projectname}"]
```

### docker-compose.yml

```yaml
services:
  {projectname}:
    image: {projectname}:latest
    container_name: {projectname}
    restart: unless-stopped
    environment:
      - PORT=8080
    volumes:
      - ./data:/data
      - ./config:/config
      - ./logs:/logs
    ports:
      - "8080:8080"
    healthcheck:
      test: ["{projectname}", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
```

---

## 🧪 Testing

### Test Structure

```go
// test/integration/api_test.go
package integration

import (
	"net/http/httptest"
	"testing"

	"{module}/src/server"
	"{module}/src/{domain}"
)

func setupTestServer(t *testing.T) *httptest.Server {
	// Initialize services
	svc, err := {domain}.NewService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Create server
	srv := server.New(svc, false)
	return httptest.NewServer(srv.Router())
}

func TestAPIEndpoints(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	tests := []struct {
		name       string
		endpoint   string
		wantStatus int
	}{
		{"Health check", "/healthz", 200},
		{"Get resource", "/api/v1/resource/1", 200},
		{"Not found", "/api/v1/resource/notfound", 404},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tt.endpoint)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -v -race ./...

# Run benchmarks
go test -v -bench=. -benchmem ./...
```

---

## 🚀 Deployment

### Binary Deployment

```bash
# 1. Build for target platform
make build

# 2. Copy binary to server
scp binaries/{projectname}-linux-amd64 server:/usr/local/bin/{projectname}

# 3. Make executable
ssh server "chmod +x /usr/local/bin/{projectname}"

# 4. Run
ssh server "{projectname} --port 8080"
```

### Systemd Service

```ini
# /etc/systemd/system/{projectname}.service
[Unit]
Description={ProjectName} API Server
After=network.target

[Service]
Type=simple
User={projectname}
Group={projectname}
ExecStart=/usr/local/bin/{projectname} --port 8080
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Docker Deployment

```bash
# Build image
make docker

# Run container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

---

## ⚡ Optional Features

### Task Scheduler

```go
// src/scheduler/scheduler.go
package scheduler

type Scheduler struct {
	tasks  []*Task
	stopCh chan struct{}
}

func New() *Scheduler {
	return &Scheduler{
		tasks:  make([]*Task, 0),
		stopCh: make(chan struct{}),
	}
}

func (s *Scheduler) AddTask(name, schedule string, handler func() error) {
	// Schedule format: "0 3 * * 0" = Sunday at 3 AM
	task := &Task{
		Name:     name,
		Schedule: schedule,
		Handler:  handler,
	}
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Start() {
	go s.run()
}
```

### GeoIP Integration

```go
// Optional: Add geographic features
import "github.com/oschwald/geoip2-golang"

type GeoIPService struct {
	cityDB *geoip2.Reader
}

func NewGeoIPService(configDir string) (*GeoIPService, error) {
	// Download GeoLite2-City.mmdb on first run
	dbPath := filepath.Join(configDir, "geoip", "GeoLite2-City.mmdb")

	if !fileExists(dbPath) {
		if err := downloadGeoIPDatabase(dbPath); err != nil {
			return nil, err
		}
	}

	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}

	return &GeoIPService{cityDB: db}, nil
}
```

### Prometheus Metrics

```go
// Optional: Expose metrics (disabled by default)
import "github.com/prometheus/client_golang/prometheus/promhttp"

func (s *Server) setupMetrics() {
	if s.metricsEnabled {
		s.router.Get("/metrics", promhttp.Handler().ServeHTTP)
	}
}
```

---

## 📋 Checklist for New Projects

### Initial Setup
- [ ] Replace {projectname} throughout codebase
- [ ] Replace {organization} in go.mod
- [ ] Update README.md with project-specific details
- [ ] Create initial data structures
- [ ] Define API endpoints

### Core Implementation
- [ ] Implement data loading/indexing
- [ ] Create HTTP handlers
- [ ] Set up authentication
- [ ] Add admin interface
- [ ] Implement health checks

### Build & Test
- [ ] Update Makefile with correct names
- [ ] Create Dockerfile
- [ ] Write integration tests
- [ ] Add CI/CD pipeline
- [ ] Generate documentation

### Deployment
- [ ] Create installation scripts
- [ ] Test on target platforms
- [ ] Set up systemd service
- [ ] Configure monitoring
- [ ] Document deployment process

### Optional
- [ ] Add task scheduler
- [ ] Integrate GeoIP (if needed)
- [ ] Enable Prometheus metrics
- [ ] Add GraphQL endpoint
- [ ] Implement caching

---

## 📚 Additional Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Chi Router](https://github.com/go-chi/chi)
- [SQLite](https://www.sqlite.org/)
- [Docker Documentation](https://docs.docker.com/)

### Similar Projects
- Reference the airports implementation in this repository
- Adapt patterns for your specific domain
- Maintain the core architecture principles

---

**Version**: 2.0
**Last Updated**: 2024-10-14
**Based On**: Airports API Server Implementation

This specification is a living document. Update it as your project evolves while maintaining compatibility with the core template principles.
