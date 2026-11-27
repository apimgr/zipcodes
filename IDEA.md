## Project description

Zipcodes is a public service for finding United States ZIP codes and exploring their geographic relationships. It serves residents, developers, and organizations that need geographic ZIP data through a browser and machine-readable interfaces. The project also provides a companion client for people who work from a terminal or desktop.

## Project variables

project_name: zipcodes
project_org: apimgr
internal_name: zipcodes
internal_org: apimgr
app_name: Zipcodes API
repo: https://github.com/apimgr/zipcodes
license: MIT
binary: zipcodes
client_binary: zipcodes-cli
owner_token: zipcodes_owner_token_qwhzkfta

## Business logic

### Product scope & non-goals

**The product must provide:**

- Geographic records for 33,100 United States ZIP codes, including city, state, county, coordinates, and IANA time zone.
- Exact ZIP lookup; search by ZIP, city, and state; autocomplete; state and city listings; random lookup; dataset statistics; and geographic radius search.
- A public, read-only REST API and GraphQL API for ZIP data, with interactive API documentation.
- A mobile-friendly browser interface for searching ZIP codes and viewing public service information. Core search and forms must work without JavaScript.
- Optional approximate caller location lookup using GeoIP data. GeoIP results are informational and must not change ZIP records.
- Public privacy, terms, security disclosure, contact, and service information pages.
- A companion client for exact lookups and requests to a server selected by the user. The client must support command-line use and interactive terminal use; separate optional builds may provide a native GUI.
- English, Spanish, French, Arabic, German, and Italian interface translations. Arabic uses a right-to-left layout.
- Public feeds and discovery documents that accurately describe the instance's available interfaces.

The checked-in dataset comes from Ready APIs curated-us-zips at commit `64ac28b620a39d9efabe39f91dc03181a13a8a76` and is licensed CC BY 4.0. It excludes APO/FPO/DPO ZIPs. The source does not contain telephone area-code assignments, so `area_codes` is empty until a suitable source is adopted.

**The product must not provide:**

- User accounts, registration, or end-user login.
- A web-based administration panel. Server administration remains an operator function.
- Paid plans, paid API keys, or differentiated access tiers.
- Non-US postal-code lookup or real-time ZIP dataset updates.
- Guessed ZIP records or area codes inferred from geographic data.
- A project-hosted public API endpoint. Users select the server to which their client connects.

### Roles & permissions

There are no end-user accounts or human roles in the product. Anonymous visitors may use public search and service information. Anonymous API credentials may be used for owner-scoped client features, but they do not create accounts, profiles, or login sessions. A server operator configures and maintains an instance without a web administration panel. ZIP lookup data and GeoIP responses are publicly readable; security report contents and operator maintenance data remain private to maintainers.

### Data model & sensitivity

- **ZIP record:** five-digit ZIP code, city, state abbreviation and name, county, latitude, longitude, IANA time zone, and area-code list. The source does not provide telephone area codes, so the list is empty. These records contain no personal information.
- **GeoIP result:** queried IP address and any available country, region, city, coordinates, and ASN information. Location is approximate and must not be retained as a side effect of lookup.
- **Contact and security report:** report category, submitted description, optional contact details, and tracking/status information. Security report contents and researcher contact details are confidential; public status contains progress only.

**Business rules:** ZIP records are read-only; lookups return only source-backed records and never infer missing ZIP or area-code data. Radius searches are bounded by the selected server's configured maximum. Missing GeoIP data does not affect ZIP lookup.

### Privacy, trust, and abuse requirements

- Treat submitted requests and uploaded contact information as untrusted.
- Do not store or expose a visitor's location as a side effect of GeoIP lookup. Missing GeoIP data must not prevent ZIP lookups.
- Keep private vulnerability report contents confidential. Public report status may expose progress, but not the submitted report body.
- Protect service availability from abusive request volume and unbounded geographic queries.
- Keep public discovery and documentation free of private configuration and credentials.
- Do not imply that an instance provides a transport, contact channel, or public overlay address unless that capability is configured and available.

ZIP records contain no personal information. The product may process a request's network address for service operation and optional approximate GeoIP lookup. Security reporting may contain sensitive researcher contact details and vulnerability information; those reports must remain confidential to maintainers.

### Trust boundaries & external services

