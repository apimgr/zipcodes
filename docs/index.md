# Zipcodes

Zipcodes is a public, read-only service for U.S. ZIP lookup, city/state search, and bounded geographic radius queries. The server renders the browser interface and exposes versioned JSON routes, GraphQL, and a companion CLI.

## Quick start

Run the `zipcodes` server binary and open `/` on its listener. The page supports server-rendered search without requiring JavaScript. The companion CLI connects to an explicitly configured server.

## Product data

The checked-in dataset contains 33,100 geographic ZIP records from Ready APIs curated-us-zips, pinned to commit `64ac28b620a39d9efabe39f91dc03181a13a8a76` and licensed CC BY 4.0. It excludes APO/FPO/DPO ZIPs. Telephone area-code assignments are unavailable from the source, so `area_codes` is empty. See the [project license](https://github.com/apimgr/zipcodes/blob/main/LICENSE.md).

## Guides

- [Installation](installation.md)
- [Configuration](configuration.md)
- [API reference](api.md)
- [CLI reference](cli.md)
- [Security](security.md)
- [Integrations](integrations.md)
- [Development](development.md)
- [Translating](translating.md)

The project is MIT licensed. The ZIP dataset is separately licensed CC BY 4.0; see the [project license](https://github.com/apimgr/zipcodes/blob/main/LICENSE.md) and [dataset attribution](https://github.com/apimgr/zipcodes/blob/main/data/licenses/zipcodes-data.md).
