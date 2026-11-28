# Zipcodes

[![CI](https://github.com/apimgr/zipcodes/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/zipcodes/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/apimgr/zipcodes)](https://github.com/apimgr/zipcodes/releases)
[![License](https://img.shields.io/github/license/apimgr/zipcodes)](LICENSE.md)

## About

Zipcodes is a public, read-only service for searching United States ZIP codes by code, city, state, or geographic distance. It provides a browser search page, REST and GraphQL APIs, optional caller GeoIP lookup, and a companion client.

## Features

- Exact ZIP lookup, city/state search, autocomplete, statistics, and radius search.
- Public read-only REST and GraphQL APIs with interactive documentation.
- Server-rendered, mobile-friendly search pages that work without JavaScript.
- Optional caller GeoIP lookup and optional Tor and I2P service exposure.
- Companion CLI with interactive terminal mode and separate optional GUI builds.
- Six embedded interface languages: English, Spanish, French, Arabic, German, and Italian.

## Production

### Binary installation

Download the latest release from [GitHub Releases](https://github.com/apimgr/zipcodes/releases/latest).

### Linux

| Architecture | Server | CLI |
|---|---|---|
| amd64 | `zipcodes-linux-amd64` | `zipcodes-cli-linux-amd64` |
| arm64 | `zipcodes-linux-arm64` | `zipcodes-cli-linux-arm64` |

```sh
arch=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -q -LSsf "https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-linux-${arch}" -o zipcodes
chmod +x zipcodes
```

Move the binary to a directory on your `PATH` if you want to run `zipcodes` from any directory.

### macOS

| Architecture | Server | CLI |
|---|---|---|
| Intel (x86_64) | `zipcodes-darwin-amd64` | `zipcodes-cli-darwin-amd64` |
| Apple Silicon (arm64) | `zipcodes-darwin-arm64` | `zipcodes-cli-darwin-arm64` |

```sh
arch=$(uname -m | sed 's/x86_64/amd64/')
curl -q -LSsf "https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-darwin-${arch}" -o zipcodes
chmod +x zipcodes
xattr -d com.apple.quarantine zipcodes 2>/dev/null || true
```

Move the binary to a directory on your `PATH` if you want to run `zipcodes` from any directory.

### Windows

| Architecture | Server | CLI |
|---|---|---|
| amd64 | `zipcodes-windows-amd64.exe` | `zipcodes-cli-windows-amd64.exe` |
| arm64 | `zipcodes-windows-arm64.exe` | `zipcodes-cli-windows-arm64.exe` |

Download the matching executable from the release page and add its directory to `PATH`.

### FreeBSD

| Architecture | Server | CLI |
|---|---|---|
| amd64 | `zipcodes-freebsd-amd64` | `zipcodes-cli-freebsd-amd64` |
| arm64 | `zipcodes-freebsd-arm64` | `zipcodes-cli-freebsd-arm64` |

Download the matching binary, then run `chmod +x zipcodes-freebsd-*` before starting it.

### System service

On supported systems, install the service with `sudo zipcodes --service --install`. See the [installation guide](docs/installation.md#system-service) for service management and platform details.

### Docker

The production Compose file is `docker/docker-compose.yml`. It uses the published image and Docker-managed named volumes for configuration, database data, and logs.

```sh
docker compose -f docker/docker-compose.yml up -d
```

The default host port is `172.17.0.1:64580`. `docker compose down` keeps the named volumes; removing a volume deletes its stored data. See the [installation guide](docs/installation.md) before changing exposure or storage settings.

### Server setup

Run `zipcodes --help` for server options. On first start, the server creates its configuration and data files. Server configuration is operator-managed; the application has no web admin panel. See the [installation guide](docs/installation.md) to configure persistent directories.

## Client

Install the `zipcodes-cli` binary for your operating system and architecture from the same release page. Choose a server explicitly; this project does not define a hosted endpoint.

```sh
zipcodes-cli --server https://YOUR_SERVER 90210
zipcodes-cli --server https://YOUR_SERVER api GET '/api/v1/zipcodes/nearby?zip=90210&radius=10'
```

Replace `https://YOUR_SERVER` with the URL of the server you operate. This project does not provide a hosted API endpoint.

With no command in an interactive terminal, a GUI-enabled build opens its GUI when a local display is available; otherwise the client opens its TUI. First-run setup can save the connection to the selected config and, for a named profile, copy it to the default config. The GUI also offers per-user automatic login startup. The native GUI does not yet expose operating-system screen-reader accessibility APIs; use the CLI or TUI for screen-reader-accessible operation. `--theme dark|light|auto` selects the terminal theme, and `--update check` or `--update yes` checks or installs a verified client update. Optional GUI builds are separate `zipcodes-cli-gui-{os}-{arch}` release files. See the [CLI reference](docs/cli.md) for options, configuration, and token use.

## Configuration

The server reads `server.yml`, ordered `yaml.d/` area files, and the separate `secrets.yml`. Run `zipcodes --maintenance backup` to create a verified archive, `zipcodes --maintenance backup test [file]` to check one without changing live files, and `zipcodes --maintenance restore [file]` to restore it. When omitted, the restore and backup test select the full backup with the newest filename date, breaking ties by modification time. Backups include the server-owned application secrets required to restore project security keys; remote database contents remain provider-managed, with an optional local SQLite cache snapshot included when enabled. Root or an operator-token holder can rotate `installation_secret` or the `secrets.yml` at-rest `encryption_key` with `zipcodes --maintenance secret rotate`. See the [configuration reference](docs/configuration.md) for settings, environment variables, backups, and secret handling.

## API

- REST base: `/api/{api_version}` (default `v1`); discover the active version at `/api/autodiscover`.
- OpenAPI JSON: `/api/{api_version}/server/swagger`; browser reference: `/server/docs/swagger`.
- GraphQL: `/api/{api_version}/server/graphql`; browser interface: `/server/docs/graphql`.
- Read-only ZIP search pages are available at `/` and `/search`.

The project does not publish a hosted API URL. See the [API reference](docs/api.md), [security guide](docs/security.md), and [integrations guide](docs/integrations.md) for discovery, authentication, monitoring, and reporting details.

## Other

The checked-in dataset contains 33,100 geographic ZIP records from Ready APIs curated-us-zips, pinned to `64ac28b620a39d9efabe39f91dc03181a13a8a76` and licensed CC BY 4.0. It excludes APO/FPO/DPO ZIPs. The source does not provide telephone area codes, so `area_codes` is empty. See [LICENSE.md](LICENSE.md) and [the dataset attribution file](data/licenses/zipcodes-data.md) for attribution and terms.

## Development

Builds and Go checks run in Docker using `casjaysdev/go:latest`; do not run Go tools on the host. Use `make build` for the eight server and CLI platform combinations, `make test` for tests, vet, vulnerability checking, and coverage, and `make build BUILD_GUI=1` for optional GUI client artifacts. See the [development guide](docs/development.md) for the other check scripts and dataset update process. Contributors can use the [translation guide](docs/translating.md) to update one of the six supported languages.

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
