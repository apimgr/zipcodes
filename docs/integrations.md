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

## Email notifications

Operator email uses the configured SMTP relay. On first run, an empty SMTP host triggers ordered detection of loopback, Docker bridge, default gateway, detected FQDN, public IPv4, and the `mail.` and `smtp.` FQDNs on ports 25, 465, and 587. A successful EHLO probe selects the relay. Startup and configuration reloads probe it again. If no working relay is available, email is disabled and messages are neither sent nor queued; ZIP lookup remains available.

`server.notifications.email.smtp.tls: auto` uses implicit TLS only on port 465. On other ports the client starts with SMTP/EHLO and upgrades with STARTTLS when offered. Authentication credentials are never sent unless TLS is active. Explicit modes are `tls`, `starttls`, and `none`. Certificate verification is relaxed only for configured private IP literals; hostnames remain verified even if DNS resolves them to a private address. `SMTP_*` environment variables override the corresponding YAML settings.

Built-in operator templates are embedded under `src/server/template/email/`. Overrides belong in `{config_dir}/template/email/{template}.tmpl`; each has a `Subject:` line, a `---` separator, and a plain-text body using `{variable}` substitutions. Overrides are read on each render, so edits take effect without restarting. Removing an override restores the embedded version. Unknown variables, empty subjects/bodies, malformed syntax, and header line breaks reject the override or message. The message is delivered separately to each recipient and contains plain-text and HTML MIME alternatives. There are no account, login, password-reset, or verification emails.

Callers translate subject and body strings before rendering, then supply the translated values as ordinary variables. The template engine does not execute translation functions. Integrations should pass event context through the notification event hooks and must keep recipient selection in the owning security/contact subsystem.

Scheduled task failures send the `scheduler_error` template when that event is enabled. Backup and SSL renewal failures use their dedicated templates when enabled, preventing duplicate generic failure messages.

## Monitoring

Prometheus metrics are served at `/server/metrics` and `/server/metrics/prometheus`; `/metrics` is enabled as a scrape-friendly alias by default. Grafana dashboard JSON is available at `/server/metrics/grafana`, and recent sanitized HTTP log events are available as Loki streams at `/server/metrics/loki`. The same service routes are mounted under the configured API version and `/api/metrics` aliases.

Each service has its own Bearer token in `secrets.yml` under `server.metrics.auth.tokens`. The server generates independent `mtr_` credentials for Prometheus, Grafana, and Loki. Send them in `Authorization: Bearer ...`; query-string credentials are rejected. `server.metrics.auth.allow_unauthenticated` can disable auth only for a metrics service protected by a private network or firewall. Metrics are operational data and should not be exposed through a public reverse proxy. See [configuration](configuration.md) for histogram and retention settings.
