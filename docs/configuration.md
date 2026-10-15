# Configuration

The server reads `server.yml` from its selected configuration directory and creates defaults on first startup.

| Setting | Default | Purpose |
|---|---:|---|
| `server.address` | `0.0.0.0` | Listener address |
| `server.port` | selected at startup | Listener port |
| `server.max_radius_miles` | `100` | Maximum accepted radius |
| `server.baseurl` | `/` | Reserved path-prefix setting; route and asset prefixing remains incomplete |
| `server.limits.max_body_size` | `10MB` | Maximum request body size |
| `server.limits.read_timeout` | `30s` | HTTP read timeout |
| `server.limits.write_timeout` | `30s` | HTTP and handler timeout |
| `server.limits.idle_timeout` | `120s` | Keep-alive idle timeout |
| `server.compression.enabled` | `true` | Enables configured response compression |
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
| `web-security.cors` | `*` | CORS origin for this public API |

Rate limit counters are persisted in the configured SQLite database (`server.db` by default). GeoIP is fail-open: a failed download is logged and does not prevent ZIP lookups. See [Security](security.md).
