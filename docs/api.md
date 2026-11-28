# API reference

ZIP lookup routes are public and read-only. The default REST API base path is `/api/{api_version}` (default `/api/v1`); `GET /api/{api_version}` returns the API name, version, base URL, and endpoint links. The trailing-slash form redirects `301` to the slash-free API root. Operators can configure `server.api_version`, and clients discover the active version from `/api/autodiscover`. Separate preference and owner-token routes manage visitor settings and anonymous owner credentials; their access requirements are described below. Interactive OpenAPI documentation is served at `/server/docs/swagger`; the OpenAPI JSON is `/api/{api_version}/server/swagger` (also `/api/swagger`). The GraphQL interface is `/server/docs/graphql`; submit queries to `/api/{api_version}/server/graphql` or `/api/graphql`. The old root `/graphql` path is not served.

The OpenAPI document is generated from the registered routes and route annotations by `go generate ./src/swagger`, embedded in the server binary, and filtered at request time for disabled optional routes. Build and test workflows regenerate it; the route drift test compares the embedded document with the router and the AI.md/IDEA.md route inventories.

## Autodiscovery

`GET /api/autodiscover` returns a standard `{ok, data}` response. `data` contains `api.current_version`, the project `name` and `version`, request-derived public `urls`, `cli_versions`, and `cli_min_version`. The server publishes CLI versions only for operator-provisioned release files in `{data_dir}/cli` (or `cli.binary_download.dir`) whose entries in `sha256.txt` match the recomputed file hashes; `version.txt` supplies each published version. `GET /api/{api_version}/cli/binaries/{filename}` streams a verified release binary and records `cli.update_started`; downloads are public unless `cli.binary_download.require_auth` is enabled. The unversioned compatibility route is `GET /cli/binaries/{filename}`. `cli_min_version` is empty unless configured. Discovery adds `data.onion` or `data.i2p` only when the corresponding overlay is healthy and its `server.tor.expose` or `server.i2p.expose` setting is `true`; both settings default to `false`. Overlay addresses are separate from the clearnet `urls` list.

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

The browser-facing pages include `/`, `/search`, `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/security`, `/server/security/policy`, `/server/security/thanks`, `/server/security/report/{tracking_id}`, `/server/contact`, `/server/dpo`, `/server/help/mcp`, `/server/preferences`, `/server/preferences/tokens`, and `/server/preferences/sync`. Interactive API pages are `/server/docs/swagger` and `/server/docs/graphql`. Corresponding versioned API routes are listed in the OpenAPI document; machine-only endpoints such as MCP and token check have no paired page requirement.

The server also publishes `/robots.txt`, `/sitemap.xml`, `/favicon.ico`, `/.well-known/security.txt`, and `/.well-known/llms.txt`. The sitemap lists public HTML pages only; it does not list API endpoints. `/.well-known/pgp-key.asc` is available only when a project keypair exists and `web.security.publish_pgp_key` is enabled.

`GET /api/{api_version}/server/auth/check` is a machine-only endpoint. Send an owner or operator token in `Authorization: Bearer <token>`; a valid token returns `{ "ok": true, "data": { "kind": "owner" | "operator" } }`. Missing, invalid, or expired tokens return `401`. The endpoint never returns the token and does not update owner-token activity state.

`/feed.xml` (Atom 1.0) and `/rss.xml` (RSS 2.0) are always available. ZIP records do not have individual public pages or creation times, so both feeds contain no entries. They use request-derived absolute URLs and support ETag conditional requests. `/.well-known/llms.txt` and `/llms.txt` list the public API and documentation endpoints.

HTML pages support server-rendered search and require no JavaScript. Theme changes use a CSRF-protected same-origin form. Preferences form, sync, and API writes set theme, language, and application preference cookies only when the visitor has accepted the Preferences cookie category. Language resolution checks a supported `?lang=` value, then the language cookie, then the best supported `Accept-Language` match, and finally English. A saved language is managed from the preferences page and follows the same consent rule as other preference cookies.

