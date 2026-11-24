# Security

ZIP lookup routes are public and read-only; they do not require accounts or API keys. ZIP input and radius/coordinate values are validated, SQL uses bound parameters, and CORS defaults to `*` without credentialed requests. Forwarded client-IP headers are honored only from trusted peers.

Per-IP request limits default to 120 read requests/minute, 10 non-read requests/minute, 120 health requests/minute, and a 240 request/minute per-IP ceiling across endpoint classes. Counters are stored in `server.db` and keyed by a SHA-256 digest of the client IP and endpoint class; expired events are removed as requests arrive.

`/server/healthz`, `/server/privacy`, `/server/terms`, `/server/about`, `/server/help`, and the API documentation pages are public. The root `/healthz` alias is disabled by default and, when enabled, uses the same health handler and content negotiation. Theme changes require a CSRF-protected same-origin form. `/api/v1/me` uses the caller IP from the TCP peer or a configured trusted proxy. GeoIP is skipped for private/reserved addresses and fails open to an IP-only response when location data is unavailable. City records come from DB-IP, and country/ASN data comes from NRO, using the CC BY 4.0 ip-location-db MMDB files. The update scheduler refreshes the enabled databases at the configured interval.

Responses include content-type, framing, referrer, content-security, cross-origin, permissions, request ID, and permitted cross-domain policy headers. HSTS is returned only for direct TLS requests. Same-origin state-changing browser forms require a CSRF token.

## Vulnerability reporting

Report a vulnerability privately through [GitHub Security Advisories](https://github.com/apimgr/zipcodes/security/advisories/new). `/.well-known/security.txt` (also `/security.txt`) advertises the configured `web-security.admin` mail contact when present, otherwise the GitHub advisory form. It includes an annual expiration time, `Preferred-Languages: en`, and a request-derived `Canonical` URL. The project has no PGP key endpoint; `/.well-known/pgp-key.asc` returns `404`.

## Public and discovery endpoints

- `/server/healthz` is public; its root `/healthz` alias is disabled by default and can be enabled with `server.healthz.root.enabled`. `/api/v1/server/healthz` and `/api/healthz` are machine endpoints. See [API reference](api.md) for content negotiation.
- `/.well-known/llms.txt` and `/llms.txt` list the public API and documentation routes without private configuration.
- `/.well-known/security.txt` is the canonical security contact file.
- `robots.txt` and `sitemap.xml` expose public crawling guidance and public HTML routes.
- Unknown `/.well-known/*` paths return `404`. No WebFinger, OAuth/OIDC metadata, app association, or MTA-STS endpoint is enabled.
