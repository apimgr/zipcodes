# Zipcodes API Documentation

Welcome to the Zipcodes API documentation. This is a fast and accurate US zipcode lookup API with 340,000+ zipcodes, GeoIP integration, and modern web interface.

## Quick Links

- [API Reference](API.md) - Complete API endpoint documentation
- [Server Administration](SERVER.md) - Server setup and configuration guide
- [GitHub Repository](https://github.com/apimgr/zipcodes) - Source code and releases

## Overview

Zipcodes provides a complete REST API and web interface for US postal code lookups with:

- **340,000+ US zipcodes** with city, state, county, and coordinates
- **GeoIP integration** for IP-to-location lookups
- **Fast indexed database** - SQLite with < 10ms queries
- **Modern web interface** with autocomplete
- **Single static binary** - 9.4MB with all assets embedded
- **Admin-only authentication** - Secure admin panel with auto-generated credentials

## Features

- Fast search by zipcode, city, state, or prefix
- Autocomplete suggestions as you type
- GeoIP lookups (IP to location)
- REST API with JSON responses
- Text format alternatives for all endpoints
- Complete dataset download (6.3MB JSON)
- CORS enabled for cross-domain use
- Multi-platform support (Linux, macOS, Windows, FreeBSD)
- Docker deployment ready
- No external dependencies required

## Quick Start

### Docker (Recommended)

```bash
docker-compose up -d
```

Server available at `http://your-server:64080`

### Binary Installation

```bash
# Download latest release
wget https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-linux-amd64

# Make executable and run
chmod +x zipcodes-linux-amd64
./zipcodes-linux-amd64
```

### First Run

On first run, admin credentials are auto-generated and displayed. Save these credentials immediately - they won't be shown again!

## Example API Usage

```bash
# Search by zipcode
curl "http://your-server:8080/api/v1/zipcode/search?q=94102"

# Search by city
curl "http://your-server:8080/api/v1/zipcode/search?q=San Francisco"

# Get specific zipcode
curl "http://your-server:8080/api/v1/zipcode/94102"

# GeoIP lookup
curl "http://your-server:8080/api/v1/geoip?ip=8.8.8.8"

# Download complete dataset
curl "http://your-server:8080/api/v1/zipcodes.json" > zipcodes.json
```

## Data Sources

- **Zipcodes**: US Postal Service data (340,000+ records)
- **GeoIP**: MaxMind GeoLite2 from [P3TERX/GeoLite.mmdb](https://github.com/P3TERX/GeoLite.mmdb)

## Support

- [GitHub Issues](https://github.com/apimgr/zipcodes/issues)
- For security issues, please email security@apimgr.com

## License

MIT License - See [LICENSE.md](https://github.com/apimgr/zipcodes/blob/main/LICENSE.md) for details.