Search terms, ZIP and coordinate parameters, request headers, contact forms, and uploaded report material are untrusted. The embedded ZIP dataset is the pinned Ready APIs curated-us-zips source. Optional GeoIP databases come from ip-location-db and DB-IP and may be unavailable or stale. Configured email, webhook, and overlay services are external and are not implied to exist unless configured and healthy.

### Threat model & abuse cases

Protect report contents and researcher contact data from public disclosure; protect service availability from high-volume requests, oversized inputs, and geographic searches with excessive radius; prevent request-controlled HTML or paths from becoming executable or exposing private files; and keep caller location from being retained or treated as an identity signal. The ZIP dataset itself is public.

### Security decisions & exceptions

No product-specific exception to AI.md. Public ZIP and GeoIP lookups are read-only; owner-scoped credentials do not create user accounts or grant access to operator administration. GeoIP is informational and must not alter ZIP data or act as the sole access-control signal.

### API routes

- `GET /search`: no authentication; render ZIP search results as HTML from optional `q`, `zip`, `radius`, `lat`, and `lng` query values; return the search page and any matching ZIP records.
- `GET /api/{api_version}/zipcodes/{code}`: no authentication; look up one five-digit ZIP; return the matching record in the standard API envelope or a validation/not-found error.
- `GET /api/{api_version}/zipcodes/{code}.txt`: no authentication; return the matching ZIP record as plain text.
- `GET /api/{api_version}/zipcodes/search?q={query}`: no authentication; search ZIP, city, state, or city/state and return matching records.
- `GET /api/{api_version}/zipcodes/autocomplete?q={query}`: no authentication; return bounded matching suggestions.
- `GET /api/{api_version}/zipcodes/city/{city}`: no authentication; list ZIP records for a city.
- `GET /api/{api_version}/zipcodes/state/{state}`: no authentication; list ZIP records for a state abbreviation.
- `GET /api/{api_version}/zipcodes/nearby?zip={zip}&radius={miles}` or `?lat={lat}&lng={lng}&radius={miles}`: no authentication; return ZIP records within the allowed radius of a ZIP or coordinate origin.
- `GET /api/{api_version}/zipcodes/stats` and `GET /api/{api_version}/zipcodes/stats.txt`: no authentication; return dataset counts as JSON or plain text.
- `GET /api/{api_version}/zipcodes/random` and `GET /api/{api_version}/zipcodes/random.txt`: no authentication; return one random record as JSON or plain text.
- `GET /api/{api_version}/zipcodes.json`: no authentication; download the complete embedded ZIP dataset as JSON.
- `GET /api/{api_version}/me`: no authentication; return the request's caller IP and available approximate GeoIP fields.
- `GET /api/{api_version}/geoip?ip={address}` and `GET /api/{api_version}/geoip.txt?ip={address}`: no authentication; look up an IP and return available location data as JSON or plain text; private/reserved IPs are not geolocated.
- `POST /api/{api_version}/geoip/batch`: no authentication; accept a JSON list of up to 100 IP addresses and return a corresponding bounded list of location results.
- `GET /api/{api_version}/server/healthz.txt`: no authentication; return the health check in plain text.
- `GET /api/{api_version}/server/preferences/tokens`: primary-token Bearer authentication; return masked primary and linked token metadata.
- `POST /api/{api_version}/server/preferences/tokens/primary`: no authentication; issue a new owner and primary token.
- `POST /api/{api_version}/server/preferences/tokens/primary/reveal` and `POST /api/{api_version}/server/preferences/tokens/primary/rotate`: primary-token Bearer authentication; reveal or rotate the primary token.
- `DELETE /api/{api_version}/server/preferences/tokens/primary`: primary-token Bearer authentication; revoke the owner and its resources.
- `POST /api/{api_version}/server/preferences/tokens`: primary-token Bearer authentication; create a linked token.
- `POST /api/{api_version}/server/preferences/tokens/{id}/reveal` and `POST /api/{api_version}/server/preferences/tokens/{id}/rotate`: primary-token Bearer authentication; reveal or rotate a linked token.
- `PATCH /api/{api_version}/server/preferences/tokens/{id}`: primary-token Bearer authentication; rename a linked token or change its expiry.
- `DELETE /api/{api_version}/server/preferences/tokens/{id}`: primary-token Bearer authentication; revoke a linked token.
- `GET /server/docs/graphql`: no authentication; render the GraphQL documentation and query form.
- `POST /server/docs/graphql`: CSRF-protected no-JavaScript query form; execute a read-only GraphQL query and render the JSON result in the page.
