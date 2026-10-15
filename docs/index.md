# Zipcodes

Zipcodes is a public, read-only service for U.S. ZIP lookup, city/state search, and bounded geographic radius queries. The server renders the browser interface and exposes versioned JSON routes, GraphQL, and a companion CLI.

## Quick start

Run the server binary, then open `/` on the configured listener. API documentation is available at `/api/v1/server/swagger`; GraphQL is at `/graphql` and `/api/v1/graphql`.

## Product data

The checked-in dataset contains 33,100 geographic ZIP records from Ready APIs curated-us-zips, pinned to commit `64ac28b620a39d9efabe39f91dc03181a13a8a76` and licensed CC BY 4.0. It excludes military APO/FPO/DPO ZIPs. Telephone area-code assignments are currently unavailable from a redistributable ZIP mapping source; `area_codes` is empty. See [data attribution and limitations](../data/licenses/zipcodes-data.md).

## Guides

- [Installation](installation.md)
- [Configuration](configuration.md)
- [API reference](api.md)
- [CLI reference](cli.md)
- [Security](security.md)
- [Integrations](integrations.md)
- [Development](development.md)

The project is MIT licensed. The ZIP dataset is separately licensed CC BY 4.0; see [LICENSE.md](../LICENSE.md) and [dataset attribution](../data/licenses/zipcodes-data.md).
