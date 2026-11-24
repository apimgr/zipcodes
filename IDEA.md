## Project description

Zipcodes is a Go web application providing US ZIP-code lookups and geographic queries through a versioned REST API, GraphQL endpoint, and server-rendered search UI. The checked-in source currently contains 33,100 geographic ZIP records with city, state, county, coordinates, and timezone. Its source does not provide telephone area-code assignments, so `area_codes` is empty. Optional GeoIP lookup resolves the caller's approximate location. ZIP records are loaded from the embedded dataset into SQLite by default; an operator may configure the supported remote libSQL/Turso backend. A companion CLI enables lookups from the terminal.

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
owner_token: zipcodes_owner_token_qwhzkfta

## Business logic

### Product scope & non-goals

**In scope:**
- 33,100 US geographic ZIP records from the checked-in Ready APIs curated-us-zips source, with city, state, county, coordinates, and IANA timezone identifier
- No telephone area-code dataset is currently available from the checked-in source; `area_codes` remains empty
- ZIP code lookup by postal code (exact match)
- City name search (optional state filter)
- State-level listing by abbreviation
- Radius-based geographic query: all ZIP codes within N miles of a given ZIP code or coordinate
- Caller GeoIP lookup (country, city, approximate coordinates) via ip-location-db datasets (DB-IP city under CC BY 4.0; iptoasn country and ASN under PDDL)
- Full web frontend (server-side Go templates, dark/light/auto theme, PWA, mobile-first)
- Server-rendered pages include `/`, `/search`, `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/security`, `/server/security/policy`, `/server/security/thanks`, `/server/security/report/{tracking_id}`, `/server/dpo`, `/server/contact`, `/server/help/mcp`, `/server/preferences`, `/server/preferences/tokens`, and `/server/preferences/sync`; `/server/docs/swagger` and `/server/docs/graphql` provide the interactive API interfaces. The shared page registry and each page's behavior follow AI.md.
- CLI client (`zipcodes-cli`) for terminal lookup: `zipcodes-cli 90210`
- OpenAPI JSON at `/api/{api_version}/server/swagger` with a server-rendered reference at `/server/docs/swagger`
- GraphQL API at `/api/{api_version}/server/graphql` and a no-JavaScript form at `/server/docs/graphql`

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- ZIP lookup and GeoIP data are read-only; AI.md's shared token and preference management endpoints are template-owned control flows and do not mutate ZIP data
- No non-US postal codes (US only; see `citylist` for international city data)
- No real-time data updates (ZIP code dataset is static, updated via releases)
- No paid tiers, no paid API keys, no differentiated access tiers

### Roles & permissions

There are no user roles. Lookup and server information endpoints are public and read-only; the AI.md-defined operator and owner tokens apply only to the server management features that use them.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Public search and server information pages |
| **Anonymous API client (curl/CLI)** | Public read-only lookup routes subject to rate limits |
| **Server operator** | Configures server through `server.yml` and operator CLI; no admin web panel |

AI.md also defines anonymous owner/linked API tokens and shared server pages. Those template-owned capabilities are separate from product accounts and must not be described as user registration or login.

### Data model & sensitivity

**ZIP code record** (stored in local SQLite by default, or remote libSQL/Turso when configured; no PII):

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

ZIP dataset records contain no PII. GeoIP results are computed per-request from the caller's IP and are not retained. The required security-report workflow stores report content and researcher contact details encrypted at rest; public status pages expose only report state, maintainer comments, and disclosure dates.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| Embedded ZIP dataset | Trusted | Static source data compiled into the binary and used to seed an empty database |
| Local SQLite database | Trusted local state | Default storage; WAL and a single writer with separate read connections |
| Remote libSQL/Turso database | Operator-selected external service | Authenticated remote storage; table and index names are prefixed; the provider owns backups and only an empty remote database may be seeded from the embedded dataset |
| ip-location-db MMDB datasets (downloaded at first run) | Trusted — HTTPS transport plus release SHA-256 verification; city data from DB-IP (CC BY 4.0), country/ASN from iptoasn.com (PDDL) | Caller lookup and optional country-deny signal; never stored |
| Incoming HTTP requests | **Untrusted** | ZIP code validated as 5-digit numeric string; coordinate inputs bounds-checked; radius capped |

No external API calls are made on behalf of ZIP lookup requests. Optional GeoIP data refreshes, configured notification delivery, and researcher-supplied OpenPGP key retrieval are separate server workflows governed by AI.md.

Failure mode for GeoIP: if the managed ip-location-db files are unavailable, lookups return the IP without location fields and country filtering fails open. Private/internal addresses are never looked up or country-blocked. Country restrictions are disabled by default and, when configured, they run after rate limiting and before authentication as a deny-only signal; they do not replace authentication or authorization. All ZIP code lookup endpoints remain read-only.

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

- **No end-user accounts or login**: intentional. The operator token remains for server management features defined by AI.md.
- **GeoIP databases downloaded at runtime**: intentional for size and freshness. DB-IP city data and iptoasn.com country/ASN source licenses follow AI.md PART 19.
- **GeoIP lookup data is not retained**: configured country restrictions are a deny-only risk signal. Authentication and authorization still apply to every request that passes the country check.
- **CORS defaults to `*`**: public read-only data is cross-origin accessible by default; operators may configure an explicit allow-list or disable CORS under `server.cors`.

### API routes

Project-defined lookup routes use `/api/{api_version}` (default `v1`, reported by `/api/autodiscover`); they are public and read-only, with rate limits defined by AI.md. Template-owned server routes are specified by AI.md and are not repeated in this inventory.

- GET `/api/{api_version}/zipcodes/{code}`: exact ZIP lookup; a five-digit code returns one ZIP record or a not-found response.
- GET `/api/{api_version}/zipcodes/search?q={query}`: search ZIP, city, state, or city/state; returns paginated matching records.
- GET `/api/{api_version}/zipcodes/autocomplete?q={query}`: return ZIP/city/state suggestions.
- GET `/api/{api_version}/zipcodes/nearby?zip={zip}&radius={miles}` or `?lat={lat}&lng={lng}&radius={miles}`: return nearby ZIP records and distances within the configured radius cap.
- GET `/api/{api_version}/zipcodes/city/{city}`: return paginated ZIP records for a city.
- GET `/api/{api_version}/zipcodes/state/{state}`: return paginated ZIP records for a state abbreviation.
- GET `/api/{api_version}/zipcodes/stats`: return dataset counts.
- GET `/api/{api_version}/zipcodes/random`: return one random ZIP record.
- GET `/api/{api_version}/zipcodes.json`: return the embedded ZIP dataset.
- GET `/api/{api_version}/zipcodes/{code}.txt`, `/api/{api_version}/zipcodes/stats.txt`, `/api/{api_version}/zipcodes/random.txt`: text representations of the corresponding lookup, statistics, and random record.
- GET `/api/{api_version}/me`: return caller IP and available GeoIP fields; location fields may be absent when GeoIP data is unavailable.
- GET `/api/{api_version}/geoip?ip={address}` and GET `/api/{api_version}/geoip.txt?ip={address}`: return GeoIP lookup data in JSON or text.
- POST `/api/{api_version}/geoip/batch`: resolve a bounded batch of IP addresses and return their GeoIP results.
Browser pages provide ZIP search, public service information and status, privacy and terms disclosures, security reporting, preferences, and interactive API documentation. The complete current page inventory is listed under Product scope; AI.md governs shared shell, consent, error-page, registry, and accessibility behavior.
