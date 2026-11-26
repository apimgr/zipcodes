# Installation

## Docker Compose

The production Compose file uses the published image and Docker-managed named volumes for configuration, database data, and logs:

```sh
docker compose -f docker/docker-compose.yml up -d
docker compose -f docker/docker-compose.yml logs -f zipcodes
```

The default host port is `172.17.0.1:64580`; change the `ports` mapping in the Compose file if that address is unavailable or the service should be exposed elsewhere. Review the bind address and firewall before exposing a public service. Stop it with `docker compose -f docker/docker-compose.yml down`; this keeps the named volumes. Removing a volume also removes its stored data.

## Binary

Download the server and, if needed, CLI binaries for the target OS and architecture from [GitHub Releases](https://github.com/apimgr/zipcodes/releases). On Linux AMD64, for example:

```sh
curl -q -LSsfO https://github.com/apimgr/zipcodes/releases/latest/download/zipcodes-linux-amd64
chmod 755 zipcodes-linux-amd64
./zipcodes-linux-amd64 --address 127.0.0.1 --port 8080
```

The server initializes its SQLite database from the embedded dataset. It creates its configuration, data, and log directories on first startup; see [Configuration](configuration.md) for directory overrides and defaults. `zipcodes --help` lists the supported options.

## System service

On a supported host, install and start the service with `sudo zipcodes --service --install`; then manage it through the host service manager or use `--service start`, `--service stop`, `--service restart`, and `--service reload`. Use `--service --disable` to stop and disable it, or `--service --uninstall` to remove it. Alternatively, create a service unit that runs the binary as a dedicated unprivileged account with writable configuration, data, and log directories. Set the listener address and port explicitly for deployment. Run `zipcodes --service --help` for platform-specific details.

## CLI

The companion `zipcodes-cli` binary requires `--server URL`, `ZIPCODES_SERVER_PRIMARY`, or `server.primary` in its per-user `cli.yml`. No hosted service URL is defined by this repository; see the [CLI reference](cli.md).
