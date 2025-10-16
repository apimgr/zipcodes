# Zipcodes - US Postal Code API

Fast and accurate US zipcode lookup API with 340,000+ zipcodes, GeoIP integration, and modern web interface.

**Single static binary with embedded assets - no dependencies required.**

## About

Zipcodes provides a complete REST API and web interface for US postal code lookups with:
- 340,000+ US zipcodes with city, state, county, and coordinates
- GeoIP integration for IP-to-location lookups
- Fast indexed SQLite database (< 10ms queries)
- Modern web interface with autocomplete
- Single 9.4MB static binary with all assets embedded
- Admin-only authentication system

### Features

- **340,000+ US Zipcodes** with city, state, county, coordinates
- **Fast Search** - Indexed database, < 10ms queries
- **Multiple Search Types** - By zipcode, city, state, prefix
- **Autocomplete** - Smart suggestions as you type
- **GeoIP Integration** - IP to location lookups
- **REST API** - Full JSON API with text alternatives
- **Admin Authentication** - Secure admin panel with auto-generated credentials
- **Web Interface** - Modern, responsive, dark/light themes
- **Download Dataset** - Complete JSON file available
- **CORS Enabled** - Use from any domain
- **Single Binary** - 9.4MB static binary, no dependencies, all assets embedded
- **Multi-Platform** - Linux, macOS, Windows, FreeBSD (amd64/arm64)

### Data Sources

