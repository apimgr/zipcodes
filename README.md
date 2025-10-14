# Zipcodes - US Postal Code API

Fast and accurate US zipcode lookup API with 340,000+ zipcodes, GeoIP integration, and modern web interface.

## Installation

### Quick Start

```bash
# Download latest release
wget https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-linux-amd64

# Make executable
chmod +x zipcodes-linux-amd64

# Run
./zipcodes-linux-amd64
```

Server starts on http://localhost:8080 (default)

### Docker

```bash
# Using docker-compose
docker-compose up -d

# Or with docker directly
docker run -d \
  -p 8080:80 \
  -v ./data:/data \
  ghcr.io/apimgr/zipcodes:latest
```

### From Source

```bash
# Clone repository
git clone https://github.com/apimgr/zipcodes.git
cd zipcodes

# Build
make build

# Run
./binaries/zipcodes
```

## Usage

### Web Interface

Open http://localhost:8080 in your browser:
- Search by zipcode, city, or state
- Real-time autocomplete
- Dark/light theme toggle
- Mobile responsive

### API Quick Examples

```bash
# Download complete dataset (6.3MB, 340K+ records)
curl "http://localhost:8080/api/v1/zipcodes.json" > zipcodes.json

# Search by zipcode
curl "http://localhost:8080/api/v1/zipcode/search?q=94102"

# Search by city
curl "http://localhost:8080/api/v1/zipcode/search?q=San Francisco"

# Search by city and state
curl "http://localhost:8080/api/v1/zipcode/search?q=New York, NY"

# Get specific zipcode details
curl "http://localhost:8080/api/v1/zipcode/94102"

# Autocomplete suggestions
curl "http://localhost:8080/api/v1/zipcode/autocomplete?q=San&limit=10"

# GeoIP lookup
curl "http://localhost:8080/api/v1/geoip?ip=8.8.8.8"
```

## API Endpoints

### Raw Dataset

```
GET /api/v1/zipcodes.json
```
Returns complete zipcodes.json file (340,000+ records, 6.3MB)

### Search

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

### Get Specific Zipcode

```
GET /api/v1/zipcode/{code}      # JSON
GET /api/v1/zipcode/{code}.txt  # Plain text
```

### Get by Location

```
GET /api/v1/zipcode/city/{city}
GET /api/v1/zipcode/state/{state}
```

### Autocomplete

```
GET /api/v1/zipcode/autocomplete?q={query}&limit={count}
```

Returns city, state suggestions (default limit: 10, max: 50)

### Statistics

```
GET /api/v1/zipcode/stats
```

Returns total zipcodes, states, and cities in database

### GeoIP Lookups

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

### Health Check

```
GET /healthz
```

Returns server status, database info, and feature availability

## Configuration

### Command Line Options

```bash
--help            Show help message
--version         Show version information
--status          Check server status
--port PORT       Set port (default: 8080)
--data DIR        Set data directory
--address ADDR    Listen address (default: 0.0.0.0)
--dev             Development mode
```

### Data Storage

```
~/.local/share/zipcodes/  (default)
├── zipcodes.db           # SQLite database (340K+ records)
└── geoip/               # GeoIP databases (auto-downloaded)
    ├── GeoLite2-City.mmdb
    ├── GeoLite2-Country.mmdb
    └── GeoLite2-ASN.mmdb
```

## Features

- **340,000+ US Zipcodes** with city, state, county, coordinates
- **Fast Search** - Indexed database, < 10ms queries
- **Multiple Search Types** - By zipcode, city, state, prefix
- **Autocomplete** - Smart suggestions as you type
- **GeoIP Integration** - IP to location lookups
- **REST API** - Full JSON API with text alternatives
- **Web Interface** - Modern, responsive, dark/light themes
- **Download Dataset** - Complete JSON file available
- **CORS Enabled** - Use from any domain
- **Single Binary** - No dependencies, embedded assets

## Data Sources

- **Zipcodes**: US Postal Service data (340,000+ records)
- **GeoIP**: MaxMind GeoLite2 from [P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb)
  - Auto-downloads on first run
  - Updates available via admin (future)

## Development

### Requirements

- Go 1.21 or later
- Make

### Build

```bash
# Development build
make dev

# Production build (all platforms)
make build

# Run tests
make test

# Docker image
make docker
```

### Project Structure

```
.
├── src/
│   ├── main.go          # Entry point
│   ├── server/          # HTTP server & routes
│   ├── database/        # SQLite operations
│   ├── api/             # API handlers
│   ├── geoip/           # GeoIP integration
│   └── data/            # zipcodes.json source
├── Makefile             # Build automation
├── Dockerfile           # Container image
└── docker-compose.yml   # Deployment config
```

## Response Format

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

## Performance

- **Search Speed**: < 10ms average
- **Database**: SQLite with indexes (zipcode, city, state)
- **Throughput**: 1000+ req/s on modern hardware
- **Memory**: ~50MB baseline + databases (~100MB total)
- **Dataset Size**: 6.3MB JSON, ~15MB SQLite database

## License

MIT License - See LICENSE.md for details

## Support

- Issues: https://github.com/apimgr/zipcodes/issues

## Security

For security issues, please email security@apimgr.com instead of using the issue tracker.
