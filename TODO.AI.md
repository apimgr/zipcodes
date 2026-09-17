# TODO.AI.md

## CI/CD follow-ups (from initial workflow bootstrap)

- [ ] **Add `docker/Dockerfile.dev` and `docker-compose.dev.yml`.** AI.md's CI/CD Rules require a `build-devel` job in `docker.yml` that builds the `:devel` tag from `docker/Dockerfile.dev` (daily schedule + non-tag push + `workflow_dispatch`). Sibling repos (`api`, `ipgaze`) have this file; zipcodes does not, so `docker.yml` currently ships only the `build-standard` job. Add the dev Dockerfile/compose file, then add the `build-devel` job to `.github/workflows/docker.yml`.
- [ ] **Write Go tests — zero coverage today.** `go test -cover ./...` (run in `casjaysdev/go:latest`) shows 0.0% coverage across every package (`admin`, `api`, `config`, `data`, `database`, `geoip`, `mode`, `paths`, `scheduler`, `server`, `service`, `ssl`, `utils`) — `tests/` holds only `.gitkeep`, no `*_test.go` files exist anywhere. `ci.yml`'s `test` job enforces a 60% coverage threshold (project default, no override in IDEA.md), so CI's `test`/`coverage` jobs will fail until real tests are added.
- [ ] **Add `tests/run_tests.sh`.** AI.md's Container-Only Development section references a Phase 2 binary-validation script at this path; it does not exist (`tests/` has only `.gitkeep`).
