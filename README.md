# Zipcodes

[![CI](https://github.com/apimgr/zipcodes/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/zipcodes/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/apimgr/zipcodes)](LICENSE.md)
[![Documentation](https://readthedocs.org/projects/apimgr-zipcodes/badge/?version=latest)](https://apimgr-zipcodes.readthedocs.io/)

Zipcodes is a public, read-only US ZIP code lookup service. It provides a server-rendered search page, a versioned REST API, a read-only GraphQL endpoint, optional caller GeoIP lookup, and a companion CLI.

## Features

- Exact ZIP lookup, city/state search, autocomplete, and radius search.
- Public read-only REST and GraphQL APIs.
- Server-rendered search interface and installable PWA shell.
- Optional GeoIP lookup for the caller or supplied IP addresses.
- Companion CLI and embedded, indexed ZIP data.
- Six embedded locale catalogs: English (`en`), Spanish (`es`), French (`fr`), Arabic (`ar`), German (`de`), and Italian (`it`). See the [translation guide](docs/translating.md).

The checked-in dataset contains 33,100 geographic ZIP records from Ready APIs curated-us-zips, pinned to `64ac28b620a39d9efabe39f91dc03181a13a8a76` and licensed CC BY 4.0. It excludes APO/FPO/DPO ZIPs. The source does not provide telephone area codes, so `area_codes` is empty. See [LICENSE.md](LICENSE.md) for attribution.

## Production

Download a release for your operating system and architecture from [GitHub Releases](https://github.com/apimgr/zipcodes/releases), then run the server binary. It creates its configuration and data files on first start. The default HTTP port is selected at startup; use `--help` for supported options.

The production Compose file is `docker/docker-compose.yml`; its host port is bound to `172.17.0.1:64580`. For native installs, use the release binary and run it under a service manager with dedicated writable configuration, data, and log directories. Persistent data should be stored outside the repository checkout. See [installation](docs/installation.md) and [configuration](docs/configuration.md).

## Client

The `zipcodes-cli` client performs exact ZIP lookups against an operator-selected server:

```sh
zipcodes-cli --server http://localhost:8080 90210
```

With no command in an interactive terminal, `zipcodes-cli` opens its lookup interface automatically. In a non-interactive shell, pass one five-digit ZIP code. Display mode is configured in `cli.yml`; `--theme dark|light|auto` selects a theme for the current run. Private deployments may supply `--token` to validate and use an owner or operator Bearer token; the client stores it only in the mode-`0600` per-user config file.

Set `ZIPCODES_SERVER_PRIMARY` or `server.primary` in `cli.yml` to save the server choice. This project does not define a hosted endpoint.

## Search

Open the server root page to search by ZIP code, city, or state. Radius search accepts an origin ZIP or coordinates and is bounded by the server's configured maximum (100 miles by default).

```sh
curl -H 'Accept: application/json' http://localhost:8080/api/v1/zipcodes/90210
curl -H 'Accept: application/json' 'http://localhost:8080/api/v1/zipcodes/nearby?zip=90210&radius=10'
zipcodes-cli --server http://localhost:8080 90210
```

The CLI requires an operator-selected server URL. Configure it with `--server`, `ZIPCODES_SERVER_PRIMARY`, or `server.primary` in `cli.yml`. This project does not define a hosted service URL.

## Configuration

The server creates `server.yml` with defaults on first start. Command-line options and environment variables can override settings; there is no admin web panel or runtime settings API. See [configuration reference](docs/configuration.md).

## API and pages

- REST API: `/api/{api_version}` (default `v1`, discoverable through `/api/autodiscover`; the trailing-slash form redirects with `301`; see [API reference](docs/api.md)).
- OpenAPI JSON and UI: `/api/v1/server/swagger` and `/server/docs/swagger`.
- GraphQL endpoint and browser interface: `/api/v1/server/graphql` and `/server/docs/graphql`.
- Remote MCP endpoint and connection guide: `/api/v1/server/mcp` and `/server/help/mcp` (owner/operator Bearer token required).
- Internal monitoring endpoints: Prometheus at `/server/metrics`, Grafana dashboard JSON at `/server/metrics/grafana`, and recent Loki-compatible HTTP logs at `/server/metrics/loki`; each uses its own generated service token from `secrets.yml`.
- Server-rendered pages: `/`, `/search`, `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/security`, `/server/security/policy`, `/server/security/thanks`, `/server/security/report/{tracking_id}`, `/server/contact`, `/server/dpo`, `/server/help/mcp`, `/server/preferences`, `/server/preferences/tokens`, and `/server/preferences/sync`. The optional root `/healthz` alias is disabled by default; see [configuration](docs/configuration.md).
- Security contact: `/.well-known/security.txt`.
- Discovery endpoints: `/.well-known/llms.txt` and `/llms.txt`.
- Empty Atom and RSS feeds: `/feed.xml` and `/rss.xml` (ZIP records have no per-record permalink or publication timestamp).

All ZIP lookup endpoints are public and read-only. GeoIP lookup is informational for callers; the optional operator-configured country deny/allow rules can reject requests after rate limiting, while always failing open when country data is unavailable. ZIP dataset reads are unchanged.

## Documentation

The [documentation index](docs/index.md) links installation, configuration, API, CLI, security, integration, development, and translation guides. Documentation is built with MkDocs Material and hosted through ReadTheDocs.

## Development

Follow the Docker-based build and test targets in the Makefile. The dataset updater reads its pinned upstream source; do not invent records or derive area codes from geography. See [development](docs/development.md).

## Disclaimer

This software is provided "as is" without warranty of any kind. Use it at your own risk.

- **No warranty:** The authors are not responsible for damages, data loss, or other issues arising from use of this software.
- **No professional advice:** This software does not provide legal, financial, medical, or other professional advice.
- **Third-party services:** External APIs and services are subject to their own terms.
- **Security:** The project follows security practices, but no software can be guaranteed free of vulnerabilities.
- **Production use:** Evaluate the software and its configuration before deploying it in production.

By using this software, you acknowledge that you have read and understood this disclaimer.

## License

The application is MIT licensed. The embedded ZIP dataset is separately licensed CC BY 4.0. See [LICENSE.md](LICENSE.md) for project and third-party attribution.
