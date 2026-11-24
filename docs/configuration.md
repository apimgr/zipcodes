# Configuration

The server reads `server.yml` from its selected configuration directory and creates defaults on first startup.

| Setting | Default | Purpose |
|---|---:|---|
| `server.address` | `0.0.0.0` | Listener address |
| `server.port` | selected at startup | Listener port |
| `server.fqdn` | empty | Public hostname supplied by the operator |
| `server.mode` | `production` | Application mode |
| `server.update_branch` | `stable` | Release branch selected by the update command |
| `server.max_radius_miles` | `100` | Maximum accepted radius |
| `server.baseurl` | `/` | URL prefix used for browser routes and embedded assets when the application is mounted below a reverse proxy path |
| `server.healthz.root.enabled` | `false` | Mount `/healthz` as a root-level alias to the `/server/healthz` frontend |
| `server.tokens.max_linked` | `25` | Maximum active linked API tokens per owner; zero disables the limit |
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
| `server.metrics.enabled` | `false` | Enables metrics collection |
| `server.metrics.endpoint` | `/metrics` | Metrics endpoint setting |
| `server.metrics.include_system` | `true` | Include system metrics |
| `server.metrics.include_app` | `true` | Include application metrics |
| `server.logging.access_format` | `apache` | Access log format |
| `server.logging.level` | `info` | Log level |
| `web-security.cors` | `*` | CORS origin for this public API |
| `web-ui.theme` | `dark` | Default page theme |
| `web-ui.logo` | empty | Optional configured logo value |
| `web-ui.favicon` | empty | Optional configured favicon value |
| `web-ui.notifications.enabled` | `true` | Enables configured announcement notices |
| `web-ui.notifications.announcements` | `[]` | Announcement messages |
| `web-robots.allow` | `/`, `/api` | Paths listed as allowed in `robots.txt` |
| `web-robots.deny` | `/debug` | Paths listed as disallowed in `robots.txt` |
| `web-security.admin` | empty | Security contact used by `security.txt` when configured |
| `branding.title` | `Zipcodes API` | Public application name |
| `branding.tagline` | `United States ZIP code lookup` | Public application tagline |
| `branding.description` | Search and geographic data description | Public application description |

The server reads `server.yml` in its configuration directory. `CONFIG_DIR`, `DATA_DIR`, and `LOGS_DIR` select the respective directories; `DB_PATH` overrides the SQLite database path, and `PORT` and `ADDRESS` can override listener settings. Command-line flags also accept configuration, data, log, address, and port values. On Linux, root defaults are `/etc/apimgr/zipcodes`, `/var/lib/apimgr/zipcodes`, and `/var/log/apimgr/zipcodes`; a non-root process uses XDG paths. macOS and Windows use their standard application support/data directories.

Sensitive server credentials are stored separately in `secrets.yml` beside `server.yml`, not in `server.yml`. The server creates this file with owner-only permissions (`0600`); on first startup it prints the generated operator token once. Protect and back up the secret file as sensitive data.

Server configuration is operator-controlled through files and command-line options. Browser theme preference is stored in a visitor cookie and does not change server configuration. The preferences page currently saves theme only; its language field is not yet writable. The CLI client has a separate `cli.yml` containing `server.primary`.

Rate limit counters are persisted in SQLite. The default database path is `db/server.db` beneath the data directory. GeoIP is fail-open: a failed download is logged and does not prevent ZIP lookups. See [Security](security.md).
