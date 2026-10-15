# API reference

All API routes are public and read-only. JSON routes are under `/api/v1`; use `/api/v1/server/swagger` for interactive OpenAPI documentation and `/api/v1/server/openapi.json` for the specification.

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/v1/zipcodes/{code}` | Exact five-digit ZIP lookup |
| GET | `/api/v1/zipcodes/search?q=...` | ZIP, city, state, or city/state search |
| GET | `/api/v1/zipcodes/city/{city}` | City lookup |
| GET | `/api/v1/zipcodes/state/{state}` | State listing |
| GET | `/api/v1/zipcodes/nearby?zip=90210&radius=10` | Radius search from a ZIP |
| GET | `/api/v1/zipcodes/nearby?lat=34.09&lng=-118.41&radius=10` | Radius search from coordinates |
| GET | `/api/v1/me` | Caller IP and available GeoIP fields |
| GET | `/api/v1/zipcodes/stats` | Dataset counts |
| GET | `/api/v1/zipcodes/random` | Random record |
| GET | `/api/v1/zipcodes.json` | Embedded dataset |
| POST | `/api/v1/graphql` | GraphQL query execution |

Radius values must be greater than zero and no more than the configured maximum, which is capped at 100 miles. Coordinate bounds are latitude -90..90 and longitude -180..180. The `/api/v1/me` endpoint returns the IP only when GeoIP is unavailable or the address is not globally routable. GeoIP data is sourced from DB-IP and NRO through the CC BY 4.0 ip-location-db MMDB files.

REST success responses use `ok: true` and `data`; errors use `ok: false`, an error code, a message, and details. List endpoints accept `page` (default `1`) and `limit` (default `250`, maximum `250`). Responses include a `pagination` object with `page`, `limit`, `total`, and `pages`. Exact ZIP lookups are single-record responses. `/api/v1/zipcodes.json` intentionally returns the full wrapped dataset.
