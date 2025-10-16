# 🔢 Zipcodes API Server - Project Specification

**Project**: zipcodes
**Module**: github.com/apimgr/zipcodes
**Language**: Go 1.24+
**Purpose**: Public US postal code API with admin-protected server configuration
**Data**: 340,000+ US ZIP codes (embedded), GeoIP databases (sapics/ip-location-db)

---

## 📖 Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Directory Layout](#directory-layout)
4. [Data Sources](#data-sources)
5. [Authentication](#authentication)
6. [Routes & Endpoints](#routes--endpoints)
7. [Configuration](#configuration)
8. [Build & Deployment](#build--deployment)
9. [Development](#development)
10. [Testing](#testing)

---

## 🎯 Project Overview

### What This Is

A **public US ZIP code information API** with a web frontend, built as a single self-contained Go binary.

- **Public API**: All ZIP code data is freely accessible (no authentication)
- **Admin Interface**: Server configuration protected by token/password authentication
- **Embedded Data**: zipcodes.json (6.6MB) built into binary, GeoIP databases auto-downloaded (~103MB)
- **Fast Search**: In-memory SQLite database with indexes for instant lookups
- **Geographic Queries**: Search by coordinates, city, state, county
- **Web Frontend**: Go html/template based UI with dark theme
- **Export Formats**: JSON, CSV, plain text

### Key Features

- Search by ZIP code, city, state, county, prefix
- Find ZIP codes near coordinates
- GeoIP lookup (find ZIP codes near IP address)
- Autocomplete suggestions for cities and states
- RESTful API with matching web/API routes
- Admin dashboard for server configuration
- Single binary deployment (9.4MB binary + auto-downloaded GeoIP databases)

---

## 🏗️ Architecture

### System Design

```
┌─────────────────────────────────────────┐
│     Single Static Go Binary (9.4MB)     │
│  ┌─────────────────────────────────┐   │
│  │  Embedded Assets (go:embed)     │   │
│  │  • zipcodes.json (6.6MB)        │   │
│  │  • HTML templates               │   │
│  │  • CSS/JS/Images                │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  SQLite Database (Runtime)      │   │
│  │  • 340,000+ ZIP codes indexed   │   │
│  │  • Admin credentials (hashed)   │   │
│  │  • Server settings              │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  HTTP Server (Chi Router)       │   │
│  │  • Public routes (no auth)      │   │
│  │  • Admin routes (auth required) │   │
│  │  • API v1 endpoints             │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  GeoIP Service (Auto-Download)  │   │
│  │  • sapics/ip-location-db        │   │
│  │  • IPv4/IPv6 city databases     │   │
│  │  • Country & ASN databases      │   │
│  │  • Downloaded from jsdelivr CDN │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### Technology Stack

- **Language**: Go 1.24+
- **HTTP Router**: Chi v5
- **Database**: SQLite (github.com/mattn/go-sqlite3 - CGO enabled)
- **Templates**: Go html/template
- **GeoIP**: oschwald/geoip2-golang
- **Embedding**: Go embed.FS
- **Authentication**: SHA-256 hashing, Bearer tokens, Basic Auth

---

## 📁 Directory Layout

### OS-Specific Paths

```yaml
Linux/BSD (with root privileges):
  Config:  /etc/zipcodes/
  Data:    /var/lib/zipcodes/
  Logs:    /var/log/zipcodes/
  Runtime: /run/zipcodes/

Linux/BSD (without root):
  Config:  ~/.config/zipcodes/
  Data:    ~/.local/share/zipcodes/
  Logs:    ~/.local/state/zipcodes/
  Runtime: ~/.local/run/zipcodes/

macOS (with privileges):
  Config:  /Library/Application Support/Zipcodes/
  Data:    /Library/Application Support/Zipcodes/data/
  Logs:    /Library/Logs/Zipcodes/
  Runtime: /var/run/zipcodes/

macOS (without privileges):
  Config:  ~/Library/Application Support/Zipcodes/
  Data:    ~/Library/Application Support/Zipcodes/data/
  Logs:    ~/Library/Logs/Zipcodes/
  Runtime: ~/Library/Application Support/Zipcodes/run/

Windows:
  Config:  C:\ProgramData\Zipcodes\config\
  Data:    C:\ProgramData\Zipcodes\data\
  Logs:    C:\ProgramData\Zipcodes\logs\
  Runtime: C:\ProgramData\Zipcodes\run\

Windows (user):
  Config:  %APPDATA%\Zipcodes\config\
  Data:    %APPDATA%\Zipcodes\data\
  Logs:    %APPDATA%\Zipcodes\logs\
  Runtime: %APPDATA%\Zipcodes\run\

Docker:
  Config:  /config
  Data:    /data
  Logs:    /logs
  Runtime: /tmp
```

### Directory Contents

```yaml
Config Directory:
  - admin_credentials     # Generated on first run (0600 permissions)
  - geoip/               # GeoIP databases (auto-downloaded from jsdelivr CDN)
    - geolite2-city-ipv4.mmdb
    - geolite2-city-ipv6.mmdb
    - geo-whois-asn-country.mmdb
    - asn.mmdb

Data Directory:
  - zipcodes.db          # SQLite database (340,000+ records)
  - db/                  # Database storage directory
    - zipcodes.db        # Can also be at /data/db/zipcodes.db

Logs Directory:
  - access.log           # HTTP access logs
  - error.log            # Application errors
  - audit.log            # Admin actions

Runtime Directory:
  - zipcodes.pid         # Process ID file
  - zipcodes.sock        # Unix socket (optional)
```

### Environment Variables & Flags

```yaml
Directory Overrides (in priority order):
  1. Command-line flags
  2. Environment variables
  3. OS defaults

Flags:
  --config DIR        # Configuration directory
  --data DIR          # Data directory
  --logs DIR          # Logs directory

  --port PORT         # HTTP port (default: random 64000-64999)
  --address ADDR      # Listen address (default: 0.0.0.0)

  --db-path PATH      # SQLite database path

  --dev               # Development mode
  --version           # Show version
  --status            # Health check
  --help              # Show help

Environment Variables:
  CONFIG_DIR          # Configuration directory
  DATA_DIR            # Data directory
  LOGS_DIR            # Logs directory

  PORT                # Server port
  ADDRESS             # Listen address

  DB_PATH             # SQLite database path

  ADMIN_USER          # Admin username (first run only)
  ADMIN_PASSWORD      # Admin password (first run only)
  ADMIN_TOKEN         # Admin API token (first run only)

Docker Environment:
  Mounted volumes:
    -v ./config:/config
    -v ./data:/data

  Environment:
    -e CONFIG_DIR=/config
    -e DATA_DIR=/data
    -e PORT=80
    -e ADMIN_PASSWORD=changeme
```

### Project Source Layout

```
./
├── .claude/
│   └── settings.local.json # Claude Code settings
├── .github/
│   └── workflows/
│       └── build.yml       # GitHub Actions (build on push & monthly)
├── .gitattributes          # Git attributes
├── .gitignore              # Git ignore patterns
├── .readthedocs.yml        # ReadTheDocs configuration
├── CLAUDE.md               # This file (specification)
├── Dockerfile              # Alpine-based multi-stage build
├── docker-compose.yml      # Production compose
├── docker-compose.test.yml # Development/testing compose
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── Jenkinsfile             # CI/CD pipeline (jenkins.casjay.cc)
├── LICENSE.md              # License file
├── Makefile                # Build system (4 targets)
├── README.md               # User documentation
├── release.txt             # Version tracking (auto-increment)
├── binaries/               # Built binaries (gitignored)
│   ├── zipcodes-linux-amd64
│   ├── zipcodes-linux-arm64
│   ├── zipcodes-windows-amd64.exe
│   ├── zipcodes-windows-arm64.exe
│   ├── zipcodes-darwin-amd64
│   ├── zipcodes-darwin-arm64
│   ├── zipcodes-freebsd-amd64
│   ├── zipcodes-freebsd-arm64
│   └── zipcodes            # Host platform binary
├── releases/                # Release artifacts (gitignored)
├── rootfs/                 # Docker volumes (gitignored)
│   ├── config/
│   │   └── zipcodes/       # Service config
│   ├── data/
│   │   └── zipcodes/       # Service data
│   ├── logs/
│   │   └── zipcodes/       # Service logs
│   └── db/                 # External databases
│       ├── postgres/
│       ├── mariadb/
│       └── redis/
├── docs/                   # Documentation (ReadTheDocs)
│   ├── index.md            # Documentation home
│   ├── API.md              # Complete API reference
│   ├── SERVER.md           # Server administration guide
│   ├── README.md           # Documentation index
│   ├── mkdocs.yml          # MkDocs configuration (Dracula theme)
│   ├── requirements.txt    # Python dependencies for RTD
│   ├── stylesheets/
│   │   └── dracula.css     # Dracula theme CSS
│   ├── javascripts/
│   │   └── extra.js        # Custom JavaScript
│   └── overrides/          # MkDocs theme overrides
├── scripts/                # Production scripts (optional)
│   ├── install.sh          # Installation script
│   ├── backup.sh           # Backup script
│   └── uninstall.sh        # Uninstallation script
├── tests/                  # Test & debug scripts (optional)
│   ├── test-docker.sh      # Docker testing script
│   ├── unit/               # Unit tests
│   ├── integration/        # Integration tests
│   └── e2e/                # End-to-end tests
└── src/                    # Source code
    ├── data/
    │   └── zipcodes.json   # JSON data ONLY (6.6MB, no .go files)
    ├── admin/              # Admin authentication & handlers
    │   ├── middleware.go   # Auth middleware
    │   └── handlers.go     # Admin route handlers
    ├── api/                # API handlers
    │   └── zipcode_handlers.go # ZIP code endpoints
    ├── database/           # Database package
    │   ├── schema.go       # ZIP code schema
    │   ├── zipcode.go      # ZIP code CRUD
    │   └── admin_schema.go # Admin auth schema
    ├── geoip/              # GeoIP service package
    │   ├── geoip.go        # GeoIP lookups
    │   ├── downloader.go   # Database downloads
    │   ├── updater.go      # Auto-updates
    │   └── handlers.go     # GeoIP API handlers
    ├── paths/              # OS path detection
    │   └── paths.go        # OS-specific directory resolution
    ├── utils/              # Utility functions
    │   └── address.go      # Address utilities
    ├── server/             # HTTP server package
    │   └── server.go       # Server setup & routing
    └── main.go             # Entry point
```

---

## 💾 Data Sources

### zipcodes.json

```yaml
Location: src/data/zipcodes.json
Size: 6.6MB
Records: 340,000+ US ZIP codes
Loaded: Runtime (embedded, loaded into SQLite)
Format: JSON array

Structure:
  [
    {
      "state": "CA",
      "city": "San Francisco",
      "county": "San Francisco",
      "zip_code": 94102,
      "latitude": "37.7799",
      "longitude": "-122.4203"
    }
  ]

Database Schema (SQLite):
  CREATE TABLE zipcodes (
    zip_code INTEGER PRIMARY KEY,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    county TEXT,
    latitude TEXT,
    longitude TEXT
  );

  CREATE INDEX idx_city ON zipcodes(city);
  CREATE INDEX idx_state ON zipcodes(state);
  CREATE INDEX idx_county ON zipcodes(county);
  CREATE INDEX idx_prefix ON zipcodes(zip_code);

Query Performance:
  - Exact ZIP lookup: < 1ms
  - City search: < 10ms
  - State search: < 10ms
  - Autocomplete: < 5ms
```

### GeoIP Databases

```yaml
Source: sapics/ip-location-db
Repository: https://github.com/sapics/ip-location-db
CDN: https://cdn.jsdelivr.net/npm/@ip-location-db/
Location: {CONFIG_DIR}/geoip/*.mmdb
Auto-Download: Yes (on first run if missing)
Total Size: ~103MB
Update Frequency: Daily (via jsdelivr CDN)

Databases:
  1. geolite2-city-ipv4.mmdb (~50MB)
     - City-level geolocation for IPv4
     - Coordinates, timezone, postal codes
     - MaxMind GeoLite2 data
     - URL: https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv4.mmdb

  2. geolite2-city-ipv6.mmdb (~40MB)
     - City-level geolocation for IPv6
     - Coordinates, timezone, postal codes
     - MaxMind GeoLite2 data
     - URL: https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv6.mmdb

  3. geo-whois-asn-country.mmdb (~8MB)
     - Country-level data (combined IPv4/IPv6)
     - Aggregated from WHOIS and ASN sources
     - Public domain (CC0/PDDL)
     - Daily updates
     - URL: https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb

  4. asn.mmdb (~5MB)
     - ASN/ISP information (combined IPv4/IPv6)
     - Autonomous System Numbers
     - Daily updates
     - URL: https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb

Download:
  Auto: On first run (if missing)
  Manual: Via admin interface (/admin/geoip/update)
  Source: jsdelivr CDN (daily updates)
  Timeout: 5 minutes per file

Update Schedule:
  - Manual via admin interface
  - Auto-update available (configurable)
  - Default: Weekly (Sunday 3:00 AM)

Benefits:
  - Daily updates (vs weekly from P3TERX)
  - Separate IPv4/IPv6 databases (optimized performance)
  - Public domain country data (no attribution required)
  - Multiple data sources (MaxMind, WHOIS, ASN, GeoFeed)
  - CDN delivery (fast, global, 99.9% uptime)
```

---

## 🔐 Authentication

### Overview

This project uses **admin-only authentication** - all ZIP code data is public, only server configuration requires authentication.

**Complete guide**: [docs/SERVER.md](./docs/SERVER.md)

### Authentication Methods

```yaml
1. API Token (Bearer):
   Header: Authorization: Bearer <token>
   Use: Programmatic access to admin API
   Format: 64-character hex string
   Routes: /api/v1/admin/*

2. Basic Auth:
   Header: Authorization: Basic <base64(user:pass)>
   Use: Web UI access
   Browser: Prompts automatically
   Routes: /admin/*
```

### First Run Setup

```yaml
On first startup:
  1. Check if admin credentials exist in database

  2. If not, generate:
     - Username: $ADMIN_USER or "administrator"
     - Password: $ADMIN_PASSWORD or random 16-char
     - Token: $ADMIN_TOKEN or random 64-char hex

  3. Save to database (SHA-256 hashed)

  4. Write to {CONFIG_DIR}/admin_credentials (0600)
     Example: /etc/zipcodes/admin_credentials

  5. Display credentials in console output
     ⚠️  Shown once - save securely!

Credential File Format:
  ========================================
  ZIPCODES API - ADMIN CREDENTIALS
  ========================================
  WEB UI LOGIN:
    URL:      http://server:port/admin
    Username: administrator
    Password: <password>

  API TOKEN:
    Header:   Authorization: Bearer <token>
    Token:    <64-char-hex>

  Created: 2024-01-01 12:00:00
  ========================================
```

### Environment Variables

```yaml
First Run Only (ignored after setup):
  ADMIN_USER=admin            # Default: administrator
  ADMIN_PASSWORD=secure123    # Default: random 16-char
  ADMIN_TOKEN=abc123...       # Default: random 64-char hex

After first run:
  Credentials stored in database
  Environment variables ignored
  To reset: delete database
```

---

## 🗺️ Routes & Endpoints

### Route Matching Philosophy

**Routes must mirror between web and API:**
- `/` ↔ `/api/v1`
- `/zipcode/search` ↔ `/api/v1/zipcode/search`
- `/docs` ↔ `/api/v1/docs`
- `/admin` ↔ `/api/v1/admin`

### Public Routes (No Authentication)

```yaml
Homepage:
  GET  /                      → Home page with search interface
  GET  /api/v1                → API information JSON

Search:
  GET  /zipcode/search        → Search page
  GET  /api/v1/zipcode/search → Search ZIP codes (JSON)
    Query params:
      ?q=query               - Search term (ZIP, city, state, prefix)
      ?city=name             - Filter by city
      ?state=code            - Filter by state (2-letter)
      ?county=name           - Filter by county
      ?limit=1000            - Results limit (default: 1000)
      ?offset=0              - Pagination

ZIP Code Details:
  GET  /zipcode/:code         → ZIP code detail page
  GET  /api/v1/zipcode/:code  → ZIP code data (JSON)
  GET  /api/v1/zipcode/:code.txt → ZIP code data (plain text)

Location Search:
  GET  /zipcode/city/:city    → All ZIP codes in city
  GET  /api/v1/zipcode/city/:city → JSON

  GET  /zipcode/state/:state  → All ZIP codes in state
  GET  /api/v1/zipcode/state/:state → JSON

Autocomplete:
  GET  /api/v1/zipcode/autocomplete → Suggestions
    Query params:
      ?q=query               - Search term
      ?limit=10              - Max suggestions (default: 10, max: 50)

Statistics:
  GET  /zipcode/stats         → Stats page
  GET  /api/v1/zipcode/stats  → Database statistics (JSON)
    Returns:
      - Total ZIP codes
      - States count
      - Cities count

GeoIP:
  GET  /geoip                 → GeoIP lookup page
  GET  /api/v1/geoip          → Lookup request IP (JSON)
  GET  /api/v1/geoip.txt      → Lookup request IP (plain text)
  GET  /api/v1/geoip?ip=1.2.3.4 → Lookup specific IP (JSON)
  POST /api/v1/geoip/batch    → Batch lookup (max 100 IPs)

Export:
  GET  /api/v1/zipcodes.json  → Full database JSON (6.6MB)

Documentation:
  GET  /docs                  → API documentation page

Health:
  GET  /healthz               → Health check (JSON)

Static Assets:
  GET  /static/*              → CSS, JS, images (embedded)
  GET  /favicon.ico           → Favicon
  GET  /robots.txt            → Robots file
```

### Admin Routes (Authentication Required)

```yaml
Dashboard:
  GET  /admin                 → Admin dashboard (Basic Auth)
  GET  /api/v1/admin          → Admin info (Bearer Token)

Settings:
  GET  /admin/settings        → Settings page
  POST /admin/settings        → Update settings
  GET  /api/v1/admin/settings → Get all settings (JSON)
  PUT  /api/v1/admin/settings → Update settings (JSON)

Database:
  GET  /admin/database        → Database management page
  POST /admin/database/rebuild → Rebuild database
  GET  /api/v1/admin/database → Database status (JSON)

GeoIP Management:
  GET  /admin/geoip           → GeoIP management page
  POST /admin/geoip/update    → Update GeoIP databases
  GET  /api/v1/admin/geoip    → GeoIP status (JSON)
  POST /api/v1/admin/geoip/update → Update databases (JSON)

Logs:
  GET  /admin/logs            → Logs viewer page
  GET  /admin/logs/:type      → View specific log
  GET  /api/v1/admin/logs     → List available logs (JSON)
  GET  /api/v1/admin/logs/:type → Get log content (JSON)

Health:
  GET  /admin/health          → Server health page
  GET  /api/v1/admin/health   → Detailed health (JSON)
```

### Response Format

```yaml
JSON Success:
  {
    "success": true,
    "data": { ... },
    "count": 1
  }

JSON Error:
  {
    "success": false,
    "error": "error message"
  }

Text Format (.txt endpoints):
  Plain text, human-readable
  No JSON structure

  Example:
    ZIP Code: 94102
    City: San Francisco
    State: CA
    County: San Francisco
    Latitude: 37.7799
    Longitude: -122.4203
```

---

## ⚙️ Configuration

### Database Schema

```sql
-- ZIP codes table
CREATE TABLE IF NOT EXISTS zipcodes (
  zip_code INTEGER PRIMARY KEY,
  city TEXT NOT NULL,
  state TEXT NOT NULL,
  county TEXT,
  latitude TEXT,
  longitude TEXT
);

CREATE INDEX IF NOT EXISTS idx_city ON zipcodes(city);
CREATE INDEX IF NOT EXISTS idx_state ON zipcodes(state);
CREATE INDEX IF NOT EXISTS idx_county ON zipcodes(county);

-- Admin credentials table
CREATE TABLE IF NOT EXISTS admin_credentials (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  token_hash TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Settings table
CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('string', 'number', 'boolean', 'json')),
  category TEXT NOT NULL,
  description TEXT,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Default Settings

```yaml
Server:
  server.title: "Zipcodes API"
  server.address: "0.0.0.0"
  server.http_port: 64000-64999 (random)
  server.https_enabled: false

Database:
  db.path: "{DATA_DIR}/zipcodes.db"
  db.auto_rebuild: false

GeoIP:
  geoip.enabled: true
  geoip.auto_update: false
  geoip.update_schedule: "0 3 * * 0" # Sunday 3 AM

Logging:
  log.level: "info"
  log.format: "json"
  log.access: true

Features:
  features.autocomplete: true
  features.batch_geoip: true
  features.export: true
```

### Modifying Settings

```yaml
Web UI:
  1. Navigate to /admin/settings
  2. Change values in form
  3. Click Save (applies immediately)

API:
  PUT /api/v1/admin/settings
  {
    "settings": {
      "server.title": "My ZIP Code API",
      "geoip.auto_update": true
    }
  }

Environment (first run only):
  DB_PATH=/data/zipcodes.db
  PORT=8080
```

---

## 🔨 Build & Deployment

### Makefile Targets

```makefile
Targets:
  make build             # Build all platforms
  make test              # Run tests
  make docker            # Build and push multi-arch Docker images
  make docker-dev        # Build development Docker image (local only)
  make release           # Create GitHub release
  make clean             # Remove build artifacts

Build Process:
  1. go mod download
  2. go build for all platforms:
     - Linux: amd64, arm64
     - Windows: amd64, arm64
     - macOS: amd64, arm64 (Apple Silicon)
     - FreeBSD: amd64, arm64
  3. Create binaries/ directory with outputs
  4. Auto-increment version in release.txt

Platforms:
  binaries/zipcodes-linux-amd64
  binaries/zipcodes-linux-arm64
  binaries/zipcodes-windows-amd64.exe
  binaries/zipcodes-windows-arm64.exe
  binaries/zipcodes-darwin-amd64
  binaries/zipcodes-darwin-arm64
  binaries/zipcodes-freebsd-amd64
  binaries/zipcodes-freebsd-arm64
  binaries/zipcodes              # Current platform
```

### Docker

```yaml
Dockerfile:
  Multi-stage build (Go builder → Alpine runtime)
  CGO_ENABLED=0 for static binary
  Binary Size: 9.4MB static binary
  Runtime Tools: curl, bash, ca-certificates, tzdata
  GeoIP Databases: Auto-downloaded (~103MB from jsdelivr CDN)
  Health check: /healthz endpoint via --status flag
  Volumes: /config, /data, /logs
  User: 65534:65534 (nobody)
  Exposed port: 80

Building:
  docker build -t zipcodes:latest .

  With version:
    docker build \
      --build-arg VERSION=1.0.0 \
      --build-arg COMMIT=$(git rev-parse --short HEAD) \
      --build-arg BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
      -t zipcodes:1.0.0 .

Production Deployment:
  Uses docker-compose.yml with ./rootfs for persistent storage

  Start:
    docker-compose up -d

  Volumes mounted to ./rootfs:
    - ./rootfs/config/zipcodes → /config (in container)
    - ./rootfs/data/zipcodes → /data (in container)
    - ./rootfs/db/sqlite → /data/db (in container)

  Default configuration:
    - Internal port: 80 (Docker)
    - External port: 64080
    - Bind address: 172.17.0.1:64080:80

  For public access:
    Change port mapping in docker-compose.yml:
      - "64080:80"      # Public HTTP

  Access:
    http://172.17.0.1:64080         # Homepage
    http://172.17.0.1:64080/admin   # Admin UI (Basic Auth)
    http://172.17.0.1:64080/api/v1  # API endpoints

  Check credentials:
    cat ./rootfs/config/zipcodes/admin_credentials

  View logs:
    docker-compose logs -f zipcodes
    cat ./rootfs/logs/zipcodes/access.log

  Set admin credentials (first run):
    Edit docker-compose.yml environment:
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=strong-password

  Stop:
    docker-compose down

Testing/Debugging:
  Uses docker-compose.test.yml with /tmp for ephemeral storage

  Test:
    cd tests
    ./test-docker.sh

  Or manually:
    docker-compose -f docker-compose.test.yml up -d

  Volumes in /tmp/zipcodes/rootfs (automatically cleaned):
    - /tmp/zipcodes/rootfs/config/zipcodes → /config
    - /tmp/zipcodes/rootfs/data/zipcodes → /data
    - /tmp/zipcodes/rootfs/logs/zipcodes → /logs

  Access:
    http://localhost:64181         # Test server

  Cleanup:
    docker-compose -f docker-compose.test.yml down
    sudo rm -rf /tmp/zipcodes/rootfs

Docker Run (Manual):
  # Production (with ./rootfs)
  docker run -d \
    --name zipcodes \
    -p 172.17.0.1:64080:80 \
    -v $(pwd)/rootfs/config/zipcodes:/config \
    -v $(pwd)/rootfs/data/zipcodes:/data \
    -v $(pwd)/rootfs/logs/zipcodes:/logs \
    -e PORT=80 \
    -e ADMIN_PASSWORD=changeme \
    --restart unless-stopped \
    ghcr.io/apimgr/zipcodes:latest

  # Testing (with /tmp)
  docker run -d \
    --name zipcodes-test \
    -p 127.0.0.1:64181:80 \
    -v /tmp/zipcodes/rootfs/config/zipcodes:/config \
    -v /tmp/zipcodes/rootfs/data/zipcodes:/data \
    -v /tmp/zipcodes/rootfs/db/sqlite:/data/db \
    -e PORT=80 \
    -e ADMIN_PASSWORD=testpass \
    ghcr.io/apimgr/zipcodes:latest
```

### Manual Installation

```bash
# Download binary
wget https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-linux-amd64
chmod +x zipcodes-linux-amd64
sudo mv zipcodes-linux-amd64 /usr/local/bin/zipcodes

# Create directories (as root)
sudo mkdir -p /etc/zipcodes /var/lib/zipcodes /var/log/zipcodes

# Run
sudo zipcodes --port 80

# Or run as user (random port)
zipcodes
# Check output for port and credentials
```

### CI/CD

```yaml
GitHub Actions (.github/workflows/build.yml):
  Triggers:
    - Push to main branch
    - Monthly schedule (1st at 3 AM UTC)

  Actions:
    - Build all platform binaries
    - Create multi-arch Docker images (amd64/arm64)
    - Publish to ghcr.io/apimgr/zipcodes
    - Create GitHub releases

  Platforms:
    - Linux: amd64, arm64
    - macOS: amd64, arm64
    - Windows: amd64, arm64
    - FreeBSD: amd64, arm64

Jenkins Pipeline (Jenkinsfile):
  Server: jenkins.casjay.cc

  Stages:
    - Build multi-architecture (amd64/arm64)
    - Parallel testing on both architectures
    - Docker image publishing
    - GitHub releases

  Features:
    - Multi-stage parallel builds
    - Cross-platform testing
    - Automated releases

ReadTheDocs (.readthedocs.yml):
  Documentation: https://zipcodes.readthedocs.io

  Build:
    - MkDocs with Material theme
    - Dracula color scheme
    - Custom CSS/JS

  Pages:
    - index.md - Home
    - API.md - Complete API reference
    - SERVER.md - Server administration
```

---

## 🛠️ Development

### Development Mode

```yaml
Enable:
  --dev flag
  OR DEV=true environment variable

Features:
  - Hot reload templates (no restart)
  - Detailed logging (SQL queries, stack traces)
  - Debug endpoints enabled
  - CORS allow all origins
  - Fast session timeout (5 min)

Debug Endpoints:
  GET /debug/routes          - List all routes
  GET /debug/config          - Show configuration
  GET /debug/db              - Database stats
  GET /debug/zipcodes        - ZIP code stats
  POST /debug/reset          - Reset to fresh state
```

### Local Development

```bash
# Install dependencies
go mod download

# Run with hot reload
go run ./src --dev --port 8080

# Or with Makefile
make build
./binaries/zipcodes --dev --port 8080

# Server starts on http://localhost:8080
# Admin credentials displayed in console
```

### Database Management

```yaml
Initial Load:
  - On first run, loads zipcodes.json into SQLite
  - Creates indexes for fast lookups
  - Takes ~5 seconds for 340,000 records

Rebuild:
  - Via admin interface: /admin/database
  - Drops and recreates database
  - Reloads from embedded zipcodes.json

Performance:
  - Indexed queries: < 10ms
  - Full table scan: avoided
  - Database size: ~15MB
```

---

## ✅ Testing

### Test Structure

```
tests/
├── unit/
│   ├── database_test.go       # Database tests
│   └── geoip_test.go          # GeoIP service tests
├── integration/
│   ├── api_test.go            # API endpoint tests
│   └── admin_test.go          # Admin auth tests
└── e2e/
    └── scenarios_test.go      # End-to-end tests
```

### Running Tests

```bash
# All tests
make test

# Or manually
go test -v -race ./...

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmarks
go test -v -bench=. -benchmem ./...
```

### Test Coverage Requirements

```yaml
Minimum Coverage: 80%

Critical Paths (100% coverage):
  - Admin authentication
  - Database initialization
  - ZIP code search/indexing
  - GeoIP lookups
  - Autocomplete
```

---

## 📊 Performance

### Benchmarks

```yaml
Search Performance:
  Exact ZIP lookup: < 1ms
  City search: < 10ms
  State search: < 10ms
  Autocomplete: < 5ms
  Prefix search: < 10ms

Database:
  Records: 340,000+
  Size: ~15MB (SQLite)
  Indexes: ZIP code, city, state, county
  Query cache: Enabled

Memory:
  Baseline: ~50MB
  With GeoIP: ~150MB
  Peak: ~200MB

Throughput:
  Sequential: 1,000+ req/s
  Concurrent: 5,000+ req/s (with caching)
```

---

## 🔒 Security

### Best Practices

```yaml
Credentials:
  - Change default admin password immediately
  - Rotate API tokens periodically
  - Use HTTPS in production
  - Restrict admin routes to internal network

File Permissions:
  admin_credentials: 0600 (owner read/write only)
  Database: 0644 (owner write, all read)
  Logs: 0644

Network:
  - Bind to 0.0.0.0 only if needed
  - Use 127.0.0.1 for local-only access
  - Configure firewall rules
  - Use reverse proxy (nginx/Caddy) for HTTPS

Database:
  - Passwords hashed with SHA-256
  - Tokens hashed with SHA-256
  - SQL injection protection (prepared statements)
  - Input validation on all endpoints
```

---

## 📝 License

MIT License - See LICENSE file

### Embedded Data Licenses

```yaml
zipcodes.json:
  Source: US Postal Service data
  License: Public Domain
  Records: 340,000+ US ZIP codes

GeoIP Databases (sapics/ip-location-db):
  Source: https://github.com/sapics/ip-location-db
  CDN: https://cdn.jsdelivr.net/npm/@ip-location-db/
  License:
    - geolite2-city (IPv4/IPv6): CC BY-SA 4.0 (MaxMind GeoLite2)
    - geo-whois-asn-country: CC0/PDDL (Public domain)
    - asn: Various open sources
  Attribution:
    - MaxMind GeoLite2 (city databases)
    - Aggregated via sapics/ip-location-db
```

---

**Zipcodes API Server v1.0** - A focused, production-ready US postal code information API with admin-only authentication. Built for simplicity, performance, and ease of deployment.
