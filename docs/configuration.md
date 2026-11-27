# Configuration

The server reads `server.yml` and the `yaml.d/` area files from its selected configuration directory and creates defaults on first startup. CLI binary distribution settings belong in `yaml.d/maintenance.yml` under the top-level `cli` key. Mode precedence is `--mode`, `MODE`, `server.mode`, then `production`. Debug precedence is explicit `--debug`, `DEBUG`, the `debug` mode default, then disabled; all boolean configuration and flag inputs use the project's shared truthy/falsy values.

| Setting | Default | Purpose |
|---|---:|---|
| `server.address` | `[::]` | Listener address; falls back to IPv4 wildcard when IPv6 is unavailable |
| `server.port` | random unused port in `64000`–`64999`; `80` in an application container | Listener port; the chosen value is saved to `server.yml` |
| `server.pidfile` | `true` | Create and validate the server PID file on startup; container runs skip PID files |
| `server.daemonize` | `false` | Detach the server process from the terminal when daemonization is enabled for the start mode |
| `server.fqdn` | empty | Public hostname supplied by the operator |
| `server.tor.binary` | auto-detected | Tor executable used for the app-owned hidden service; when no executable is available Tor remains disabled |
| `server.tor.expose` | `false` | Whether `/api/autodiscover` advertises the `.onion` address |
| `server.tor.contact_email` | empty | Email shown only on `.onion` responses; no clearnet contact is used as a fallback |
| `server.tor.use_network` | `false` | Routes app-determined outbound requests through the managed Tor SOCKS listener |
| `server.tor.max_circuits` | `32` | Tor pending-circuit limit, from 1 to 128 |
| `server.tor.startup_timeout` | `60` | Seconds to wait for the Tor control-port file, from 10 to 300 |
| `server.tor.circuit_timeout` | `60` | Tor circuit build timeout in seconds, from 10 to 300 |
| `server.tor.bootstrap_timeout` | `180` | Tor network bootstrap timeout in seconds, from 30 to 600 |
| `server.tor.safe_logging` | `true` | Scrubs sensitive data from Tor logs |
| `server.tor.max_streams_per_circuit` | `100` | Hidden-service stream limit, from 10 to 500 |
| `server.tor.close_circuit_on_stream_limit` | `true` | Closes a circuit when its stream limit is exceeded |
| `server.tor.bandwidth_rate` / `server.tor.bandwidth_burst` | `1 MB` / `2 MB` | Tor bandwidth rate and burst limits |
| `server.tor.max_monthly_bandwidth` | `100 GB` | Monthly Tor accounting limit, or `unlimited` |
| `server.tor.num_intro_points` | `3` | Hidden-service introduction points, from 3 to 10 |
| `server.tor.virtual_port` | `80` | Virtual port exposed by the hidden service, from 1 to 65535 |
| `server.i2p.enabled` | `true` in a container; otherwise `false` | Enables an I2P eepsite when an i2pd binary or SAM bridge is available; `--i2p` / `--i2p=false` overrides `I2P_ENABLED`, which overrides this setting |
| `server.i2p.expose` | `false` | Advertises a healthy `.b32.i2p` address in `/api/autodiscover` when enabled |
| `server.i2p.binary` | auto-detected | Optional path to the dedicated i2pd binary; when unavailable, the server may use the configured external SAM bridge |
| `server.i2p.sam_address` | `127.0.0.1:7656` | External SAMv3 bridge address used when no i2pd binary is available |
| `server.i2p.virtual_port` | `80` | Port exposed by the I2P eepsite |
| `server.i2p.inbound_length`, `server.i2p.outbound_length` | `3` | I2P tunnel hop counts, each from 0 to 7 |
| `server.i2p.inbound_quantity`, `server.i2p.outbound_quantity` | `5` | Parallel I2P tunnels, each from 1 to 16 |
| `server.i2p.signature_type` | `7` | SAM destination signature type |
| `server.i2p.bootstrap_timeout` | `300` | SAM session bootstrap timeout in seconds, from 30 to 600 |
| `server.mode` | `production` | Application mode; overridden by `MODE`, then `--mode` |
| `server.database.driver` | `sqlite` | Database backend: `sqlite`/`sqlite2`/`sqlite3` or remote `libsql`/`turso` |
| `server.database.url` | `{data_dir}/db/server.db` | SQLite path, or `libsql://`/`https://` remote connection URL; private-network `http://` is allowed only for loopback/private hosts; credential URLs belong in `secrets.yml` |
| `server.database.pool.max_open` | `25` | Maximum pooled database connections (SQLite read pool; remote pool) |
| `server.database.pool.max_idle` | `5` | Maximum idle database connections |
| `server.database.pool.max_lifetime` | `5m` | Maximum pooled connection lifetime |
| `server.database.pool.max_idle_time` | `1m` | Maximum idle connection time |
| `server.database.cache.enabled` | `false` | Enable the local SQLite read cache for a remote database |
| `server.database.cache.path` | `{data_dir}/db/cache.db` | Local SQLite cache file; resolved relative to the current data directory when omitted |
| `server.database.cache.sync_interval` | `30s` | Interval for complete transactional reconciliation from remote tables |
| `server.database.cache.offline_mode` | `false` | Allow stale, read-only cache reads while the remote is unavailable |
| `server.database.cache.max_age` | `1h` | Maximum age for normal cache reads before querying the remote |
| `server.update.branch` | `stable` | Release channel: `stable`, `beta`, or `devel`; aliases normalize to `devel` |
| `server.update.auto_install` | `false` | Let the scheduled update check stage eligible releases when the process started elevated |
| `server.update.defer_days` | `0` | Delay scheduled stable/beta updates by 0–365 days; manual checks ignore this value |
| `server.update.cli_min_version` | `""` | Operator-managed minimum client version advertised by autodiscovery |
| `cli.binary_download.dir` | `{data_dir}/cli` | Operator-provisioned companion CLI release files (`version.txt`, `sha256.txt`, and platform binaries) |
| `cli.binary_download.require_auth` | `false` | Require a valid Bearer token to download companion CLI binaries |
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
| `server.api_version` | `v1` | Version prefix used by the canonical API routes |
| `server.geoip.deny_countries` | `[]` | Deny requests from these ISO 3166-1 alpha-2 countries |
| `server.geoip.allow_countries` | `[]` | When nonempty, allow only these ISO country codes; takes precedence over `deny_countries` |
| `server.geoip.presets` | `{}` | Operator-authored named ISO country lists; presets are never applied automatically |
| `server.security.allowlist` | `[]` | Trusted IP/CIDR entries that bypass GeoIP country restrictions |
| `server.security.reports.retention_days` | `365` | Days to retain closed security reports; open reports are never purged |
| `server.trusted_proxies.additional` | `[]` | Additional proxy IPs/CIDRs whose forwarding headers are trusted |
| `server.rate_limit.read.requests` | `120` | Read requests per IP, or per valid token, per window |
| `server.rate_limit.read.window` | `60` | Read window in seconds |
| `server.rate_limit.write.requests` | `10` | Write requests per IP, or per valid token, per window |
| `server.rate_limit.write.window` | `60` | Write window in seconds |
| `server.rate_limit.health.requests` | `120` | Health requests per IP, or per valid token, per window |
| `server.rate_limit.health.window` | `60` | Health window in seconds |
| `server.rate_limit.strict.requests` | `10` | Sensitive token actions per IP or owner per window |
| `server.rate_limit.strict.window` | `900` | Sensitive-action window in seconds |
| `server.rate_limit.auth_failures.requests` | `5` | Failed token-authentication attempts per IP before lockout |
| `server.rate_limit.auth_failures.window` | `900` | Failed-authentication window in seconds |
| `server.rate_limit.global_burst` | `240` | Global per-IP or per-token ceiling per minute |
| `server.metrics.enabled` | `true` | Enables the internal Prometheus, Grafana, and Loki endpoints |
| `server.metrics.root.enabled` | `true` | Enables root aliases at `/metrics` and `/metrics/{service}` |
| `server.metrics.auth.allow_unauthenticated` | `false` | Disables service-token checks for firewalled metrics deployments only |
| `server.metrics.include_system` | `true` | Includes host CPU, memory, and disk metrics |
| `server.metrics.include_runtime` | `true` | Includes Go runtime gauges and GC counters |
| `server.metrics.loki.max_entries` | `1000` | Maximum recent structured HTTP log entries returned by Loki endpoint |
| `server.metrics.loki.max_age` | `1h` | Maximum age of entries returned by the Loki endpoint |
| `server.backup.retention.max_backups` | `1` | Number of daily full backups to retain (minimum 1) |
| `server.backup.retention.keep_weekly` | `0` | Additional Sunday full backups to retain |
| `server.backup.retention.keep_monthly` | `0` | First-of-month full backups to retain |
| `server.backup.retention.keep_yearly` | `0` | January 1 full backups to retain |
| `server.backup.retention.max_total_size` | `10%` | Backup-volume size cap, as a percentage or byte size; `0` disables it |
| `server.restore.max_archive_size` | `1G` | Maximum accepted backup archive size |
| `server.restore.max_expanded_size` | `4G` | Maximum total extracted size |
| `server.restore.max_files` | `10000` | Maximum archive file count |
| `server.restore.max_file_size` | `512M` | Maximum size of one extracted file |
| `server.restore.max_ratio` | `100` | Maximum expanded-to-compressed ratio |
| `server.scheduler.timezone` | `America/New_York` | IANA timezone used for scheduled backups |
| `server.scheduler.catch_up_window` | `1h` | Maximum age of a missed scheduled run to execute on startup |
| `server.scheduler.tasks.backup_daily.enabled` | `true` | Enables the daily backup task |
| `server.scheduler.tasks.backup_daily.schedule` | `0 2 * * *` | Daily backup schedule in the configured timezone |
| `server.scheduler.tasks.backup_hourly.enabled` | `false` | Enables hourly incremental backups |
| `server.scheduler.tasks.backup_hourly.schedule` | `@hourly` | Hourly incremental schedule |
| `server.scheduler.tasks.update_check.enabled` | `true` | Enables the daily operator-only update check |
| `server.scheduler.tasks.update_check.schedule` | `0 6 * * *` | Daily update check schedule |
| `server.notifications.email.smtp.host` | autodetected | SMTP relay hostname or IP; an empty value triggers first-run detection |
| `server.notifications.email.smtp.port` | `587` | SMTP relay port; automatic TLS uses implicit TLS on 465 and STARTTLS negotiation on other ports |
| `server.notifications.email.smtp.username` | empty | SMTP username; AUTH is used only after TLS is established |
| `server.notifications.email.smtp.password` | empty | SMTP password, stored as a secret and never logged |
| `server.notifications.email.smtp.tls` | `auto` | SMTP TLS mode: `auto`, `tls`, `starttls`, or `none` |
| `server.notifications.email.from.name` | application name | Sender display name |
| `server.notifications.email.from.email` | `no-reply@{fqdn}` | Sender address; defaults to `no-reply@localhost` when no FQDN is available |
| `server.notifications.email.reply_to` | empty | Optional Reply-To address |
| `server.notifications.email.events.*` | event-specific | Per-event operator email switches; SMTP must be configured and reachable |
| `server.contact.admin.email` | empty | Operator notification recipient; if unset, admin email notifications are skipped |
| `server.contact.security.email` | `security@{fqdn}` | Secondary security contact; an explicitly empty value falls back to the private admin recipient for internal delivery only and is omitted from public pages and `security.txt` |
| `server.contact.abuse.email` | empty | Public abuse report recipient; delivery falls back to general, then admin |
| `server.contact.general.email` | empty | Fallback recipient for abuse, DMCA, and privacy reports; also the DPO page's email fallback |
| `server.contact.dpo.name` | empty | Optional name shown on the data protection contact page |
| `server.contact.dpo.email` | empty | DPO email shown publicly; an empty value falls back to `server.contact.general.email`, never admin |
| `server.contact.{admin,security,abuse,general}.webhooks.*` | empty | Optional per-role HTTPS transports (`telegram`, `discord`, `slack`, `mattermost`, `pushover`, `gotify`, or `generic`); signing secrets are generated under the matching role in `secrets.yml` |
| `server.pages.contact.enabled` | `true` | Enables `/server/contact` and the instance-report line in security.txt |
| `server.pages.contact.captcha` | `simple` | CAPTCHA provider: self-hosted `simple`, `recaptcha`, or `hcaptcha`; external providers require site and secret keys and use the simple fallback over Tor and I2P |
| `server.pages.contact.captcha_site_key` | empty | Public widget key for reCAPTCHA or hCaptcha |
| `server.pages.contact.captcha_secret_key` | empty | Provider verification key; migrate it to `secrets.yml` and keep it out of responses and logs |
| `server.pages.contact.antispam.*` | project defaults | Fill and review timing, per-IP/email/global limits, markup/link rules, block lists, strike bans, and duplicate suppression |
| `server.privacy.data.sold` | `false` | Controls data-sharing statements and whether the conditional CCPA opt-out form is shown |
| `server.privacy.data.stored_on_server` | `true` | States whether data is stored on the instance |
| `server.privacy.data.sharing[]` | project defaults | Publicly listed data-sharing conditions, timing, and data categories |
| `server.privacy.retention.*` | project defaults | Public retention statement and whether export/deletion options are available |
| `server.privacy.consent.*` | project defaults | Consent banner visibility, sold/not-sold message, policy link, button labels, and preferences link |
| `server.privacy.cookies.{essential,preferences,analytics}.*` | project defaults | Cookie category descriptions and whether each nonessential category is offered |
| `server.privacy.third_party.services[]` | `[]` | Third-party services listed by name, purpose, data sent, and HTTPS privacy-policy URL |
| `server.maintenance.cleanup.disk_threshold` | `90` | Skip scheduled backups when disk use exceeds this percentage |
| `server.metrics.duration_buckets` | `[0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]` | Prometheus request duration histogram buckets in seconds |
| `server.metrics.size_buckets` | `[100, 1000, 10000, 100000, 1000000, 10000000]` | Prometheus request and response size histogram buckets in bytes |
| `server.logging.access_format` | `apache` | Reserved setting; currently has no effect |
| `server.logging.level` | `info` | Reserved setting; currently has no effect |
| `server.cors.allowed_origins` | `[*]` | Allowed browser origins; `[*]` permits any origin without credential headers, and `[""]` disables CORS |
| `server.cors.allow_credentials` | `true` | Allow credential headers only for explicitly matched origins |
| `server.cors.max_age` | `86400` | Preflight cache lifetime in seconds |
| `server.seo.*` | Project defaults | Canonical URLs, robots directives, generated SEO/social tags, verification codes, sitemap limits, and feed title/item count |
| `web.ui.theme` | `dark` | Default page theme (`dark`, `light`, or `auto`) |
| `web.ui.accent` | empty | Optional accessible `#RRGGBB` accent color |
| `web.robots.allow` | `[/]` | Public paths allowed by `robots.txt` and eligible for the sitemap |
| `web.robots.deny` | `[/debug]` | Paths disallowed in `robots.txt` and excluded from the sitemap |
| `web.announcements.enabled` | `true` | Enables the server-rendered site announcement banner |
| `web.announcements.messages[]` | `[]` | Ordered announcements with an ID, type, text, optional UTC window, and dismissal setting |
| `web.security.publish_pgp_key` | `true` | Publishes the project security public key at `/.well-known/pgp-key.asc` and includes it in security.txt while a keypair exists |
| `web.security.keyservers` | `[keys.openpgp.org]` | Keyserver host names used for project key publication |
| `web.security.report_url` | empty | Optional private vulnerability reporting URL used by the security overview |
| `web.security.policy` | empty | Optional Markdown text replacing the default disclosure policy |
| `server.branding.title` | `Zipcodes API` | Public application name |
| `server.branding.tagline` | `United States ZIP code lookup` | Public application tagline |
| `server.branding.description` | Search and geographic data description | Public application description |
| `server.branding.logo` | empty | Optional HTTPS/data URI or asset filename used as the site logo |
| `server.branding.favicon` | empty | Optional HTTPS/data URI or asset filename used as the favicon |
| `server.pages.help.content` | Project-specific help summary in Markdown | Appended to the server-rendered help page and returned as Markdown in its JSON API |
| `server.pages.help.mcp` | Project-specific MCP tools and usage examples | Markdown shown in the MCP help page and JSON mirror |

