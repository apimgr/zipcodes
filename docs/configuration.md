# Configuration

The server reads `server.yml` from its selected configuration directory and creates defaults on first startup.

| Setting | Default | Purpose |
|---|---:|---|
| `server.address` | `0.0.0.0` | Listener address |
| `server.port` | selected at startup | Listener port |
| `server.fqdn` | empty | Public hostname supplied by the operator |
| `server.mode` | `production` | Application mode |
| `server.update_branch` | `stable` | Reserved setting; update commands currently select their branch from the command argument |
| `server.max_radius_miles` | `100` | Maximum accepted radius |
| `server.baseurl` | `/` | URL prefix used for browser routes and embedded assets when the application is mounted below a reverse proxy path |
| `server.healthz.root.enabled` | `false` | Mount `/healthz` as a root-level alias to the `/server/healthz` frontend |
| `server.tokens.max_linked` | `25` | Maximum active linked API tokens per owner; zero disables the limit |
| `server.tokens.cookie_max_age` | `1y` | Lifetime of the owner-token browser cookie |
| `server.tokens.idle_purge_days` | `0` | Reserved setting; currently accepted but has no effect |
| `server.limits.max_body_size` | `10MB` | Maximum request body size |
| `server.limits.read_timeout` | `30s` | HTTP read timeout |
| `server.limits.write_timeout` | `30s` | HTTP and handler timeout |
| `server.limits.idle_timeout` | `120s` | Keep-alive idle timeout |
| `server.compression.enabled` | `true` | Enables configured response compression |
| `server.compression.level` | `5` | Compression level from 1 to 9 |
| `server.compression.types` | HTML, CSS, JavaScript, JSON, XML | MIME types eligible for compression |
| `server.geoip.enabled` | `true` | Downloads and opens GeoIP databases |
| `server.geoip.dir` | `security/geoip` | GeoIP database directory relative to data directory |
| `server.geoip.update_interval` | `weekly` | GeoIP database refresh interval |
| `server.geoip.city` | `true` | Download IPv4 and IPv6 city databases |
| `server.geoip.country` | `true` | Download the country-code database |
| `server.geoip.asn` | `true` | Download ASN organization data |
| `server.trusted_proxies.additional` | `[]` | Additional proxy IPs/CIDRs whose forwarding headers are trusted |
| `server.rate_limit.enabled` | `true` | Enables request abuse limits |
| `server.rate_limit.read.requests` | `120` | Read requests per IP per window |
| `server.rate_limit.read.window` | `60` | Read window in seconds |
| `server.rate_limit.write.requests` | `10` | Write requests per IP per window |
| `server.rate_limit.write.window` | `60` | Write window in seconds |
| `server.rate_limit.health.requests` | `120` | Health requests per IP per window |
| `server.rate_limit.health.window` | `60` | Health window in seconds |
| `server.rate_limit.global_burst` | `240` | Global per-IP ceiling per minute |
| `server.metrics.enabled` | `false` | Reserved setting; metrics collection is not currently implemented |
| `server.metrics.endpoint` | `/metrics` | Reserved setting; no metrics endpoint is currently served |
| `server.metrics.include_system` | `true` | Reserved setting; currently has no effect |
| `server.metrics.include_app` | `true` | Reserved setting; currently has no effect |
| `server.logging.access_format` | `apache` | Reserved setting; currently has no effect |
| `server.logging.level` | `info` | Reserved setting; currently has no effect |
| `web-security.cors` | `*` | CORS origin for this public API |
| `web-ui.theme` | `dark` | Default page theme |
| `web-ui.logo` | empty | Reserved setting; the pages currently use the application name as text |
| `web-ui.favicon` | empty | Reserved setting; the pages currently use the bundled favicon |
| `web-ui.notifications.enabled` | `true` | Reserved setting; announcements are not currently rendered |
| `web-ui.notifications.announcements` | `[]` | Reserved announcement messages; currently not rendered |
| `web-robots.allow` | `/`, `/api` | Paths listed as allowed in `robots.txt` |
| `web-robots.deny` | `/debug` | Paths listed as disallowed in `robots.txt` |
| `web-security.admin` | empty | Reserved setting; the current `security.txt` always points to GitHub Security Advisories |
| `branding.title` | `Zipcodes API` | Public application name |
| `branding.tagline` | `United States ZIP code lookup` | Public application tagline |
| `branding.description` | Search and geographic data description | Public application description |
| `server.pages.help.mcp` | Project-specific MCP tools and usage examples | Markdown shown in the MCP help page and JSON mirror |

The server reads `server.yml` in its configuration directory. `CONFIG_DIR`, `DATA_DIR`, and `LOGS_DIR` select the respective directories; `DB_PATH` overrides the SQLite database path, and `PORT` and `ADDRESS` override listener settings. `DOMAIN` supplies a public host used when generating absolute URLs if `server.fqdn` is empty. Command-line flags also accept configuration, data, log, address, and port values. On Linux, root defaults are `/etc/apimgr/zipcodes`, `/var/lib/apimgr/zipcodes`, and `/var/log/apimgr/zipcodes`; a non-root process uses XDG paths. macOS and Windows use their standard application support/data directories.

Sensitive server credentials are stored separately in `secrets.yml` beside `server.yml`, not in `server.yml`. It contains `server.token` (the operator token) and `server.security.encryption_key` plus `server.security.encryption_key_version` (the key and version used to encrypt owner-token values). The server creates this file with owner-only permissions (`0600`); on first startup it prints the generated operator token once. Protect and back up the secret file as sensitive data. Losing the encryption key prevents stored owner tokens from being revealed.

Server configuration is operator-controlled through files and command-line options. Browser theme preference is stored in a visitor cookie and does not change server configuration. The preferences page currently saves theme only; its language field is not yet writable. The CLI client has a separate `cli.yml` containing `server.primary`.

Rate limit counters are persisted in SQLite. The default database path is `db/server.db` beneath the data directory. `server.tokens.cookie_max_age` accepts a Go duration such as `24h` or `1y`; the latter is treated as 365 days. GeoIP is fail-open: a failed download is logged and does not prevent ZIP lookups. See [Security](security.md).
