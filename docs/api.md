# API reference

ZIP lookup routes are public and read-only. The default REST API base path is `/api/{api_version}` (default `/api/v1`); operators can configure `server.api_version`, and clients discover the active version from `/api/autodiscover`. Separate preference and owner-token routes manage visitor settings and anonymous owner credentials; their access requirements are described below. Interactive OpenAPI documentation is served at `/server/docs/swagger`; the OpenAPI JSON is `/api/{api_version}/server/swagger` (also `/api/swagger`). The GraphQL interface is `/server/docs/graphql`; submit queries to `/api/{api_version}/server/graphql` or `/api/graphql`. The old root `/graphql` path is not served.

## ZIP code routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/{api_version}/zipcodes/{code}` | Look up one five-digit ZIP code |
| GET | `/api/{api_version}/zipcodes/{code}.txt` | Plain text ZIP lookup |
| GET | `/api/{api_version}/zipcodes/search?q={query}` | Search ZIP, city, state, or city/state |
| GET | `/api/{api_version}/zipcodes/autocomplete?q={query}` | Return suggestions |
| GET | `/api/{api_version}/zipcodes/city/{city}` | List ZIP codes for a city |
| GET | `/api/{api_version}/zipcodes/state/{state}` | List ZIP codes for a state abbreviation |
| GET | `/api/{api_version}/zipcodes/nearby?zip={zip}&radius={miles}` | Search around an origin ZIP |
| GET | `/api/{api_version}/zipcodes/nearby?lat={lat}&lng={lng}&radius={miles}` | Search around coordinates |
| GET | `/api/{api_version}/zipcodes/stats` | Return dataset counts |
| GET | `/api/{api_version}/zipcodes/stats.txt` | Return dataset counts as text |
| GET | `/api/{api_version}/zipcodes/random` | Return a random ZIP record |
| GET | `/api/{api_version}/zipcodes/random.txt` | Return a random record as text |
| GET | `/api/{api_version}/zipcodes.json` | Download the full embedded dataset |

Radius must be greater than zero and may not exceed the server's configured maximum (100 miles by default). Coordinate inputs are bounded to latitude -90..90 and longitude -180..180. ZIP codes must contain exactly five digits.

## GeoIP routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/{api_version}/me` | Return the caller IP and available location data |
| GET | `/api/{api_version}/geoip?ip={address}` | Look up an IP address |
| GET | `/api/{api_version}/geoip.txt?ip={address}` | Text form of an IP lookup |
| POST | `/api/{api_version}/geoip/batch` | Look up a bounded batch of IP addresses |

GeoIP lookups return IP-only results when data is unavailable. Private/internal addresses are never looked up or country-blocked. Optional country restrictions are disabled by default and act only as a deny-only risk signal; ZIP lookup data is unaffected.

## Server information routes

The matching browser pages include `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/security`, `/server/security/policy`, `/server/security/thanks`, `/server/dpo`, `/server/help/mcp`, and the preferences and token pages. Corresponding versioned API routes are listed in the OpenAPI document; MCP is a machine-only endpoint and has no paired page requirement.

The server also publishes `/robots.txt`, `/sitemap.xml`, `/favicon.ico`, and `/.well-known/security.txt`. The sitemap lists public HTML pages only; it does not list API endpoints.

`/feed.xml` (Atom 1.0) and `/rss.xml` (RSS 2.0) are always available. ZIP records do not have individual public pages or creation times, so both feeds contain no entries. They use request-derived absolute URLs and support ETag conditional requests. `/.well-known/llms.txt` and `/llms.txt` list the public API and documentation endpoints. `/.well-known/pgp-key.asc` returns 404 because this deployment has no security-report keypair.

HTML pages support server-rendered search and require no JavaScript. Theme changes use a CSRF-protected same-origin form and are stored in a first-party cookie.

## MCP

The remote Streamable HTTP endpoint is `/api/{api_version}/server/mcp`; `/api/mcp` and `/mcp` are direct aliases. It accepts authenticated JSON-RPC 2.0 `POST` requests using the same owner or operator Bearer tokens as the REST API. It targets protocol revision `2026-07-28`, has no initialize handshake or session, and requires `MCP-Protocol-Version`, `Mcp-Method`, and `Accept: application/json, text/event-stream`. Tool names, descriptions, and input schemas are generated from the OpenAPI operations. Token-management operations and the MCP endpoint itself are excluded from the tool list. On a running instance, `/server/help/mcp` provides the client setup guide.

## Contact and security reporting

`GET /api/{api_version}/server/contact` describes the available report types, form fields, fill token, and CAPTCHA challenge. `POST` uses a submit, review, and confirm flow; it does not deliver on the first request. Abuse, DMCA, and privacy submissions require SMTP. Security reports use a rotating `security_id`, are encrypted with the server AES-256-GCM key before database storage, and still succeed when SMTP is unavailable. The response returns a `tracking_id`; the status URL requires a separate one-use-per-day token sent only in the acknowledgment email. `/.well-known/security.txt` is the only security.txt route.

## Preferences