The server merges built-in defaults, `server.yml`, regular `yaml.d/*.yml` and `yaml.d/*.yaml` files in byte-wise filename order, then `secrets.yml`. Branding, SEO, and `web.ui` belong in `yaml.d/branding.yml`; other `web.*` settings and CORS belong in `yaml.d/web.yml`; rate limits and token policy belong in `yaml.d/security.yml`. Startup automatically splits area sections out of older `server.yml` files and preserves a `0600` backup. GeoIP data sources, database categories, storage location (`{data_dir}/security/geoip`), and weekly update schedule are fixed by the application; only the GeoIP enabled flag, country lists, and operator-authored presets are configurable. Example `yaml.d/branding.yml`:

The privacy page and `/api/{api_version}/server/privacy` mirror read `server.privacy.*`. `server.branding.content.privacy` is sanitized and appended after the built-in privacy sections. When `server.privacy.data.sold` is false (the default), no CCPA form is rendered and `POST /server/ccpa` returns `404`; when true, the privacy page exposes a CSRF-protected opt-out/opt-in form. The instance does not currently configure analytics tracking; list any other third-party service only when that service is actually enabled.

SMTP settings can be overridden by `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_TLS`, `SMTP_FROM_NAME`, and `SMTP_FROM_EMAIL`. The SMTP password is sensitive configuration and must remain in `secrets.yml` or the environment rather than being exposed in status output or logs.

