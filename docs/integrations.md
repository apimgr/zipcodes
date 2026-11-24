# Integrations

Zipcodes exposes public read-only ZIP and GeoIP lookups, plus authenticated owner-token management and MCP connections:

- REST API under `/api/{api_version}` (default `/api/v1`, discoverable from `/api/autodiscover`) with an OpenAPI document at `/api/{api_version}/server/swagger` and browser UI at `/server/docs/swagger`.
- Remote MCP Streamable HTTP at `/api/v1/server/mcp`, with `/api/mcp` and `/mcp` aliases. It uses the existing owner/operator Bearer tokens and exposes OpenAPI operations other than token management; see [API reference](api.md#mcp).
- Read-only GraphQL at `/api/{api_version}/server/graphql`, with a browser interface at `/server/docs/graphql`.
- The PWA manifest at `/manifest.json` and service worker at `/sw.js`.
- `/robots.txt` and the generated `/sitemap.xml` for public page discovery.
- Empty Atom (`/feed.xml`) and RSS (`/rss.xml`) feeds because ZIP records have no public per-record page or publication timestamp.
- `/.well-known/llms.txt` with request-derived absolute links to the API and its documentation.
- A companion `zipcodes-cli` that connects to an operator-selected server.
- `/security.txt` and `/.well-known/security.txt` for vulnerability reporting information.

API clients can optionally use anonymous owner credentials to associate activity with an owner. A primary owner token manages linked tokens; linked tokens are for automation and cannot manage tokens. This is credential-based identity only: there are no user accounts, login sessions, or federated identity flows. See [API reference](api.md) and [Security](security.md) for token routes and handling.

The project does not provide webhooks, identity federation, user accounts, or a hosted API address. The CLI server URL is supplied with `--server`, `ZIPCODES_SERVER_PRIMARY`, or `server.primary` in its `cli.yml`.

GeoIP database downloads are an optional feature used for caller/IP lookup. The server fetches the fixed ip-location-db assets over HTTPS into `{data_dir}/security/geoip` and verifies release SHA-256 checksums before replacement; an unavailable dataset leaves ZIP lookup functionality available and country restrictions fail open.

## Monitoring

Prometheus metrics are served at `/server/metrics` and `/server/metrics/prometheus`; `/metrics` is enabled as a scrape-friendly alias by default. Grafana dashboard JSON is available at `/server/metrics/grafana`, and recent sanitized HTTP log events are available as Loki streams at `/server/metrics/loki`. The same service routes are mounted under the configured API version and `/api/metrics` aliases.

Each service has its own Bearer token in `secrets.yml` under `server.metrics.auth.tokens`. The server generates independent `mtr_` credentials for Prometheus, Grafana, and Loki. Send them in `Authorization: Bearer ...`; query-string credentials are rejected. `server.metrics.auth.allow_unauthenticated` can disable auth only for a metrics service protected by a private network or firewall. Metrics are operational data and should not be exposed through a public reverse proxy. See [configuration](configuration.md) for histogram and retention settings.