The browser preferences page is `/server/preferences`; its form currently saves the `theme` value (`dark`, `light`, or `auto`). The page displays the detected language but language changes are not implemented yet. `GET /api/{api_version}/server/preferences` returns the current `theme`, `lang`, and `prefs` values. `POST /api/{api_version}/server/preferences` accepts JSON or form values and applies valid theme changes; unsupported values, including `lang`, are returned in `ignored`.

Preference export/import is available at `/server/preferences/sync` and `/server/preferences/sync/import`. The sync page exports theme and language values as a URL and a base64url code. Import currently applies theme only; language is included in the data but is ignored. The API equivalents are `GET /api/{api_version}/server/preferences/sync` and `GET /api/{api_version}/server/preferences/sync/import`.

## Owner tokens

Owner tokens identify an anonymous owner; they do not create user accounts. The primary token can manage linked tokens. Linked tokens have the same owner scope but cannot manage tokens. API token routes accept a primary owner token in `Authorization: Bearer …`; cookies are not used for these API routes. The web token page uses the owner cookie and CSRF-protected POST forms.

The server-rendered token page supports these form actions:

| Method | Route | Purpose |
|---|---|---|
| POST | `/server/preferences/token/create` | Issue a primary token for a new owner |
| POST | `/server/preferences/token/restore` | Restore the browser using a primary token |
| POST | `/server/preferences/token/forget` | Clear this browser's owner-token cookie without revoking the token |
| POST | `/server/preferences/token/reveal` | Show the primary token |
| POST | `/server/preferences/token/rotate` | Rotate the primary token |
| POST | `/server/preferences/token/revoke` | Revoke the primary token and delete its owner |
| POST | `/server/preferences/tokens` | Create a linked token |
| POST | `/server/preferences/tokens/{id}/{action}` | Reveal, rotate, or revoke a linked token |

| Method | Route | Access | Purpose |
|---|---|---|---|
| GET | `/server/preferences/tokens` | Owner cookie for management; otherwise shows create/restore forms | Token management page |
| POST | `/api/{api_version}/server/preferences/tokens/primary` | No token required | Create an owner and issue its primary token; returns `201`, or `409` if a valid token is already supplied |
| GET | `/api/{api_version}/server/preferences/tokens` | Primary owner Bearer token | List tokens with values masked/omitted |
| POST | `/api/{api_version}/server/preferences/tokens/primary/reveal` | Primary owner Bearer token | Reveal the primary token |
| POST | `/api/{api_version}/server/preferences/tokens/primary/rotate` | Primary owner Bearer token; JSON `confirm` must match its last four characters | Replace the primary token |
| DELETE | `/api/{api_version}/server/preferences/tokens/primary` | Primary owner Bearer token; JSON `confirm` must match its last four characters | Delete the owner and all its tokens/resources |
| POST | `/api/{api_version}/server/preferences/tokens` | Primary owner Bearer token; JSON `name` and optional `expiry` | Create a linked token |
| POST | `/api/{api_version}/server/preferences/tokens/{id}/reveal` | Primary owner Bearer token | Reveal a linked token |
| POST | `/api/{api_version}/server/preferences/tokens/{id}/rotate` | Primary owner Bearer token; JSON `confirm` must match the token's last four characters | Replace a linked token |
| POST | `/api/{api_version}/server/preferences/tokens/{id}/update` | Primary owner Bearer token; JSON `name` and `expiry` | Rename a linked token or change its expiry |
| DELETE | `/api/{api_version}/server/preferences/tokens/{id}` | Primary owner Bearer token | Revoke a linked token |

Linked-token expiry accepts `never`, `7d`, `30d`, `90d`, `1y`, or a future `YYYY-MM-DD` date. The default limit is 25 active linked tokens per owner. Raw token values are returned only by issue, reveal, and rotate actions; save them securely. The primary token may be restored in a browser from `/server/preferences/tokens`. Revoking the primary deletes the owner and its associated data; rotating it preserves the owner and linked tokens.

`/server/healthz` is the canonical health frontend: browsers receive HTML, HTTP tools receive human-readable text, and `Accept: application/json` receives JSON. `/api/{api_version}/server/healthz` and `/api/healthz` are machine endpoints; they return JSON by default and support text with `.txt` or `Accept: text/plain`. The root `/healthz` alias is disabled by default and can be enabled with `server.healthz.root.enabled: true`; it uses the same content negotiation as `/server/healthz`. Health responses are bare and are not wrapped in the standard API envelope. Health checks are cached for at most 60 seconds. A live configuration change that requires a restart adds `pending_restart` and `restart_reason` to health output and changes the status to `restart_required` while the service remains available. Request totals are persisted in batches of at most 60 seconds, requests are grouped into rolling hourly buckets, and health and metrics scrapes are excluded from request totals. GeoIP checks report enabled datasets individually; unavailable non-critical collectors make the overall status `degraded`. The OpenAPI and GraphQL endpoints are described above. Other JSON API responses use `ok` and `data` on success; errors use `ok`, `error`, `message`, and `details`.