Webhook URLs are configured under their recipient role in `yaml.d/contact.yml` or `server.yml`. At startup, the server generates a 32-byte signing secret for each configured transport when its matching `server.contact.{role}.webhooks.{transport}_secret` entry is empty, then stores that value in `secrets.yml` with mode `0600`. Share the secret only with the receiving endpoint. The server never includes webhook URLs or signing secrets in public responses or log messages.

While running, the configuration manager checks `server.yml`, `yaml.d/`, `secrets.yml`, and `assets/` every five seconds. A valid change reloads the merged configuration. A malformed file keeps its last successfully parsed contents; malformed first-time files are ignored. Changes to `server.port`, `server.address`, `server.daemonize`, or the SSL, database, and Tor settings are reported as requiring a restart. Check `/server/healthz` or run `zipcodes --status` to see pending restart state and the affected settings.

```yaml
server:
  branding:
    title: Zipcodes API
    tagline: United States ZIP code lookup
    description: Search ZIP codes and explore their geographic data through a public API.
  seo:
    title_separator: " - "
    opengraph:
      enabled: true
      type: website
    twitter:
      enabled: true
      card: summary_large_image
  cache:
    type: memory
    url: ""
    host: localhost
    port: 6379
    username: ""
    password: ""
    db: 0
    tls: false
    tls_skip_verify: true
    tls_ca: ""
    tls_cert: ""
    tls_key: ""
    pool_size: 10
    min_idle: 2
    timeout: 5s
    prefix: "zipcodes:"
    ttl: 1h
web:
  ui:
    theme: dark
    accent: ""
```

