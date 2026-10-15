# Installation

## Docker

Build the production image from the repository using its Dockerfile and run it with persistent storage for the configured data and configuration directories. The Compose files document development and test layouts.

## Binary

Download the release for the target OS and architecture from the GitHub Releases page. Run `zipcodes --help` for supported flags. The server initializes its SQLite database from embedded data.

## System service

Run the binary as a dedicated unprivileged account with writable user data/configuration directories and a restart policy managed by the host service manager. Set the listener address and port explicitly for deployment.

## CLI

The companion `zipcodes-cli` binary requires `--server URL` or `ZIPCODES_SERVER_PRIMARY` / `cli.yml` configuration. No hosted service URL is defined by this repository.
