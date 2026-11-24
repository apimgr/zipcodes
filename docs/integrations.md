# Integrations

Zipcodes exposes public read-only interfaces for browser users and API clients:

- REST API under `/api/v1` with an OpenAPI document at `/api/v1/server/swagger` and browser UI at `/server/docs/swagger`.
- Read-only GraphQL at `/api/v1/server/graphql`, with a browser interface at `/server/docs/graphql` and `/graphql`.
- The PWA manifest at `/manifest.json` and service worker at `/sw.js`.
- `/robots.txt` and the generated `/sitemap.xml` for public page discovery.
- Empty Atom (`/feed.xml`) and RSS (`/rss.xml`) feeds because ZIP records have no public per-record page or publication timestamp.
- `/.well-known/llms.txt` with request-derived absolute links to the API and its documentation.
- A companion `zipcodes-cli` that connects to an operator-selected server.
- `/security.txt` and `/.well-known/security.txt` for vulnerability reporting information.

API clients can optionally use anonymous owner credentials to associate activity with an owner. A primary owner token manages linked tokens; linked tokens are for automation and cannot manage tokens. This is credential-based identity only: there are no user accounts, login sessions, or federated identity flows. See [API reference](api.md) and [Security](security.md) for token routes and handling.

The project does not provide webhooks, identity federation, user accounts, or a hosted API address. The CLI server URL is supplied with `--server`, `ZIPCODES_SERVER_PRIMARY`, or `server.primary` in its `cli.yml`.

GeoIP database downloads are an optional data integration used for caller/IP lookup. They are fetched by the server over HTTPS from the configured ip-location-db source; an unavailable database leaves ZIP lookup functionality available.
