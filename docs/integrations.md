# Integrations

Zipcodes exposes public read-only interfaces for browser users and API clients:

- REST API under `/api/v1` with an OpenAPI document at `/api/v1/server/swagger` and browser UI at `/server/docs/swagger`.
- Read-only GraphQL at `/api/v1/server/graphql`, with a browser interface at `/server/docs/graphql` and `/graphql`.
- The PWA manifest at `/manifest.json` and service worker at `/sw.js`.
- A companion `zipcodes-cli` that connects to an operator-selected server.
- `/.well-known/security.txt` for vulnerability reporting information.

The project does not provide webhooks, identity federation, user accounts, or a hosted API address. The CLI server URL is supplied with `--server`, `ZIPCODES_SERVER_PRIMARY`, or `server.primary` in its `cli.yml`.

GeoIP database downloads are an optional data integration used for caller/IP lookup. They are fetched by the server over HTTPS from the configured ip-location-db source; an unavailable database leaves ZIP lookup functionality available.