Header and footer content accepts sanitized HTML. The about, privacy, and terms entries accept `html`, `md`, or `txt`; each body may be a single string or a language map with a `default` fallback. `CONFIG_DIR`, `DATA_DIR`, and `LOG_DIR` select the respective directories; `DATABASE_DIR` selects the SQLite database directory, while `PORT` and `LISTEN` override the listener port and address. `PID_FILE` selects the PID file path unless `--pid` is provided; `server.pidfile: false` disables PID checks and writes, and container and Windows runs skip PID files. `MODE` and `DEBUG` follow the precedence described above, and `I2P_ENABLED` overrides `server.i2p.enabled` unless `--i2p` is explicitly provided. `DOMAIN` supplies a public host used when generating absolute URLs if `server.fqdn` is empty. Command-line flags also accept configuration, data, log, address, and port values. On Linux, root defaults are `/etc/apimgr/zipcodes`, `/var/lib/apimgr/zipcodes`, and `/var/log/apimgr/zipcodes`; a non-root process uses XDG paths. macOS and Windows use their standard application support/data directories.

SQLite is the default. Its URL is a local path and defaults to `{data_dir}/db/server.db` on native installs and `/data/db/sqlite/server.db` in containers. `DATABASE_DIR` overrides the SQLite database directory on every start. SQLite uses separate read-only and single-writer pools, WAL, `synchronous=NORMAL`, and a 10 second busy timeout. Keep the database on local storage; network and FUSE filesystems are unsupported.

