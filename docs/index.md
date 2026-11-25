# Zipcodes

Zipcodes is a public, read-only service for U.S. ZIP lookup, city/state search, and bounded geographic radius queries. It includes server-rendered search pages, a versioned REST API, GraphQL, OpenAPI documentation, optional caller GeoIP, and a companion CLI.

## Quick start

With Docker Compose:

```sh
docker compose -f docker/docker-compose.yml up -d
```

Open `http://172.17.0.1:64580/` on the Docker host. For a native Linux AMD64 run, download `zipcodes-linux-amd64`, then start it with `./zipcodes-linux-amd64 --address 127.0.0.1 --port 8080` and open `http://127.0.0.1:8080/`. The search page works without JavaScript. See the [installation guide](installation.md) for persistent storage and service setup.

The companion CLI connects to an explicitly configured server:

```sh
zipcodes-cli --server http://127.0.0.1:8080 90210
```

Run `zipcodes-cli` without arguments in an interactive terminal for repeated lookups. In scripts or non-interactive shells, provide one five-digit ZIP code.

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

## Links

- [Repository](https://github.com/apimgr/zipcodes)
- On a running instance: `/server/docs/swagger` (OpenAPI) and `/server/docs/graphql` (GraphQL)
