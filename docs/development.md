# Development

## Tooling

Use Docker and Make. Go builds and checks run in `casjaysdev/go:latest`; the repository does not require a host Go installation. The Makefile sets `CGO_ENABLED=0` and uses Docker-mounted module and build caches.

## Build

- `make dev` builds development binaries into a project-scoped temporary directory.
- `make local` builds the server and CLI for the host platform under the gitignored `binaries/` directory.
- `make build` cross-compiles the server and CLI for the supported OS/architecture matrix.
- `make build BUILD_GUI=1` additionally creates separate GUI-tagged CLI binaries for Linux, macOS, and Windows. These optional binaries dynamically load platform graphics libraries; the default server, CLI, and TUI builds remain static and do not include GUI dependencies.

## Tests and static analysis

`make test` runs `go vet`, `govulncheck`, Go tests with coverage, and the configured coverage threshold in the Go container. Use `make test E2E=1` for the Docker browser suite or `make test INCUS=1` for the systemd and OpenRC lifecycle matrix. These modes use the existing `tests/e2e.sh` and `tests/incus.sh` runners. `./tests/run_tests.sh` runs the Docker-based basic checks. The GitHub CI workflow also runs `staticcheck`, build, vulnerability, and secret scanning.

## Documentation

Build the MkDocs site with strict link checking before submitting documentation changes:

```sh
mkdocs build --strict
```

Install the Python packages from `docs/requirements.txt` in an isolated Python environment before running MkDocs locally.

## Dataset

The dataset updater is `scripts/update_zipcodes_data.py`. It reads the pinned Ready APIs curated-us-zips source revision and writes `src/data/zipcodes.json`. The data is licensed CC BY 4.0; source details, attribution, and exclusions are in `data/licenses/zipcodes-data.md`.

Do not add guessed ZIP records or infer area codes from coordinates, states, or counties. Update source attribution whenever the pinned dataset changes.

## Contributing

Keep REST, GraphQL, OpenAPI, server-rendered pages, README, and `docs/` synchronized when an API or user-facing behavior changes. Add regression coverage for changed behavior, keep generated/build files out of commits, and review the complete diff before submitting a change.
