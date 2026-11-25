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

## Business logic

### Product scope and non-goals

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

### Actors and access

There are no end-user accounts or human roles in the product. Anonymous visitors may use public search and service information. Anonymous API credentials may be used for owner-scoped client features, but they do not create accounts, profiles, or login sessions. A server operator configures and maintains an instance without a web administration panel.

ZIP lookup data is public and read-only. Operator-controlled reporting and maintenance features do not grant anonymous visitors access to server administration.

### Privacy, trust, and abuse requirements

- Treat submitted requests and uploaded contact information as untrusted.
- Do not store or expose a visitor's location as a side effect of GeoIP lookup. Missing GeoIP data must not prevent ZIP lookups.
- Keep private vulnerability report contents confidential. Public report status may expose progress, but not the submitted report body.
- Protect service availability from abusive request volume and unbounded geographic queries.
- Keep public discovery and documentation free of private configuration and credentials.
- Do not imply that an instance provides a transport, contact channel, or public overlay address unless that capability is configured and available.

ZIP records contain no personal information. The product may process a request's network address for service operation and optional approximate GeoIP lookup. Security reporting may contain sensitive researcher contact details and vulnerability information; those reports must remain confidential to maintainers.
