# Integrations

Zipcodes exposes public read-only ZIP and GeoIP lookups, plus authenticated owner-token management and MCP connections:

- REST API under `/api/{api_version}` (default `/api/v1`, discoverable from `/api/autodiscover`) with an OpenAPI document at `/api/{api_version}/server/swagger` and browser UI at `/server/docs/swagger`.
- Remote MCP Streamable HTTP at `/api/{api_version}/server/mcp`, with `/api/mcp` and `/mcp` aliases. It uses the existing owner/operator Bearer tokens and exposes eligible OpenAPI operations, excluding token management and the cookie-bound `/server/preferences/**` family; see [API reference](api.md#mcp).
- Read-only GraphQL at `/api/{api_version}/server/graphql`, with a browser interface at `/server/docs/graphql`.
- The PWA manifest at `/manifest.json` and service worker at `/sw.js`.
- `/robots.txt` and the generated `/sitemap.xml` for public page discovery.
- Empty Atom (`/feed.xml`) and RSS (`/rss.xml`) feeds because ZIP records have no public per-record page or publication timestamp.
- `/.well-known/llms.txt` with request-derived absolute links to the API and its documentation.
- A companion `zipcodes-cli` that connects to an operator-selected server.
- `/.well-known/security.txt` for vulnerability reporting information.

API clients can optionally use anonymous owner credentials to associate activity with an owner. A primary owner token manages linked tokens; linked tokens are for automation and cannot manage tokens. This is credential-based identity only: there are no user accounts, login sessions, or federated identity flows. See [API reference](api.md) and [Security](security.md) for token routes and handling.

The project does not provide identity federation, user accounts, native app association links, or a hosted API address. The CLI server URL is supplied with `--server`, `ZIPCODES_SERVER_PRIMARY`, or `server.primary` in its `cli.yml`.

## I2P eepsite

`server.i2p.enabled` controls the optional I2P eepsite. It defaults on in containers and off elsewhere; an explicit `--i2p` or `--i2p=false` takes precedence over `I2P_ENABLED`, which takes precedence over `server.yml`. When enabled, the server prefers a dedicated i2pd process and otherwise uses a reachable SAMv3 bridge. The destination key is stored under `{data_dir}/i2p/site/site-keys.dat` with private permissions, and the `.b32.i2p` address remains stable across restarts. I2P forwarding uses a separate loopback backend listener and does not add a PROXY header. Runtime state is exposed by the server health response; see [configuration](configuration.md#configuration) for the settings and limits.

## Tor hidden service

When a Tor executable is available, the server starts its own Tor process with app-scoped configuration and data directories. The generated v3 hidden service forwards to a dedicated loopback backend that requires Tor's PROXY protocol header; it does not share the public listener. Outbound SOCKS access is disabled by default; when enabled it binds to loopback and explicitly accepts only loopback clients. Tor is configured as a client and hidden service, never as a relay or exit node. Bootstrap runs separately from the clearnet listeners, and Tor startup failures leave the public server available. The server checks the Tor control endpoint every 30 seconds and restarts the managed process if the check fails. Set `server.tor.expose: true` to include a healthy `.onion` address in `/api/autodiscover`; Tor and I2P addresses are omitted by default. See [configuration](configuration.md#configuration) for the Tor options and limits.

Tor requests use onion URLs and UTC timestamps, suppress clearnet email addresses and project-host links, and identify clients with Tor's opaque circuit token (`tor:{id}`) rather than logging the local daemon address. The fallback identity is `tor`. If `server.tor.contact_email` is set, it is the only email shown on onion responses. Clearnet HTML document responses advertise the matching onion URL with `Onion-Location`; JSON, assets, redirects, and onion responses do not receive that header.

## Contact webhooks

Contact roles can deliver notifications through `server.contact.{admin,security,abuse,general}.webhooks`. Supported names are `telegram`, `discord`, `slack`, `mattermost`, `pushover`, `gotify`, and `generic`. Each configured HTTPS destination receives the adapter-specific request format described in [configuration](configuration.md), plus an HMAC-SHA256 signature, Unix timestamp, UUIDv7 delivery ID, event name, and application User-Agent. Signing secrets are generated per role and transport in `secrets.yml`; they are not exposed in public pages, APIs, or logs.

Security-role webhook bodies contain the already encrypted report. Abuse, DMCA, and privacy submissions use their resolved role destinations and include submitted report content, so configure those endpoints as private destinations. Admin notifications contain summaries and must not include user-submitted content.

Webhook delivery runs asynchronously. Failed sends retry after 1 minute, 5 minutes, 15 minutes, 1 hour, 6 hours, and 24 hours, reusing the same delivery ID; after the final failure, the server logs `notify.webhook_failed` without printing the endpoint URL or signing secret.

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
