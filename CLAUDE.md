# 🔢 Zipcodes API Server - Project Specification

**Project**: zipcodes
**Module**: github.com/apimgr/zipcodes
**Language**: Go 1.24+
**Purpose**: Public US postal code API with admin-protected server configuration
**Data**: 340,000+ US ZIP codes (embedded), GeoIP databases (sapics/ip-location-db via jsdelivr CDN)

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
11. [Performance](#performance)
12. [Security](#security)
13. [License](#license)

---

## 🎯 Project Overview

### What This Is

A **public US ZIP code information API** with a web frontend, built as a single self-contained Go binary.

- **Public API**: All ZIP code data is freely accessible (no authentication)
- **Admin Interface**: Server configuration protected by token/password authentication
- **Embedded Data**: zipcodes.json (6.4MB) built into binary
- **External GeoIP**: ~103MB GeoIP databases auto-downloaded on first run
- **Fast Search**: In-memory SQLite database with indexes for instant lookups
- **Geographic Queries**: Search by coordinates, city, state, county
- **Web Frontend**: Go html/template based UI with dark theme
- **Export Formats**: JSON, CSV, plain text

### Key Features

- Search by ZIP code, city, state, county, prefix
- Find ZIP codes near coordinates
- GeoIP lookup (find ZIP codes near IP address) with IPv4/IPv6 support
- Autocomplete suggestions for cities and states
- RESTful API with matching web/API routes
- Admin dashboard for server configuration
- Single binary deployment (~16MB with embedded data)
- Multi-platform support (Linux, macOS, Windows, FreeBSD × amd64/arm64)

### Binary Characteristics

```yaml
Binary Size: ~16MB (static binary with CGO disabled)
Embedded Assets:
  - zipcodes.json: 6.4MB (340,000+ records)
  - HTML templates
  - CSS/JS/Images
  - Static files

External Assets (auto-downloaded):
  - GeoIP databases: ~103MB total
  - Downloaded from jsdelivr CDN on first run
  - Stored in DATA_DIR/geoip/

Total Runtime Footprint:
  - Binary: 16MB
  - SQLite DB: ~15MB (generated from embedded JSON)
  - GeoIP: ~103MB (optional, auto-downloaded)
  - Memory: ~50MB baseline, ~150MB with GeoIP
```

---

## 🏗️ Architecture

### System Design

```
┌─────────────────────────────────────────┐
│         Single Go Binary (~16MB)        │
│  ┌─────────────────────────────────┐   │
│  │  Embedded Assets (go:embed)     │   │
│  │  • zipcodes.json (6.4MB)        │   │
│  │  • HTML templates               │   │
│  │  • CSS/JS/Images                │   │
│  └─────────────────────────────────┘   │
│  ┌─────────────────────────────────┐   │
│  │  SQLite Database (runtime)      │   │
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
└─────────────────────────────────────────┘
            ↓ Auto-downloads on first run
┌─────────────────────────────────────────┐
│    External GeoIP Databases (~103MB)    │
│  ┌─────────────────────────────────┐   │
│  │  sapics/ip-location-db          │   │
│  │  (via jsdelivr CDN)             │   │
│  │  • geolite2-city-ipv4.mmdb      │   │
│  │  • geolite2-city-ipv6.mmdb      │   │
│  │  • geo-whois-asn-country.mmdb   │   │
│  │  • asn.mmdb                     │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### Technology Stack

- **Language**: Go 1.24+
- **HTTP Router**: Chi v5
- **Database**: SQLite (github.com/mattn/go-sqlite3)
- **Templates**: Go html/template
- **GeoIP**: oschwald/geoip2-golang
- **Embedding**: Go embed.FS
- **Authentication**: SHA-256 hashing, Bearer tokens, Basic Auth

### Build Configuration

```yaml
CGO: Disabled (CGO_ENABLED=0)
Static Binary: Yes
LDFLAGS: -w -s (strip debug info)
Platforms:
  - Linux: amd64, arm64
  - macOS: amd64, arm64 (Apple Silicon)
  - Windows: amd64, arm64
  - FreeBSD: amd64, arm64
```

---

## 📁 Directory Layout

### OS-Specific Paths

```yaml
Linux/BSD (with root privileges):
  Config:  /etc/zipcodes/
  Data:    /var/lib/zipcodes/
  Logs:    /var/log/zipcodes/

Linux/BSD (without root):
  Config:  ~/.config/zipcodes/
  Data:    ~/.local/share/zipcodes/
  Logs:    ~/.local/state/zipcodes/

macOS (with privileges):
  Config:  /Library/Application Support/Zipcodes/
  Data:    /Library/Application Support/Zipcodes/data/
  Logs:    /Library/Logs/Zipcodes/

macOS (without privileges):
  Config:  ~/Library/Application Support/Zipcodes/
  Data:    ~/Library/Application Support/Zipcodes/data/
  Logs:    ~/Library/Logs/Zipcodes/

Windows:
  Config:  C:\ProgramData\Zipcodes\config\
  Data:    C:\ProgramData\Zipcodes\data\
  Logs:    C:\ProgramData\Zipcodes\logs\

Windows (user):
  Config:  %APPDATA%\Zipcodes\
  Data:    %LOCALAPPDATA%\Zipcodes\
  Logs:    %LOCALAPPDATA%\Zipcodes\logs\

Docker:
  Config:  /config
  Data:    /data
  Logs:    /logs
```

### Directory Contents

```yaml
Config Directory:
  - admin_credentials     # Generated on first run (0600 permissions)

Data Directory:
  - zipcodes.db          # SQLite database (340,000+ records)
  - db/                  # Database storage directory (optional)
    - zipcodes.db        # Alternative location: /data/db/zipcodes.db
  - geoip/               # GeoIP databases (auto-downloaded)
    - geolite2-city-ipv4.mmdb       (~50MB)
    - geolite2-city-ipv6.mmdb       (~40MB)
    - geo-whois-asn-country.mmdb    (~8MB)
    - asn.mmdb                      (~5MB)

Logs Directory:
  - access.log           # HTTP access logs (future)
  - error.log            # Application errors (future)
  - audit.log            # Admin actions (future)
```

### Environment Variables & Flags

```yaml
Directory Overrides (in priority order):
  1. Command-line flags (highest priority)
  2. Environment variables
  3. OS defaults (lowest priority)

Command-line Flags:
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
    -v ./rootfs/config/zipcodes:/config
    -v ./rootfs/data/zipcodes:/data
    -v ./rootfs/logs/zipcodes:/logs

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
│       └── release.yml     # GitHub Actions (build on push & monthly)
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
├── LICENSE.md              # MIT License
├── Makefile                # Build system
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
├── releases/               # Release artifacts (gitignored)
├── rootfs/                 # Docker volumes (gitignored)
│   ├── config/
│   │   └── zipcodes/       # Service config
│   ├── data/
│   │   └── zipcodes/       # Service data
│   └── logs/
│       └── zipcodes/       # Service logs
├── docs/                   # Documentation (MkDocs)
│   ├── index.md            # Documentation home
│   └── mkdocs.yml          # MkDocs configuration
├── scripts/                # Production scripts (optional)
│   └── install.sh          # Installation helper
├── tests/                  # Test scripts (optional)
│   └── test-docker.sh      # Docker testing
└── src/                    # Source code
    ├── data/
    │   └── zipcodes.json   # JSON data ONLY (6.4MB, no .go files)
    ├── admin/              # Admin authentication & handlers
    │   ├── middleware.go   # Auth middleware
    │   └── handlers.go     # Admin route handlers
    ├── api/                # API handlers
    │   └── zipcode_handlers.go # ZIP code endpoints
    ├── database/           # Database package
    │   ├── schema.go       # Universal schema (users, sessions, etc.)
    │   ├── zipcode.go      # ZIP code CRUD operations
    │   └── admin_schema.go # Admin-only auth schema
    ├── geoip/              # GeoIP service package
    │   ├── geoip.go        # GeoIP lookups
    │   ├── downloader.go   # Database downloads (sapics via jsdelivr)
    │   ├── updater.go      # Auto-updates
    │   └── handlers.go     # GeoIP API handlers
    ├── paths/              # OS path detection
    │   └── paths.go        # OS-specific directory resolution
    ├── utils/              # Utility functions
    │   └── address.go      # Address utilities (GetDisplayAddress)
    ├── server/             # HTTP server package
    │   ├── server.go       # Server setup & routing
    │   ├── docs_handlers.go # OpenAPI/GraphQL handlers
    │   ├── static/         # Static assets (embedded)
    │   └── templates/      # HTML templates (embedded)
    │       └── index.html
    └── main.go             # Entry point
```

---

## 💾 Data Sources

### zipcodes.json

```yaml
Location: src/data/zipcodes.json
Size: 6.4MB
Records: 340,000+ US ZIP codes
Lines: 341,930
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
Distribution: jsdelivr CDN (daily updates)
Location: {DATA_DIR}/geoip/*.mmdb
Auto-Download: Yes (on first run if missing)
Total Size: ~103MB

Databases:
  1. geolite2-city-ipv4.mmdb (~50MB)
     URL: https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv4.mmdb
     - IPv4 city-level geolocation
     - Coordinates, timezone, city names

  2. geolite2-city-ipv6.mmdb (~40MB)
     URL: https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv6.mmdb
     - IPv6 city-level geolocation
     - Coordinates, timezone, city names

  3. geo-whois-asn-country.mmdb (~8MB)
     URL: https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb
     - Country-level fallback
     - Combined IPv4/IPv6

  4. asn.mmdb (~5MB)
     URL: https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb
     - ASN information
     - ISP organization names

Download:
  Auto: On first run (if missing)
  Manual: Via admin interface (future)
  CDN: jsdelivr.net (GitHub-backed, daily updates)
  Timeout: 300 seconds (5 minutes) per file

Features:
  - Separate IPv4/IPv6 databases for better performance
  - Automatic IP version detection
  - Fallback from city -> country
  - ASN lookups for ISP information

Update Schedule:
  - Manual via admin interface (future)
  - Auto-update available (configurable, future)
  - Daily updates available from CDN
```

---

## 🔐 Authentication

### Overview

This project uses **admin-only authentication** - all ZIP code data is public, only server configuration requires authentication.

**No user accounts, no user registration, admin-only access.**

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
     - Password: $ADMIN_PASSWORD or random 16-char hex
     - Token: $ADMIN_TOKEN or random 64-char hex

  3. Save to database (SHA-256 hashed)

  4. Write to {CONFIG_DIR}/admin_credentials (0600)
     Example: ~/.config/zipcodes/admin_credentials

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
    URL:      http://server:port/api/v1/admin
    Header:   Authorization: Bearer <token>
    Token:    <64-char-hex>

  Created: 2024-01-01 12:00:00
  ========================================
```

### Environment Variables

```yaml
First Run Only (ignored after setup):
  ADMIN_USER=admin            # Default: administrator
  ADMIN_PASSWORD=secure123    # Default: random 16-char hex
  ADMIN_TOKEN=abc123...       # Default: random 64-char hex

After first run:
  Credentials stored in database (SHA-256 hashed)
  Environment variables ignored
  To reset: delete database and restart
```

### Security Implementation

```yaml
Password Hashing:
  Algorithm: SHA-256
  Storage: Hex-encoded hash
  Function: crypto/sha256

Token Hashing:
  Algorithm: SHA-256
  Storage: Hex-encoded hash
  Function: crypto/sha256

Database Schema:
  CREATE TABLE admin_credentials (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );

Verification:
  - Compare SHA-256 hash of provided credentials
  - Single admin account only (id=1)
  - No user registration
```

---

## 🗺️ Routes & Endpoints

### Route Matching Philosophy

**Routes mirror between web and API:**
- `/` ↔ `/api/v1`
- `/zipcode/search` ↔ `/api/v1/zipcode/search`
- `/admin` ↔ `/api/v1/admin`

### Public Routes (No Authentication)

```yaml
Homepage:
  GET  /                      → Home page with search interface
  GET  /api/v1                → API information JSON (future)

Search:
  GET  /zipcode/search        → Search page (future)
  GET  /api/v1/zipcode/search → Search ZIP codes (JSON)
    Query params:
      ?q=query               - Search term (ZIP, city, state, prefix)
      ?city=name             - Filter by city
      ?state=code            - Filter by state (2-letter)
      ?county=name           - Filter by county
      ?limit=1000            - Results limit (default: 1000)
      ?offset=0              - Pagination

ZIP Code Details:
  GET  /zipcode/:code         → ZIP code detail page (future)
  GET  /api/v1/zipcode/:code  → ZIP code data (JSON)
  GET  /api/v1/zipcode/:code.txt → ZIP code data (plain text)

Location Search:
  GET  /zipcode/city/:city    → All ZIP codes in city (future)
  GET  /api/v1/zipcode/city/:city → JSON

  GET  /zipcode/state/:state  → All ZIP codes in state (future)
  GET  /api/v1/zipcode/state/:state → JSON

Autocomplete:
  GET  /api/v1/zipcode/autocomplete → Suggestions
    Query params:
      ?q=query               - Search term
      ?limit=10              - Max suggestions (default: 10, max: 50)

Statistics:
  GET  /zipcode/stats         → Stats page (future)
  GET  /api/v1/zipcode/stats  → Database statistics (JSON)
    Returns:
      - Total ZIP codes
      - States count
      - Cities count

GeoIP:
  GET  /geoip                 → GeoIP lookup page (future)
  GET  /api/v1/geoip          → Lookup request IP (JSON)
  GET  /api/v1/geoip.txt      → Lookup request IP (plain text)
  GET  /api/v1/geoip?ip=1.2.3.4 → Lookup specific IP (JSON)
  POST /api/v1/geoip/batch    → Batch lookup (max 100 IPs)

Export:
  GET  /api/v1/zipcodes.json  → Full database JSON (6.4MB, embedded file)

Documentation:
  GET  /openapi               → OpenAPI/Swagger UI (future)
  GET  /graphql               → GraphQL Playground (future)
  GET  /api/v1/openapi        → OpenAPI spec (future)
  GET  /api/v1/openapi.json   → OpenAPI JSON spec (future)
  GET  /api/v1/graphql        → GraphQL endpoint (future)
  POST /api/v1/graphql        → GraphQL queries (future)

Health:
  GET  /healthz               → Health check (JSON)
  GET  /api/v1/health         → Health check (JSON)

Static Assets:
  GET  /static/*              → CSS, JS, images (embedded)
```

### Admin Routes (Authentication Required)

```yaml
Dashboard:
  GET  /admin                 → Admin dashboard (Basic Auth)
  GET  /api/v1/admin          → Admin info (Bearer Token)

Settings:
  GET  /admin/settings        → Settings page (Basic Auth)
  POST /admin/settings        → Update settings (Basic Auth)
  GET  /api/v1/admin/settings → Get all settings (Bearer Token)
  PUT  /api/v1/admin/settings → Update settings (Bearer Token)

Database:
  GET  /admin/database        → Database management page (Basic Auth)
  POST /admin/database/test   → Test database connection (Basic Auth)

Logs:
  GET  /admin/logs            → Logs viewer page (Basic Auth)

Audit:
  GET  /admin/audit           → Audit log viewer (Basic Auth)

System:
  POST /api/v1/admin/reload   → Reload configuration (Bearer Token)
  GET  /api/v1/admin/stats    → Admin statistics (Bearer Token)
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
-- ZIP codes table (zipcode data)
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

-- Admin credentials table (admin-only auth)
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
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Audit log table
CREATE TABLE IF NOT EXISTS audit_log (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
  username TEXT,
  action TEXT NOT NULL,
  resource TEXT NOT NULL,
  old_value TEXT,
  new_value TEXT,
  ip_address TEXT NOT NULL,
  user_agent TEXT,
  success INTEGER NOT NULL,
  error_message TEXT,
  timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Scheduled tasks table
CREATE TABLE IF NOT EXISTS scheduled_tasks (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
  name TEXT UNIQUE NOT NULL,
  cron_expression TEXT NOT NULL,
  command TEXT NOT NULL,
  enabled INTEGER DEFAULT 1,
  last_run DATETIME,
  next_run DATETIME NOT NULL,
  last_status TEXT,
  last_error TEXT,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Universal schema tables (currently unused, for future extensibility)
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  display_name TEXT,
  avatar_url TEXT,
  bio TEXT,
  role TEXT NOT NULL CHECK (role IN ('administrator', 'user', 'guest')) DEFAULT 'user',
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'pending')),
  timezone TEXT DEFAULT 'UTC',
  language TEXT DEFAULT 'en',
  theme TEXT DEFAULT 'dark',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_login DATETIME,
  failed_login_attempts INTEGER DEFAULT 0,
  locked_until DATETIME,
  metadata TEXT
);

CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token TEXT UNIQUE NOT NULL,
  ip_address TEXT NOT NULL,
  user_agent TEXT,
  device_name TEXT,
  expires_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_activity DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  remember_me INTEGER DEFAULT 0,
  is_active INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS tokens (
  id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  token_hash TEXT UNIQUE NOT NULL,
  last_used DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  revoked_at DATETIME
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_audit_log_timestamp ON audit_log(timestamp);
CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_tokens_token_hash ON tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_audit_log_user_id ON audit_log(user_id);
```

### Default Settings

```yaml
Server:
  server.title: "Zipcodes"
  server.tagline: "US Postal Code Lookup API"
  server.description: "Fast and accurate US zipcode lookup API with 340,000+ zipcodes, GeoIP integration, and modern web interface."
  server.address: "0.0.0.0"
  server.http_port: 64080 (default in settings, random 64000-64999 at runtime)
  server.https_enabled: false
  server.timezone: "UTC"
  server.date_format: "US"
  server.time_format: "12-hour"

Proxy:
  proxy.enabled: true
  proxy.trust_headers: true

Features:
  features.api_enabled: true

Database:
  db.path: "{DATA_DIR}/zipcodes.db"

GeoIP:
  geoip.enabled: true
  geoip.auto_update: false (future)
  geoip.update_schedule: "0 3 * * 0" # Sunday 3 AM (future)
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

Build Flags:
  CGO_ENABLED=0              # Static binary
  LDFLAGS:
    -X main.Version=$(VERSION)
    -X main.Commit=$(COMMIT)
    -X main.BuildDate=$(BUILD_DATE)
    -w -s                    # Strip debug info
```

### Docker

```yaml
Dockerfile:
  Multi-stage build (Go builder → Alpine runtime)
  CGO_ENABLED=0 for static binary
  Size: ~16MB binary in ~30MB container
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
    - ./rootfs/logs/zipcodes → /logs (in container)

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

  Set admin credentials (first run):
    Edit docker-compose.yml environment:
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=strong-password

  Stop:
    docker-compose down

Testing/Debugging:
  Uses docker-compose.test.yml with /tmp for ephemeral storage

  Test:
    docker-compose -f docker-compose.test.yml up -d

  Volumes in /tmp (automatically cleaned):
    - /tmp/zipcodes/rootfs/config → /config
    - /tmp/zipcodes/rootfs/data → /data
    - /tmp/zipcodes/rootfs/logs → /logs

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
    -v /tmp/zipcodes/rootfs/config:/config \
    -v /tmp/zipcodes/rootfs/data:/data \
    -v /tmp/zipcodes/rootfs/logs:/logs \
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
GitHub Actions (.github/workflows/release.yml):
  Triggers:
    - Push to main branch
    - Monthly schedule (1st at 3 AM UTC)
    - Manual workflow dispatch

  Actions:
    - Run tests
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
    - Custom CSS/JS

  Pages:
    - index.md - Home
```

---

## 🛠️ Development

### Development Mode

```yaml
Enable:
  --dev flag
  OR DEV=true environment variable

Features (future):
  - Hot reload templates (no restart)
  - Detailed logging (SQL queries, stack traces)
  - Debug endpoints enabled
  - CORS allow all origins
  - Fast session timeout (5 min)
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

Rebuild (future):
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

Binary:
  Size: ~16MB (static, stripped)
  Embedded: ~6.4MB zipcodes.json
  Total runtime: ~16MB binary + ~15MB DB + ~103MB GeoIP (optional)

Throughput:
  Sequential: 1,000+ req/s
  Concurrent: 5,000+ req/s (future, with caching)
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

CORS:
  - Enabled for all origins (public API)
  - Allow methods: GET, POST, PUT, DELETE, OPTIONS
  - Allow headers: Content-Type, Authorization
```

### Security Headers

```yaml
HTTP Headers:
  X-Frame-Options: DENY
  X-Content-Type-Options: nosniff
  X-XSS-Protection: 1; mode=block
  Referrer-Policy: strict-origin-when-cross-origin
  Access-Control-Allow-Origin: * (public API)
```

---

## 📝 License

MIT License - See LICENSE.md for details

### Embedded Data Licenses

```yaml
zipcodes.json:
  Source: US Postal Service data
  License: Public Domain
  Records: 340,000+ US ZIP codes

GeoIP Databases (sapics/ip-location-db):
  Source: https://github.com/sapics/ip-location-db
  Distribution: jsdelivr CDN
  License: CC BY-SA 4.0 (database)
  Attribution: MaxMind GeoLite2 + aggregated sources
  Updates: Daily via CDN

Go Dependencies:
  github.com/go-chi/chi/v5: MIT
  github.com/mattn/go-sqlite3: MIT
  github.com/oschwald/geoip2-golang: ISC
  golang.org/x/crypto: BSD-3-Clause
```

---

**Zipcodes API Server v0.0.1** - A focused, production-ready US postal code information API with admin-only authentication. Built for simplicity, performance, and ease of deployment.

**Key Features:**
- Single 16MB static binary
- 340,000+ US ZIP codes embedded
- GeoIP integration (~103MB, auto-downloaded)
- Admin-only authentication
- Multi-platform support
- Production-ready Docker deployment