For a remote libSQL/Turso database, set `server.database.driver: libsql` and use either `libsql://host` with `?authToken=...` in the URL or an `https://host` URL plus `server.database.token` in `secrets.yml`. Plain `http://` is accepted only for loopback or private hosts, including hostnames that resolve exclusively to private addresses. The URL takes precedence when it already contains `authToken`. `DATABASE_URL` initializes the saved connection on first run only; set `DATABASE_DRIVER` when the URL does not select a driver by scheme. Remote tables and indexes use the `srv_` prefix. The URL token and separate token are secrets: never include them in logs or public configuration output. Remote contents are managed and backed up by the provider. The optional SQLite read cache is disabled by default; enable `server.database.cache.enabled` to use a separate local cache file. The default path is `{data_dir}/db/cache.db`, the sync interval is 30 seconds, and the maximum cache age is one hour. Set `path`, `sync_interval`, `offline_mode`, and `max_age` under `server.database.cache` to change those values. Setting `offline_mode: true` allows stale, read-only cache results when the remote is unavailable. Cached data is read-only: writes go to the remote database first, and remote state wins each cache reconciliation. When enabled, local backups include the cache snapshot as `server.db`; restoring seeds an empty remote database from that snapshot but leaves a populated remote database unchanged.

Sensitive server credentials are stored separately in `secrets.yml` beside `server.yml`, not in `server.yml`. It contains `server.token` (the operator token), `server.database.token` (remote libSQL token when kept separate from the URL), `server.cache.password` and credential-bearing `server.cache.url`, `server.security.encryption_key` plus `server.security.encryption_key_version` (the key and version used to encrypt owner-token values and security-report contact metadata), and separate `server.metrics.auth.tokens.prometheus`, `.grafana`, and `.loki` credentials. The server creates these metric service tokens using the `mtr_` prefix and writes the file with owner-only permissions (`0600`); on first startup it prints the generated operator token once, but never prints metrics tokens. Protect and back up the secret file as sensitive data. Losing the encryption key prevents stored owner tokens from being revealed and prevents security-report status notices from being sent.

