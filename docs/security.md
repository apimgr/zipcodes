# Security

All zipcode routes are public, read-only, and unauthenticated by product design. ZIP input and radius/coordinate values are validated, SQL uses bound parameters, CORS is open for the public dataset, and forwarded client-IP headers are honored only from trusted peers.

Per-IP request limits default to 120 read requests/minute, 10 non-read requests/minute, 120 health requests/minute, and a 240 request/minute per-IP ceiling across endpoint classes. Counters are stored in `server.db` and keyed by a SHA-256 digest of the client IP and endpoint class; expired events are removed as requests arrive.

`/server/healthz`, `/server/privacy`, `/server/terms`, and `/server/about` are public pages. `/api/v1/me` uses the caller IP from the TCP peer or a configured trusted proxy. GeoIP is skipped for private/reserved addresses and fails open to an IP-only response when location data is unavailable. City records come from DB-IP, and country/ASN data comes from NRO, using the CC BY 4.0 ip-location-db MMDB files. The update scheduler refreshes the enabled databases at the configured interval.

Report vulnerabilities privately through [GitHub Security Advisories](https://github.com/apimgr/zipcodes/security/advisories/new). The public `/.well-known/security.txt` endpoint lists the current contact and expiration date.
