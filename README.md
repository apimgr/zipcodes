# Zipcodes - US Postal Code API

Fast and accurate US zipcode lookup API with 33,100 geographic ZIP records, GeoIP integration, and modern web interface.

**Static Go binary with embedded assets and ZIP data. GeoIP data is downloaded separately at runtime when enabled.**

## About

Zipcodes provides a complete REST API and web interface for US postal code lookups with:
- 33,100 U.S. ZIP records with city, state, county, and coordinates
- GeoIP integration for IP-to-location lookups
- Fast indexed SQLite database (< 10ms queries)
- Modern web interface with autocomplete
- Single static binary with all web assets and the ZIP dataset embedded

### Features

- **33,100 U.S. ZIP records** with city, state, county, coordinates
- **Fast Search** - Indexed database, < 10ms queries
- **Multiple Search Types** - By zipcode, city, state, prefix
- **Autocomplete** - Smart suggestions as you type
- **GeoIP Integration** - IP to location lookups
- **REST API** - Full JSON API with text alternatives
- **Open read-only access** - No accounts, API keys, or write endpoints
- **Web Interface** - Modern, responsive, dark/light themes
- **Download Dataset** - Complete JSON file available
- **CORS Enabled** - Use from any domain
- **Single Binary** - static binary with embedded web assets and ZIP records
- **Multi-Platform** - Linux, macOS, Windows, FreeBSD (amd64/arm64)

### Data Sources

- **Zipcodes**: Ready APIs curated-us-zips (33,100 records; CC BY 4.0; see data/licenses/zipcodes-data.md)
- **GeoIP**: City data from DB-IP and country/ASN data from NRO are downloaded at runtime from ip-location-db under CC BY 4.0 and refreshed on the configured schedule. GeoIP may return only the caller IP when a database is unavailable.
- **ZIP records**: Ready APIs curated-us-zips, pinned at `64ac28b620a39d9efabe39f91dc03181a13a8a76`, under CC BY 4.0. The upstream dataset includes 33,100 geographic ZIP records and excludes APO/FPO/DPO ZIPs.
- **Known data gap**: The embedded dataset has no verified ZIP-to-area-code mapping; area-code arrays are empty. See `data/licenses/zipcodes-data.md`.

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
- `zipcodes-cli-linux-amd64` - CLI for Linux (Intel/AMD)
- `zipcodes-cli-linux-arm64` - CLI for Linux (ARM)
- `zipcodes-cli-darwin-amd64` - CLI for macOS (Intel)
- `zipcodes-cli-darwin-arm64` - CLI for macOS (Apple Silicon)
- `zipcodes-cli-windows-amd64.exe` - CLI for Windows (Intel/AMD)
- `zipcodes-cli-windows-arm64.exe` - CLI for Windows (ARM)
- `zipcodes-cli-freebsd-amd64` - CLI for FreeBSD (Intel/AMD)
- `zipcodes-cli-freebsd-arm64` - CLI for FreeBSD (ARM)

### First Run

The server creates its configuration and data directories on first startup. All lookup endpoints are public and read-only; no account setup is required. The CLI needs `--server URL` or `ZIPCODES_SERVER_PRIMARY`; this project does not define a default hosted API.

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
├── db/
│   └── server.db       # ZIP records, rate limits, and scheduler state
└── security/geoip/     # GeoIP databases (downloaded from ip-location-db sources via HTTPS)
    ├── dbip-city-ipv4.mmdb
    ├── dbip-city-ipv6.mmdb
    ├── geo-whois-asn-country.mmdb # ~8MB
    └── asn.mmdb                   # ~5MB

config/
```

## CLI Usage

The CLI performs exact ZIP lookups against a configured public server:

```bash
zipcodes-cli --server https://your-zipcodes-server.example 90210
```

Set `ZIPCODES_SERVER_PRIMARY` or configure `server.primary` in `cli.yml` to avoid repeating `--server`.

## Docker Deployment

### Production (`docker/docker-compose.yml`)

```bash
docker compose -f docker/docker-compose.yml up -d
```

The production stack binds the HTTP service to `172.17.0.1:64580`, stores configuration and data in `./volumes/`, and uses the binary's production defaults.

### Development (`docker/docker-compose.dev.yml`)

```bash
docker compose -f docker/docker-compose.dev.yml up -d
```

The development image is compiled into the `:devel` image; it does not use source bind mounts or hot reload.

### Automated testing (`docker/docker-compose.test.yml`)

Copy this file and its `./volumes/` paths into an isolated temporary deployment directory before starting it. Never run the test compose file from the repository checkout.

### Docker CLI

```bash
# Production setup
docker run -d \
  --name zipcodes \
  -p 172.17.0.1:64080:80 \
  -v ./volumes/data/zipcodes:/data \
  -v ./volumes/config/zipcodes:/config \
  -v ./volumes/logs/zipcodes:/logs \
  --restart unless-stopped \
  ghcr.io/apimgr/zipcodes:latest

