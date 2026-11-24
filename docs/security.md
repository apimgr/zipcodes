# Security

ZIP lookup routes are public and read-only; they do not require accounts or API keys. Owner tokens are optional credentials for associating caller activity with an anonymous owner and managing that owner's linked tokens. There are no user accounts, profiles, or passwords. ZIP input and radius/coordinate values are validated, SQL uses bound parameters, and CORS defaults to `*` without credentialed requests. Forwarded client-IP headers are honored only from trusted peers.

Per-IP request limits default to 120 read requests/minute, 10 non-read requests/minute, 120 health requests/minute, and a 240 request/minute per-IP ceiling across endpoint classes. Counters are stored in `server.db` and keyed by a SHA-256 digest of the client IP and endpoint class; expired events are removed as requests arrive.

The SQLite database uses separate read-only and single-writer pools with WAL, a 10 second busy timeout, and a shutdown WAL checkpoint. Keep the database on local storage; network and FUSE filesystems are unsupported. This build does not support remote libSQL/Turso databases.

`/server/healthz`, `/server/privacy`, `/server/terms`, `/server/about`, `/server/help`, preferences, and the API documentation pages are public. The root `/healthz` alias is disabled by default and, when enabled, uses the same health handler and content negotiation. Browser preference and token-management mutations use same-origin POST forms with CSRF tokens. API token-management routes require the `Authorization: Bearer` header with the owner's primary token, except primary-token issuance, which is public and issues a new anonymous owner. API routes do not accept owner cookies. Linked tokens cannot manage tokens. `/api/{api_version}/me` uses the caller IP from the TCP peer or a configured trusted proxy. GeoIP skips private/reserved addresses and fails open when lookup data is unavailable. Country and ASN data use PDDL iptoasn.com files; city data uses DB-IP under CC BY 4.0. Every downloaded MMDB is checked against its release checksum before replacement, and the scheduler refreshes enabled GeoIP weekly.

The `/server/dpo` page and its JSON mirror expose `server.contact.dpo.email`, falling back to `server.contact.general.email`. They never expose the private `server.contact.admin.email`; when neither public address is configured, the page links to `/server/security` instead.

Owner and linked tokens use the `own_` prefix. The database stores a SHA-256 lookup hash and an AES-256-GCM encrypted copy so owners can reveal their tokens from the token page. Revealed values are not included in token-list responses. A token's raw value is returned by issue, reveal, or rotate actions; protect it like a password. Rotating a primary token preserves the owner and linked tokens. Revoking the primary deletes that owner and all associated tokens and resources. See [API reference](api.md) for the routes.

The operator token uses the `adm_` prefix and is stored in `secrets.yml`, separate from `server.yml`. It is generated on first startup and printed to the process console once. The file is created with `0600` permissions. The token-encryption key is stored in the same file; losing it prevents decryption of stored token values.

Backups can contain server configuration, secrets, database contents, and optional application data. Store the backup directory with restricted permissions and configure `server.backup.encryption_password` in `secrets.yml` to encrypt archives. The password is operator-managed and is not generated. Enabling any `server.compliance` standard blocks backup creation until encryption is configured. Restore is available only through the local CLI; it prompts for encrypted archive passwords without echoing them. Restoring over a populated database requires typed confirmation when running as root, or the operator token when running as the service account. Other users cannot overwrite a populated database. Restore enforces archive expansion limits and validates the SQLite database before installation. Retain old encryption passwords for archives created before a password change; rotation does not re-encrypt existing files.

Responses include content-type, framing, referrer, content-security, cross-origin, permissions, request ID, and permitted cross-domain policy headers. HSTS is returned only for direct TLS requests. Same-origin state-changing browser forms require a CSRF token.

## TLS certificates

When TLS is enabled, the server selects certificates from the configured certificate files, the system Let's Encrypt directories, the app-managed `{config_dir}/ssl/letsencrypt/{fqdn}/` directory, then `{config_dir}/ssl/local/{fqdn}/`. Certificates must be readable, unexpired, and match the configured host. Only app-managed Let's Encrypt certificates renew automatically, seven days before expiry; system and local certificates remain operator-managed. HTTP-01 and TLS-ALPN-01 use the server's existing HTTP and HTTPS listeners. DNS-01 supports lego providers; provider credentials are encrypted at rest with the server encryption key, and values are never logged or returned. See [configuration](configuration.md) for setup details.

## Vulnerability reporting

Report project-source vulnerabilities privately through [GitHub Security Advisories](https://github.com/apimgr/zipcodes/security/advisories/new). A deployed instance can also accept reports through `/server/security` and the `/server/contact?security_id=...` form. Security reports are encrypted at rest with the configured AES-256-GCM key; SMTP delivery attaches the encrypted report, and an acknowledgment contains the tracking ID and a status URL. When a researcher submits an armored OpenPGP public key or an HTTPS key URL, the acknowledgment is encrypted to that key. Without working SMTP, the report is still stored but no notification or acknowledgment is sent. The status URL exposes triage state only, not report content.

`/.well-known/security.txt` is the sole security.txt route. It includes the GitHub advisory URL, a rotating instance-report URL when contact reporting is enabled, a configured security contact when available, an annual expiration time, `Preferred-Languages: en`, and request-derived policy and `Canonical` URLs. `/.well-known/pgp-key.asc` returns `404`; project-level keypair management and encrypted-at-rest PGP reports are not implemented yet.

## Public and discovery endpoints

- `/server/healthz` is public; its root `/healthz` alias is disabled by default and can be enabled with `server.healthz.root.enabled`. `/api/v1/server/healthz` and `/api/healthz` are machine endpoints. See [API reference](api.md) for content negotiation.
- `/.well-known/llms.txt` and `/llms.txt` list the public API and documentation routes without private configuration.
- `/.well-known/security.txt` is the canonical security contact file.
- `robots.txt` and `sitemap.xml` expose public crawling guidance and public HTML routes.
- Unknown `/.well-known/*` paths return `404`. No WebFinger, OAuth/OIDC metadata, app association, or MTA-STS endpoint is enabled.