GraphQL is read-only. The root query fields are `zipcode(code: String!)`, `search(city: String, state: String)`, and `stats`. `zipcode` returns one record; `search` requires a city, a state, or both; `stats` returns total ZIP codes, cities, and states. ZIP records expose `zip`, `zipcode`, `zipCode`, `city`, `state`, `stateFull`, `county`, `timezone`, `areaCodes`, `latitude`, `longitude`, and nested `coordinates { latitude longitude }`. The API accepts a GraphQL query by `GET` query parameter or `POST` with a JSON body or `application/graphql` body. For example:

```graphql
query GetZipcode {
  zipcode(code: "90210") {
    zip
    city
    state
    county
    timezone
  }
}
```

## MCP

The remote Streamable HTTP endpoint is `/api/{api_version}/server/mcp`; `/api/mcp` and `/mcp` are direct aliases. It accepts authenticated JSON-RPC 2.0 `POST` requests using the same owner or operator Bearer tokens as the REST API. It targets protocol revision `2026-07-28`, has no initialize handshake or session, and requires `MCP-Protocol-Version`, `Mcp-Method`, and `Accept: application/json, text/event-stream`. Tool names, descriptions, and input schemas are generated from the OpenAPI operations. Token-management and issuance operations, the cookie-bound `/server/preferences/**` family, and the MCP endpoint itself are excluded from the tool list. On a running instance, `/server/help/mcp` provides the client setup guide.

## Contact and security reporting

`GET /api/{api_version}/server/contact` describes the available report types, form fields, fill token, and CAPTCHA challenge. Its simple arithmetic question follows the visitor's selected language; when reCAPTCHA or hCaptcha is configured, the descriptor includes this simple challenge as a JavaScript-free fallback. `POST` uses a submit, review, and confirm flow; it does not deliver on the first request. Abuse, DMCA, and privacy submissions use SMTP or a configured role webhook; the request fails when no transport accepts delivery. Security reports use a rotating `security_id`, are encrypted with the project OpenPGP key when present and otherwise with the server AES-256-GCM key before database storage, and still succeed when SMTP is unavailable. Security webhooks receive an event summary only; SMTP delivery attaches the encrypted report. The response returns a `tracking_id`; the status URL requires a separate one-use-per-day token sent only in the acknowledgment email. The security key fingerprint and public download URL appear in the `/api/{api_version}/server/security` mirror when publication is enabled; `/.well-known/security.txt` is the only security.txt route.

The optional `cve_requested` security-report checkbox is retained in the encrypted report body and as report metadata. When requested, maintainers initiate CVE assignment after triage.

## Preferences

The browser preferences page is `/server/preferences`. Its single form supports appearance, language, cookie consent, and application preferences. The Language section lists supported languages by native name; Automatic returns to `Accept-Language` detection. Preference cookies (`theme`, `lang`, and application preference cookies) are written only after the visitor accepts the Preferences cookie category. Without that consent, valid changes are ignored and no preference cookies are set. The header theme toggle is subject to the same gate. `GET /api/{api_version}/server/preferences` returns the current `theme`, `lang`, and `prefs` values. `POST /api/{api_version}/server/preferences` accepts JSON or form values and applies valid preferences only when consent is present; otherwise `data.consent_required` is `true` and `applied` is empty.

Preference export/import is available at `/server/preferences/sync` and `/server/preferences/sync/import`. The sync page exports the allowed preferences as a URL and a base64url code. Import is additive and accepts the shared URL, a query string, or a short code; it applies only valid values and remains subject to the destination device's Preferences consent. The API equivalents are `GET /api/{api_version}/server/preferences/sync` and `GET /api/{api_version}/server/preferences/sync/import`. The import route can set preference cookies, so the entire cookie-bound `/server/preferences/**` route family is excluded from MCP.

## Owner tokens