# Development setup
docker run -d \
  --name zipcodes \
  -p 127.0.0.1:64080:80 \
  -v ./volumes/data/zipcodes:/data \
  -v ./volumes/config/zipcodes:/config \
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

## CLI Client

The companion `zipcodes-cli` looks up a ZIP code through a configured Zipcodes server. Set the server with `--server`, the `ZIPCODES_SERVER_PRIMARY` environment variable, or `~/.config/apimgr/zipcodes/cli.yml`:

```bash
zipcodes-cli --server https://YOUR_SERVER 90210
```

The first valid `--server` value is saved in `cli.yml` with user-only file permissions.

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
# Download complete dataset (dataset size depends on formatting, 33,100 records)
curl -q -LSsf "http://your-server:8080/api/v1/zipcodes.json" > zipcodes.json

# Search by zipcode
curl -q -LSsf "http://your-server:8080/api/v1/zipcodes/search?q=94102"

# Search by city
curl -q -LSsf "http://your-server:8080/api/v1/zipcodes/search?q=San%20Francisco"

# Search by city and state
curl -q -LSsf "http://your-server:8080/api/v1/zipcodes/search?q=New%20York%2C%20NY"

# Get specific zipcode details
curl -q -LSsf "http://your-server:8080/api/v1/zipcodes/94102"

# Autocomplete suggestions
curl -q -LSsf "http://your-server:8080/api/v1/zipcodes/autocomplete?q=San&limit=10"

# GeoIP lookup
curl -q -LSsf "http://your-server:8080/api/v1/me"
```

### API Endpoints

#### Raw Dataset

```
GET /api/v1/zipcodes.json
```
Returns an `{ "ok": true, "data": [...] }` response containing the embedded ZIP records (roughly 3 MB before response formatting).

#### Search

```
GET /api/v1/zipcodes/search?q={query}
```
Universal search by zipcode, city, state, or prefix

**Examples:**
- `?q=94102` - Find zipcode 94102
- `?q=Boston` - All zipcodes in Boston
- `?q=Miami, FL` - All zipcodes in Miami, FL
- `?q=TX` - Zipcodes in Texas (full state listing)
- `?q=941` - All zipcodes starting with 941

**Response:**
```json
{
  "ok": true,
  "count": 1,
  "data": [{
    "state": "CA",
    "city": "San Francisco",
    "county": "San Francisco",
    "zip_code": 94102,
    "latitude": 37.7799,
    "longitude": -122.4203,
    "state_full": "California",
    "timezone": "America/Los_Angeles",
    "area_codes": []
  }]
}
```

#### Get Specific Zipcode

```
GET /api/v1/zipcodes/{code}      # JSON
GET /api/v1/zipcodes/{code}.txt  # Plain text
```

#### Get by Location

```
GET /api/v1/zipcodes/city/{city}
GET /api/v1/zipcodes/state/{state}
```

#### Autocomplete

```
GET /api/v1/zipcodes/autocomplete?q={query}&limit={count}
```

Returns city, state suggestions (default limit: 10, max: 50)

#### Statistics

```
GET /api/v1/zipcodes/stats
```

Returns total zipcodes, states, and cities in database

#### GeoIP Lookups

```
GET /api/v1/me  # Caller IP and available GeoIP details
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
  "ok": true,
  "data": { },
  "count": 1
}
```

**Error:**
```json
{
  "ok": false,
  "error": "REQUEST_FAILED",
  "message": "The request could not be completed.",
  "details": {}
}
```

### Performance

- **Search Speed**: < 10ms average
- **Database**: SQLite with indexes (zipcode, city, state)
- **Radius bounds**: configurable up to 100 miles
- **Rate limits**: configurable per-IP read, write, health, and global limits

## Development

### Requirements

- Go 1.25 or later
- Make (optional, for build automation)
- Docker (optional, for container builds)

### Clone & Build

```bash
# Clone repository
git clone https://github.com/apimgr/zipcodes.git
cd zipcodes

# Build server and CLI for all platforms
make build

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o zipcodes ./src

# Run locally
./binaries/zipcodes --port 8080
```

### Build Commands

```bash
make build      # Build server and CLI for linux, darwin, windows, freebsd × amd64/arm64
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
│   ├── database/        # SQLite data access
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
├── Makefile             # Docker-based Go build and test targets
├── docker/              # Production, development, and test container definitions
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

## Disclaimer

This software is provided "as is" without warranty of any kind. Use it at your own risk.

- The authors are not responsible for damages, data loss, or issues arising from use.
- This software does not provide legal, financial, medical, or other professional advice.
- GeoIP data and other third-party services remain subject to their providers' terms.
- Evaluate the software and its security posture before production deployment.

## License & Credits

MIT License - See LICENSE.md for details

### Support

- Issues: https://github.com/apimgr/zipcodes/issues
- For security issues, please email security@apimgr.com

---

**Zipcodes API v1.0** - A production-ready US postal code API with GeoIP integration. Built for simplicity, performance, and ease of deployment.
