## Project description

Zipcodes is a full-stack Go web application providing US postal code lookups and geographic queries across 340,000+ ZIP codes through a versioned REST API, GraphQL endpoint, and a server-side rendered web UI. Each ZIP code record includes city, state, county, coordinates, timezone, and telephone area codes. GeoIP integration auto-detects the caller's approximate location. ZIP code data is stored in an indexed embedded SQLite database for sub-10ms query times. A companion CLI tool enables lookups from the terminal. Deployed as a single self-contained static binary.

## Project variables

project_name: zipcodes
project_org: apimgr
internal_name: zipcodes
internal_org: apimgr
app_name: Zipcodes API
repo: https://github.com/apimgr/zipcodes
license: MIT
binary: zipcodes
client_binary: zipcodes-cli

## Business logic

### Product scope & non-goals

**In scope:**
- 340,000+ US ZIP codes with city, state, county, coordinates, IANA timezone identifier, and telephone area codes
- ZIP code lookup by postal code (exact match)
- City name search (optional state filter)
- State-level listing by abbreviation
- Radius-based geographic query: all ZIP codes within N miles of a given ZIP code or coordinate
- Caller GeoIP lookup (country, city, approximate coordinates) via MaxMind GeoLite2
- Full web frontend (server-side Go templates, dark/light/auto theme, PWA, mobile-first)
- Server pages: `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`
- CLI client (`zipcodes-cli`) for terminal lookup: `zipcodes-cli 90210`
- OpenAPI/Swagger docs at `/api/{api_version}/server/swagger`
- GraphQL at `/graphql`

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- No write/mutation API (data is read-only)
- No non-US postal codes (US only; see `citylist` for international city data)
- No real-time data updates (ZIP code dataset is static, updated via releases)
- No paid tiers, no API keys, no rate-limited access tiers

### Roles & permissions

There are no user roles. All endpoints are public and require no authentication.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Full read access to all web pages and API endpoints |
| **Anonymous API client (curl/CLI)** | Full read access to all API endpoints |
| **Server operator** | Configures server via `server.yml` only (max radius, GeoIP update schedule); no web management interface |

### Data model & sensitivity

**ZIP code record** (stored in embedded SQLite, no PII):

| Field | Type | Sensitivity |
|-------|------|-------------|
| `zip` | string — 5-digit US postal code | Public |
| `city` | string — primary city name | Public |
| `state` | string — 2-letter state abbreviation | Public |
| `state_full` | string — full state name | Public |
| `county` | string — county name | Public |
| `lat` | float — latitude | Public |
| `lng` | float — longitude | Public |
| `timezone` | string — IANA timezone identifier | Public |
| `area_codes` | string[] — telephone area codes | Public |

No PII stored or served. GeoIP results are computed per-request from the caller's IP and never stored.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| ZIP code SQLite database (embedded at build) | Fully trusted | Static, compiled into binary |
| MaxMind GeoLite2 (downloaded at first run) | Trusted — HTTPS + checksum verified | Used for caller GeoIP lookup only |
| Incoming HTTP requests | **Untrusted** | ZIP code validated as 5-digit numeric string; coordinate inputs bounds-checked; radius capped |

No external API calls are made on behalf of user requests.

Failure mode for GeoIP: if the GeoLite2 database is unavailable, the `/me` endpoint returns the caller's IP only without location fields. All ZIP code lookup endpoints are unaffected.

### Threat model & abuse cases

**Primary assets:** service availability.

**Attacker/abuser goals:**
- DoS via high-rate ZIP lookup or radius search requests
- Expensive unbounded radius queries to bulk-dump all ZIP code coordinates
- Path traversal or injection via ZIP code parameter

**Defenses:**
- Rate limiting on all endpoints
- Radius search capped at a configurable maximum (default: 100 miles) to bound result set size
- ZIP code parameter validated as a 5-digit numeric string before any database lookup — no dynamic SQL construction from raw input; parameterized queries only
- Request size limits on all query parameters
- No user accounts eliminates credential stuffing and privilege escalation entirely

### Security decisions & exceptions

- **No authentication on any endpoint**: intentional. Public read-only reference API.
- **GeoIP database downloaded at runtime**: intentional for size and freshness. Integrity verified via HTTPS.
- **GeoIP is informational only**: the caller's IP is a standard HTTP observable; GeoIP lookup results are not stored or logged.
- **All responses include `Access-Control-Allow-Origin: *`**: intentional. Public data API designed for cross-origin browser use.