Owner tokens identify an anonymous owner; they do not create user accounts. A primary token has `global` scope and can manage linked tokens. Linked tokens default to `read-write` or can be created with `read` scope; read scope permits GET, HEAD, and OPTIONS resource methods, and the MCP tool list contains only GET operations for read tokens. A read token calling an MCP write tool receives `403` with `WWW-Authenticate: Bearer error="insufficient_scope"`. Linked tokens cannot manage tokens. API token routes accept a primary owner token in `Authorization: Bearer …`; cookies are not used for these API routes. Missing credentials return `401`; valid linked credentials on a management route return `403` with the insufficient-scope challenge. The web token page uses the owner cookie and CSRF-protected POST forms, renders without JavaScript, and localizes its controls in English, Spanish, French, German, Arabic, and Italian.

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
| POST | `/server/preferences/tokens/{id}/reveal` | Reveal a linked token |
| POST | `/server/preferences/tokens/{id}/rotate` | Rotate a linked token |
| POST | `/server/preferences/tokens/{id}/update` | Rename a linked token or change its expiry |
| POST | `/server/preferences/tokens/{id}/revoke` | Revoke a linked token |

| Method | Route | Access | Purpose |
|---|---|---|---|
| GET | `/server/preferences/tokens` | Owner cookie for management; otherwise shows create/restore forms | Token management page |
| POST | `/api/{api_version}/server/preferences/tokens/primary` | No token required | Create an owner and issue its primary token; returns `201`, or `409` if a valid token is already supplied |
| GET | `/api/{api_version}/server/preferences/tokens` | Primary owner Bearer token | List tokens with values masked/omitted |
| POST | `/api/{api_version}/server/preferences/tokens/primary/reveal` | Primary owner Bearer token | Reveal the primary token |
| POST | `/api/{api_version}/server/preferences/tokens/primary/rotate` | Primary owner Bearer token; JSON `confirm` must match its last four characters | Replace the primary token |
| DELETE | `/api/{api_version}/server/preferences/tokens/primary` | Primary owner Bearer token; JSON `confirm` must match its last four characters | `200`; delete the owner and all its tokens/resources |
| POST | `/api/{api_version}/server/preferences/tokens` | Primary owner Bearer token; JSON `name`, optional `scope` (`read` or `read-write`, default `read-write`), and `expires` or `expires_on` | Create a linked token; `global` is rejected |
| POST | `/api/{api_version}/server/preferences/tokens/{id}/reveal` | Primary owner Bearer token | Reveal a linked token |
| POST | `/api/{api_version}/server/preferences/tokens/{id}/rotate` | Primary owner Bearer token; JSON `confirm` must match the token's last four characters | Replace a linked token |
| PATCH | `/api/{api_version}/server/preferences/tokens/{id}` | Primary owner Bearer token; JSON `name`, `expires` or `expires_on` (use at most one expiry field) | Rename a linked token or change its expiry |
| DELETE | `/api/{api_version}/server/preferences/tokens/{id}` | Primary owner Bearer token | `200`; revoke a linked token |

Linked-token expiry accepts `never`, `7d`, `30d`, `90d`, `1y`, or a future `YYYY-MM-DD` date through `expires_on`. The default limit is 25 active linked tokens per owner. Token objects include their scope; the primary reports `global`, while linked tokens report `read` or `read-write`. Scope cannot be changed after creation. Raw token values are returned only by issue, reveal, and rotate actions; save them securely. The primary token may be restored in a browser from `/server/preferences/tokens`. Revoking the primary deletes the owner and its associated data; rotating it preserves the owner and linked tokens.

`/server/healthz` is the canonical health frontend: browsers receive HTML, HTTP tools receive human-readable text, and `Accept: application/json` receives JSON. `/api/{api_version}/server/healthz` and `/api/healthz` are machine endpoints; they return JSON by default and support text with `.txt` or `Accept: text/plain`. The root `/healthz` alias is disabled by default and can be enabled with `server.healthz.root.enabled: true`; it uses the same content negotiation as `/server/healthz`. Health responses are bare and are not wrapped in the standard API envelope. Health checks are cached for at most 60 seconds. A live configuration change that requires a restart adds `pending_restart` and `restart_reason` to health output and changes the status to `restart_required` while the service remains available. Request totals are persisted in batches of at most 60 seconds, requests are grouped into rolling hourly buckets, and health and metrics scrapes are excluded from request totals. GeoIP checks report enabled datasets individually; unavailable non-critical collectors make the overall status `degraded`. The OpenAPI and GraphQL endpoints are described above. Other JSON API responses use `ok` and `data` on success; errors use `ok`, `error`, `message`, and `details`.
