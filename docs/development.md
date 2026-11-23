# Development

Use the project Docker images and Compose files for development; Go commands run in `casjaysdev/go:latest` containers. `make dev` builds into a temporary directory, `make local` writes server and CLI binaries to `binaries/`, `make test` runs the Go unit-test and coverage gate, and `make build` cross-compiles the release matrix. These targets use Docker-mounted Go module and build caches.

Run `./tests/run_tests.sh` for the basic binary and HTTP checks. It builds both binaries and runs them in a disposable Alpine container with runtime data under a project-scoped temporary directory. The Incus lifecycle and browser E2E scripts required by the full project specification are not yet present.

Before a commit, run the repository test and lint gates, build the MkDocs site, inspect the complete diff, and document every changed file in `.git/COMMIT_MESS`. Use `gitcommit --dir <repository> all` for the repository commit workflow.

The ZIP dataset is reproduced with `python3 scripts/update_zipcodes_data.py`, which fetches the pinned source revision and writes `src/data/zipcodes.json`. The dataset is CC BY 4.0 and its attribution and exclusions are documented in `data/licenses/zipcodes-data.md`.

Do not add guessed ZIP records or infer area codes from coordinates, states, or counties. Update source attribution whenever the pinned dataset changes.