Backup encryption uses the operator-managed `server.backup.encryption_password` key in `secrets.yml`; it is never generated automatically. Without a password, backups are unencrypted and startup logs a warning. When any compliance standard is enabled, backup creation is refused until a password is configured. Encrypted archives use a cleartext `BKUP1` header authenticated as AES-GCM additional data; it contains the per-archive salt, Argon2id parameters, nonce, key ID, and the incremental base filename when applicable. Scheduled backups use the header key ID to recognize a password change without decrypting old archives: that run writes a timestamped full and skips the incremental, while retaining the previous full and incremental pair until the next run replaces the incremental. Keep old passwords to verify or restore older archives. Backups are written under `BACKUP_DIR` when set, otherwise to the platform backup location with a data-directory fallback. `--backup` overrides that directory. Run `zipcodes --maintenance backup [filename]` for a manual backup, `zipcodes --maintenance backup test [file]` to verify one, and `zipcodes --maintenance restore [file]` for a CLI-only restore. When omitted, the restore and backup test select the full archive with the latest date in its filename; modification time breaks date ties. Restore prompts for an encrypted archive password without echoing it and applies restore size limits. A nonempty restore requires root confirmation or the service account plus the operator token. `--include-ssl` and `--include-data` add those optional content groups to a manual backup.

Every backup includes `app_secrets.json`, an export of the server-owned `app_secrets` rows, including `installation_secret` and previous-key rows. Restore replaces those rows after the verified archive files are installed. This file supports PGP key recovery when a remote database is not included in the archive. `backup test` verifies a backed-up PGP keypair against this secret snapshot without modifying live files, including for remote-database backups without a local `server.db`.

Retention is applied at startup and after backups. The daily backup scheduler defaults to 02:00 and the `security_cleanup` task defaults to 03:30 in `server.scheduler.timezone`; hourly incrementals are disabled by default. Security cleanup removes closed reports after `server.security.reports.retention_days` and never purges open reports. A scheduled backup is skipped when free space is below twice the previous full backup size or disk use is above `server.maintenance.cleanup.disk_threshold`. In remote mode the database provider owns remote backups; restore seeding is limited to an empty remote database.

Server configuration is operator-controlled through files and command-line options. Visitor theme and language preferences are managed by the server preferences pages and consent rules in AI.md. The CLI client has a separate per-user `cli.yml` with `server.primary`, `display.mode` (`auto` by default), and `theme.mode` (`auto` by default). Its display mode is detected from the environment; the client does not provide a `--tui` mode flag.

Rate limiting is always active; its request and window limits are configurable, but there is no `enabled` switch. Valid bearer tokens use token-specific buckets, while unauthenticated or invalid-token requests use the client IP. Repeated failed bearer authentication locks out the IP for the configured window, even if a later request uses a valid token. Sensitive token actions have a separate strict limit. Counters use SQLite by default. Optional Valkey/Redis backends are selected with `server.cache.type` and either `server.cache.url` or host/port settings; credential-bearing URLs and passwords are stored in `secrets.yml`. A configured but unreachable external cache falls back to in-process sliding-window counters and reports `checks.cache: degraded` in health responses. `CACHE_URL` initializes the saved cache connection on first run only. The default database path is `db/server.db` beneath the data directory. `server.tokens.cookie_max_age` accepts a Go duration such as `24h` or `1y`; the latter is treated as 365 days. GeoIP is fail-open: missing lookup data and private/internal addresses do not block requests. Country restrictions are disabled by default and, when configured, run after per-IP rate limiting as a deny-only signal; they never replace authentication or authorization. Presets only reuse operator-authored lists and are not applied automatically. See [Security](security.md).
