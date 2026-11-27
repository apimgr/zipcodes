# Development

## Tooling

Use Docker and Make. Go builds and checks run in `casjaysdev/go:latest`; the repository does not require a host Go installation. The Makefile sets `CGO_ENABLED=0` and uses Docker-mounted module and build caches.

## Build

- `make build` builds debug server and CLI binaries for the local platform and all eight supported OS/architecture combinations.
- `make build BUILD_GUI=1` additionally creates separate GUI-tagged CLI binaries for Linux, macOS, and Windows on amd64 and arm64. These optional binaries dynamically load platform graphics libraries; the default server, CLI, and TUI builds remain static and do not include GUI dependencies.
- `make release` builds stripped production binaries and release metadata in `releases/`; it does not publish a release.
- `make docker` builds the local Docker image, and `make clean` removes the ignored build-output directories.
- `make build`, `make release`, and `make test` run `go generate ./src/swagger` in the Go container. The basic Docker test runner also regenerates it before its direct server build. This regenerates `src/swagger/openapi.json` from registered API routes and annotations before compilation and checks. Do not edit that generated file directly.

## Tests and static analysis

`make test` runs `go vet`, `staticcheck`, `govulncheck`, ShellCheck, `shfmt`, Go tests with coverage, and the configured coverage threshold in the Go container, then runs the basic Docker binary checks. Run the optional browser suite directly with `./tests/e2e.sh`. Run service lifecycle tests directly with `./tests/incus.sh [systemd|openrc|all]`; `INCUS_MODE=systemd` or `INCUS_MODE=openrc` selects one service manager, and `SKIP_BUILD=1` reuses existing binaries in `BINDIR`. `./tests/run_tests.sh` runs only the basic Docker checks. The GitHub CI workflow also runs the configured build, vulnerability, and secret-scanning jobs.

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

Keep REST, GraphQL, OpenAPI, server-rendered pages, README, and `docs/` synchronized when an API or user-facing behavior changes. Add regression coverage for changed behavior. Keep build outputs out of commits; commit `src/swagger/openapi.json` with source changes because it is the generated OpenAPI superset embedded in the server. Review the complete diff before submitting a change.
