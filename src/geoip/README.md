# GeoIP Module

This module provides GeoIP lookup functionality using MaxMind GeoLite2 databases from the P3TERX/GeoLite.mmdb GitHub repository.

## Features

- **Automatic Database Download**: Downloads latest databases from GitHub releases
- **Multiple Database Support**: City, Country, and ASN databases
- **Automatic Updates**: Scheduled updates with configurable intervals
- **HTTP Handlers**: Ready-to-use HTTP handlers for API endpoints
- **Batch Lookups**: Support for bulk IP lookups
- **Thread-Safe**: Concurrent lookup support with read-write locks

## Usage

### Initialize GeoIP

```go
import "github.com/casapps/zipcodes/src/geoip"

// Download databases (if not present)
dbFiles, err := geoip.DownloadDatabases("/path/to/data")
if err != nil {
    log.Fatal(err)
}

// Initialize GeoIP
err = geoip.Initialize(dbFiles.CityDB, dbFiles.CountryDB, dbFiles.ASNDB)
if err != nil {
    log.Fatal(err)
}
```

### Lookup IP Address

```go
// Lookup single IP
location, err := geoip.LookupIP("8.8.8.8")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Country: %s\n", location.Country)
fmt.Printf("City: %s\n", location.City)
fmt.Printf("Coordinates: %f, %f\n", location.Latitude, location.Longitude)
```

### HTTP Handlers

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/casapps/zipcodes/src/geoip"
)

r := chi.NewRouter()

// JSON response
r.Get("/api/v1/geoip", geoip.LookupHandler)

// Plain text response
r.Get("/api/v1/geoip.txt", geoip.LookupTextHandler)

// Batch lookup
r.Post("/api/v1/geoip/batch", geoip.BatchLookupHandler)
```

### Automatic Updates

```go
// Create updater with configuration
updater := geoip.NewUpdater(&geoip.UpdaterConfig{
    DataDir:       "/path/to/data",
    CheckInterval: 24 * time.Hour,  // Check daily
    AutoUpdate:    true,             // Auto-download updates
    OnUpdateFunc: func() {
        log.Println("GeoIP databases updated!")
    },
    OnErrorFunc: func(err error) {
        log.Printf("Update error: %v", err)
    },
})

// Start automatic updates
updater.Start()
defer updater.Stop()
```

### Manual Update

```go
// Trigger manual update
err := updater.ManualUpdate()
if err != nil {
    log.Printf("Manual update failed: %v", err)
}
```

### Scheduled Task (Cron)

```go
// Get a function for cron scheduler
updateTask := geoip.GetScheduledTask("/path/to/data")

// Use with your scheduler
scheduler.AddTask("geoip-update", "0 3 * * *", updateTask)
```

## API Endpoints

### GET /api/v1/geoip

Lookup IP address (JSON response)

**Query Parameters:**
- `ip` (optional) - IP address to lookup. If not provided, uses client IP.

**Example:**
```bash
curl "http://localhost:8080/api/v1/geoip?ip=8.8.8.8"
```

**Response:**
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

### GET /api/v1/geoip.txt

Lookup IP address (plain text response)

**Example:**
```bash
curl "http://localhost:8080/api/v1/geoip.txt?ip=8.8.8.8"
```

**Response:**
```
IP: 8.8.8.8
Country: United States (US)
City: Mountain View
Coordinates: 37.4056, -122.0775
Timezone: America/Los_Angeles
ASN: 15169 (Google LLC)
```

### POST /api/v1/geoip/batch

Batch lookup multiple IPs

**Request Body:**
```json
{
  "ips": ["8.8.8.8", "1.1.1.1", "208.67.222.222"]
}
```

**Response:**
```json
{
  "success": true,
  "count": 3,
  "results": [
    {
      "ip": "8.8.8.8",
      "country": "United States",
      ...
    },
    ...
  ]
}
```

**Limits:**
- Maximum 100 IPs per request

## Database Sources

Databases are automatically downloaded from:
- **Repository**: https://github.com/P3TERX/GeoLite.mmdb
- **License**: Creative Commons Attribution-ShareAlike 4.0 International License
- **Update Frequency**: Regularly updated by maintainer

## Files

- `geoip.go` - Core GeoIP lookup functionality
- `downloader.go` - GitHub release database downloader
- `handlers.go` - HTTP request handlers
- `updater.go` - Automatic update scheduler

## Database Storage

Databases are stored in `{dataDir}/geoip/`:
- `GeoLite2-City.mmdb` - City-level location data
- `GeoLite2-Country.mmdb` - Country-level location data
- `GeoLite2-ASN.mmdb` - ASN and organization data

## Thread Safety

All lookup operations are thread-safe and support concurrent access. The module uses read-write locks to allow multiple simultaneous lookups while protecting database reload operations.

## Error Handling

- Invalid IP addresses return error
- Missing databases gracefully degrade (Country fallback if City missing)
- Download failures are logged but don't crash the application
- Update failures retain existing databases