- **Zipcodes**: US Postal Service data (340,000+ records)
- **GeoIP**: MaxMind GeoLite2 + aggregated sources from [sapics/ip-location-db](https://github.com/sapics/ip-location-db)
  - Auto-downloads on first run (~100MB total, 4 databases)
  - Daily updates available via jsdelivr CDN
  - Separate IPv4/IPv6 city databases for better performance

## Production Installation

### Binary Installation

```bash
# Download latest release (choose your platform)
wget https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-linux-amd64

# Make executable
chmod +x zipcodes-linux-amd64
mv zipcodes-linux-amd64 /usr/local/bin/zipcodes

# Run as systemd service (recommended)
sudo nano /etc/systemd/system/zipcodes.service
```

**systemd service:**
```ini
[Unit]
Description=Zipcodes API Server
After=network.target

[Service]
Type=simple
User=zipcodes
Group=zipcodes
ExecStart=/usr/local/bin/zipcodes --port 8080
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start
sudo systemctl enable zipcodes
sudo systemctl start zipcodes
```

### Available Binaries

All releases include binaries for:
- `zipcodes-linux-amd64` - Linux (Intel/AMD)
- `zipcodes-linux-arm64` - Linux (ARM)
- `zipcodes-darwin-amd64` - macOS (Intel)
- `zipcodes-darwin-arm64` - macOS (Apple Silicon)
- `zipcodes-windows-amd64.exe` - Windows (Intel/AMD)
- `zipcodes-windows-arm64.exe` - Windows (ARM)
- `zipcodes-freebsd-amd64` - FreeBSD (Intel/AMD)
- `zipcodes-freebsd-arm64` - FreeBSD (ARM)

### First Run

On first run, admin credentials are auto-generated and saved to:
- File: `{CONFIG_DIR}/admin_credentials`
- Console output (shown once)

**Example output:**
```
========================================
ZIPCODES API - ADMIN CREDENTIALS
========================================
WEB UI LOGIN:
  URL:      http://your-server:8080/admin
  Username: administrator
  Password: <generated>

API TOKEN:
  Token:    <generated>
========================================
```

**Save these credentials immediately - they won't be shown again!**

### Configuration

#### Command Line Options

```bash
--help            Show help message
--version         Show version information
--status          Check server status
--port PORT       Set port (default: random 64000-64999)
--address ADDR    Listen address (default: 0.0.0.0)
--config DIR      Set config directory
--data DIR        Set data directory
--logs DIR        Set logs directory
--db-path PATH    Set SQLite database path
--dev             Development mode
```

#### Environment Variables

```bash
CONFIG_DIR        Configuration directory
DATA_DIR          Data directory
LOGS_DIR          Logs directory
DB_PATH           SQLite database path
PORT              Server port
ADDRESS           Listen address
ADMIN_USER        Admin username (first run only)
ADMIN_PASSWORD    Admin password (first run only)
ADMIN_TOKEN       Admin API token (first run only)
```

#### Data Storage

**Default Locations:**

Linux/BSD (user):
```
~/.local/share/zipcodes/  # Data
~/.config/zipcodes/       # Config
~/.local/state/zipcodes/  # Logs
```

macOS:
```
~/Library/Application Support/zipcodes/  # Data & Config
~/Library/Logs/zipcodes/                 # Logs
```

Windows:
```
%LOCALAPPDATA%\zipcodes\    # Data & Config
%LOCALAPPDATA%\zipcodes\logs\  # Logs
```

**Directory Contents:**
```
data/
├── zipcodes.db           # SQLite database (340K+ records)
└── geoip/               # GeoIP databases (auto-downloaded from jsdelivr CDN)
    ├── geolite2-city-ipv4.mmdb    # ~50MB
    ├── geolite2-city-ipv6.mmdb    # ~40MB
    ├── geo-whois-asn-country.mmdb # ~8MB
    └── asn.mmdb                   # ~5MB

config/
└── admin_credentials     # Admin login info (0600 permissions)
```

## Docker Deployment

### Production (docker-compose.yml)

```bash
# Using docker-compose
docker-compose up -d

# Server available at http://172.17.0.1:64080
```

**docker-compose.yml:**
```yaml
services:
  zipcodes:
    image: ghcr.io/apimgr/zipcodes:latest
    container_name: zipcodes
    restart: unless-stopped

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/zipcodes.db
      # Uncomment and set for first deployment
      #- ADMIN_USER=administrator
      #- ADMIN_PASSWORD=changeme
      #- ADMIN_TOKEN=your-token-here

    volumes:
      - ./rootfs/config/zipcodes:/config
      - ./rootfs/data/zipcodes:/data
      - ./rootfs/logs/zipcodes:/logs

    ports:
      - "172.17.0.1:64080:80"

    networks:
      - zipcodes

    healthcheck:
      test: ["CMD", "/usr/local/bin/zipcodes", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

networks:
  zipcodes:
    name: zipcodes
    external: false
    driver: bridge
```

### Development/Testing (docker-compose.test.yml)

```bash
# Build development image
make docker-dev

# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Server available at http://localhost:64181
# Uses /tmp for ephemeral storage
```

### Docker CLI

```bash
# Production setup
docker run -d \
  --name zipcodes \
  -p 172.17.0.1:64080:80 \
  -v ./rootfs/data/zipcodes:/data \
  -v ./rootfs/config/zipcodes:/config \
  -v ./rootfs/logs/zipcodes:/logs \
  --restart unless-stopped \
  ghcr.io/apimgr/zipcodes:latest

# Development setup
docker run -d \
  --name zipcodes \
  -p 127.0.0.1:64080:80 \
  -v ./rootfs/data/zipcodes:/data \
  -v ./rootfs/config/zipcodes:/config \
  ghcr.io/apimgr/zipcodes:latest
```

### Quick Test (No Installation)

```bash
# Download and run directly
wget https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-linux-amd64
chmod +x zipcodes-linux-amd64
./zipcodes-linux-amd64
```

Server will start on a random port (64000-64999 range) and display the URL.

## API Usage

### Web Interface

Navigate to your server URL:
- Search by zipcode, city, or state
- Real-time autocomplete
- Dark/light theme toggle
- Mobile responsive design

### Quick Examples

Replace `your-server:port` with your actual server address.

```bash
# Download complete dataset (6.3MB, 340K+ records)
curl "http://your-server:8080/api/v1/zipcodes.json" > zipcodes.json

# Search by zipcode
curl "http://your-server:8080/api/v1/zipcode/search?q=94102"

# Search by city
curl "http://your-server:8080/api/v1/zipcode/search?q=San Francisco"

# Search by city and state
curl "http://your-server:8080/api/v1/zipcode/search?q=New York, NY"

# Get specific zipcode details
curl "http://your-server:8080/api/v1/zipcode/94102"

# Autocomplete suggestions
curl "http://your-server:8080/api/v1/zipcode/autocomplete?q=San&limit=10"

# GeoIP lookup
curl "http://your-server:8080/api/v1/geoip?ip=8.8.8.8"
```

### API Endpoints

#### Raw Dataset

```
GET /api/v1/zipcodes.json
```
Returns complete zipcodes.json file (340,000+ records, 6.3MB)

#### Search

```
GET /api/v1/zipcode/search?q={query}
```
Universal search by zipcode, city, state, or prefix

**Examples:**
- `?q=94102` - Find zipcode 94102
- `?q=Boston` - All zipcodes in Boston
- `?q=Miami, FL` - All zipcodes in Miami, FL
- `?q=TX` - Zipcodes in Texas (max 1000)
- `?q=941` - All zipcodes starting with 941

**Response:**
```json
{
  "success": true,
  "count": 1,
  "data": [{
    "state": "CA",
    "city": "San Francisco",
    "county": "San Francisco",
    "zip_code": 94102,
    "latitude": "37.7799",
    "longitude": "-122.4203"
  }]
}
```

#### Get Specific Zipcode

```
GET /api/v1/zipcode/{code}      # JSON
GET /api/v1/zipcode/{code}.txt  # Plain text
```

#### Get by Location

```
GET /api/v1/zipcode/city/{city}
GET /api/v1/zipcode/state/{state}
```

#### Autocomplete

```
GET /api/v1/zipcode/autocomplete?q={query}&limit={count}
```

Returns city, state suggestions (default limit: 10, max: 50)

#### Statistics

```
GET /api/v1/zipcode/stats
```

Returns total zipcodes, states, and cities in database

#### GeoIP Lookups

```
GET /api/v1/geoip?ip={address}      # JSON
GET /api/v1/geoip.txt?ip={address}  # Plain text
POST /api/v1/geoip/batch            # Batch lookup (max 100 IPs)
```

**Example Response:**
```json
{
  "ip": "8.8.8.8",
  "country": "United States",
  "country_code": "US",
  "city": "Mountain View",
  "latitude": 37.4056,
  "longitude": -122.0775,
  "timezone": "America/Los_Angeles",
  "asn": 15169,
  "asn_org": "Google LLC"
}
```

#### Health Check

```
GET /healthz
```

Returns server status, database info, and feature availability

### Response Format

All JSON responses follow this structure:

**Success:**
```json
{
  "success": true,
  "data": { },
  "count": 1
}
```

**Error:**
```json
{
  "success": false,
  "error": "error message"
}
```

### Performance

- **Search Speed**: < 10ms average
- **Database**: SQLite with indexes (zipcode, city, state)
- **Throughput**: 1000+ req/s on modern hardware
- **Memory**: ~50MB baseline + databases (~100MB total)
- **Dataset Size**: 6.3MB JSON, ~15MB SQLite database

## Development

### Requirements

- Go 1.23 or later
- Make (optional, for build automation)
- Docker (optional, for container builds)

### Clone & Build

```bash
# Clone repository
git clone https://github.com/apimgr/zipcodes.git
cd zipcodes

# Build all platforms (8 binaries + host binary)
make build

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o zipcodes ./src

# Run locally
./binaries/zipcodes --port 8080
```

### Build Commands

```bash
make build      # Build all platforms (linux, darwin, windows, freebsd × amd64/arm64)
make test       # Run tests with race detection
make docker     # Build and push multi-arch Docker images
make docker-dev # Build development Docker image (local only)
make release    # Create GitHub release (auto-increments version)
make clean      # Clean build artifacts
```

### Build Output

```
binaries/
├── zipcodes                    # Host platform
├── zipcodes-linux-amd64
├── zipcodes-linux-arm64
├── zipcodes-darwin-amd64
├── zipcodes-darwin-arm64
├── zipcodes-windows-amd64.exe
├── zipcodes-windows-arm64.exe
├── zipcodes-freebsd-amd64
└── zipcodes-freebsd-arm64
```

### Version Management

- Version stored in `release.txt` (semantic versioning)
- Auto-increments on `make release`
- Override with `VERSION=1.0.0 make build`

### Project Structure

```
.
├── src/
│   ├── main.go          # Entry point
│   ├── server/          # HTTP server & routes
│   ├── database/        # SQLite operations & admin schema
│   ├── admin/           # Admin authentication & handlers
│   ├── api/             # API handlers
│   ├── geoip/           # GeoIP integration
│   ├── paths/           # OS-specific directory detection
│   ├── utils/           # Address utilities
│   └── data/            # zipcodes.json source (JSON only)
├── docs/                # Documentation (MkDocs)
│   ├── index.md
│   ├── mkdocs.yml
│   └── requirements.txt
├── .github/workflows/   # GitHub Actions
│   └── build.yml        # Build on push & monthly
├── Makefile             # Build automation (4 targets)
├── Dockerfile           # Multi-stage Alpine container
├── docker-compose.yml   # Production deployment config
├── docker-compose.test.yml  # Development/testing config
├── Jenkinsfile          # CI/CD pipeline (jenkins.casjay.cc)
├── .readthedocs.yml     # ReadTheDocs configuration
├── CLAUDE.md            # Project specification
└── release.txt          # Version tracking
```

### CI/CD

**GitHub Actions** (`.github/workflows/build.yml`):
- Triggers: Push to main, monthly schedule (1st at 3 AM UTC)
- Builds all platform binaries
- Creates multi-arch Docker images (amd64/arm64)
- Publishes to ghcr.io/apimgr/zipcodes
- Creates GitHub releases

**Jenkins Pipeline** (`Jenkinsfile`):
- Multi-architecture builds (amd64/arm64)
- Parallel testing on both architectures
- Docker image publishing
- GitHub releases
- Server: jenkins.casjay.cc

## License & Credits

MIT License - See LICENSE.md for details

### Support

- Issues: https://github.com/apimgr/zipcodes/issues
- For security issues, please email security@apimgr.com

---

**Zipcodes API v1.0** - A production-ready US postal code API with GeoIP integration. Built for simplicity, performance, and ease of deployment.
