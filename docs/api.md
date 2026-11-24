# API reference

All lookup routes are public and read-only. The REST API base path is `/api/v1`. Interactive OpenAPI documentation is served at `/server/docs/swagger`; the OpenAPI JSON is `/api/v1/server/swagger` (also `/api/swagger`). The GraphQL interface is `/server/docs/graphql`; submit GraphQL queries to `/api/v1/server/graphql` or `/api/graphql` with GET or POST. The `/graphql` browser path is also available.

## ZIP code routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/v1/zipcodes/{code}` | Look up one five-digit ZIP code |
| GET | `/api/v1/zipcodes/{code}.txt` | Plain text ZIP lookup |
| GET | `/api/v1/zipcodes/search?q={query}` | Search ZIP, city, state, or city/state |
| GET | `/api/v1/zipcodes/autocomplete?q={query}` | Return suggestions |
| GET | `/api/v1/zipcodes/city/{city}` | List ZIP codes for a city |
| GET | `/api/v1/zipcodes/state/{state}` | List ZIP codes for a state abbreviation |
| GET | `/api/v1/zipcodes/nearby?zip={zip}&radius={miles}` | Search around an origin ZIP |
| GET | `/api/v1/zipcodes/nearby?lat={lat}&lng={lng}&radius={miles}` | Search around coordinates |
| GET | `/api/v1/zipcodes/stats` | Return dataset counts |
| GET | `/api/v1/zipcodes/stats.txt` | Return dataset counts as text |
| GET | `/api/v1/zipcodes/random` | Return a random ZIP record |
| GET | `/api/v1/zipcodes/random.txt` | Return a random record as text |
| GET | `/api/v1/zipcodes.json` | Download the full embedded dataset |

Radius must be greater than zero and may not exceed the server's configured maximum (100 miles by default). Coordinate inputs are bounded to latitude -90..90 and longitude -180..180. ZIP codes must contain exactly five digits.

## GeoIP routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/v1/me` | Return the caller IP and available location data |
| GET | `/api/v1/geoip?ip={address}` | Look up an IP address |
| GET | `/api/v1/geoip.txt?ip={address}` | Text form of an IP lookup |
| POST | `/api/v1/geoip/batch` | Look up a bounded batch of IP addresses |

GeoIP is informational. Private or reserved addresses and unavailable databases may produce IP-only results. ZIP lookups do not depend on GeoIP availability.

## Server information routes

The matching browser pages are `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, and `/server/terms`. Their versioned API routes are `/api/v1/server/{about,help,healthz,privacy,terms}`.

The server also publishes `/robots.txt`, `/sitemap.xml`, `/favicon.ico`, `/security.txt`, and `/.well-known/security.txt`. The sitemap lists public HTML pages only; it does not list API endpoints.

`/feed.xml` (Atom 1.0) and `/rss.xml` (RSS 2.0) are always available. ZIP records do not have individual public pages or creation times, so both feeds contain no entries. They use request-derived absolute URLs and support ETag conditional requests. `/.well-known/llms.txt` and `/llms.txt` list the public API and documentation endpoints. `/.well-known/pgp-key.asc` returns 404 because this deployment has no security-report keypair.

HTML pages support server-rendered search and require no JavaScript. Theme changes use a CSRF-protected same-origin form and are stored in a first-party cookie.

`/server/healthz` is the canonical health frontend: browsers receive HTML, HTTP tools receive human-readable text, and `Accept: application/json` receives JSON. `/api/v1/server/healthz` and `/api/healthz` are machine endpoints; they return JSON by default and support text with `.txt` or `Accept: text/plain`. The root `/healthz` alias is disabled by default and can be enabled with `server.healthz.root.enabled: true`; it uses the same content negotiation as `/server/healthz`. Health responses are bare and are not wrapped in the standard API envelope. The OpenAPI and GraphQL endpoints are described above. Other JSON API responses use `ok` and `data` on success; errors use `ok`, `error`, `message`, and `details`.
