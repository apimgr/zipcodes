# AI Implementation Tasks

These tasks track verified work against `AI.md`. They remain open until the implementation, documentation, and required verification are complete.

## Security and contact workflows

- [ ] Complete project OpenPGP key generation, encrypted private-key storage, export/import/delete, publication, rotation, and crash recovery from PART 11. Rotation now stores the new key fingerprint and metadata in the same DB transaction as report re-encryption, and recovery uses that marker even when no reports exist; still verify remaining key lifecycle, signing, and export contracts.
- [x] Complete security-report maintenance list/show/export/update workflows, audit events, status notifications, and retention behavior from PART 11. Verified against `src/security_commands.go`, `src/database/security_reports.go`, scheduled cleanup in `src/main.go`, command tests, and `docs/security.md`/`docs/cli.md`.
- [ ] Complete contact workflow requirements from PARTS 11, 12, and 16, including provider CAPTCHA/no-JavaScript paths, edit/retry behavior, PGP fallback, and exact field/API errors.
- [ ] Complete project privacy, preference-language, and preference-sync behavior after resolving the remaining conflicting `AI.md` directives with the coordinator. The preferences form now lists the six registered native language names, saves a selected language only after Preferences consent, and lets Automatic delete the override cookie. The `?lang=` persistence directive still conflicts with the preference-cookie consent gate and its cookie-visibility rules; cookie-consent browser E2E had passed before recent changes.

## Runtime and integration behavior

- [ ] Complete PART 4/8 startup path resolution, directory permissions and ownership, `--daemon`/`server.daemonize`, and PID lifecycle sequencing; platform/container path helpers and basic PID-file handling are implemented, but startup ordering, privilege transitions, and daemonization remain incomplete.
- [ ] Complete owner and linked token lifecycle, preferences pages, and sync contract from PARTS 11 and 16.
- [ ] Complete MCP route/tool behavior and annotations from PART 11. The tool generator now excludes the entire cookie-bound `/server/preferences/**` family, including GET sync import, so the state-changing import is not advertised as read-only.
- [ ] Complete remote database support, backup rotation/incremental behavior, and secrets rotation from PARTS 10, 11, and 21.
- [x] Add the versioned API root index and canonical trailing-slash redirect required by PARTS 1 and 14. The root returns request-aware links and is included in OpenAPI; server tests cover the response and `301` location.
- [ ] Complete service installation, update, and cross-platform binary behavior from PARTS 23–25.
- [ ] Complete TLS, Tor, I2P, and overlay lifecycle behavior from PARTS 30–31.

## Verification and documentation

- [x] Move server linker metadata and build-epoch parsing into `src/version.go` as required by PART 25; Makefile `-X main.*` symbols remain unchanged.
- [ ] Compare every applicable `AI.md` PART against source, SSR pages, tests, CI, and documentation; implement each verified gap.
- [ ] Refactor `src/main.go` to the thin flag/config entry and move complete startup, signal handling, and graceful shutdown into a cancellable, unit-testable `run(ctx, cfg)` path as required by PART 28. One-listener HTTP serving and context-driven shutdown now live in `src/app.Run` with a request/readiness test, but the main startup, signal, reload, and overlay orchestration still remain in `main()`.
- [ ] Complete the Incus lifecycle scenarios required by PART 28. The browser E2E suite passed before the latest client/auth-check changes, but has not been rerun against the current tree.
- [ ] Complete the PART 32 client contract: command tree for project functionality, Cobra/Viper configuration, and matching GUI/TUI views for every server-facing command. Automatic display selection, `auth.token` / `ZIPCODES_TOKEN` precedence, Bubble Tea setup and lookup screens, `/server/auth/check`, shared terminal palettes, per-invocation autodiscovery, SHA-256 validation, and staged CLI replacement/re-exec are implemented; the CLI still exposes only exact lookup and has no GUI screens.
- [ ] Bring aggregate test coverage to the `Makefile` threshold and pass all prescribed tests and lint gates. Latest wrapped `make test` passes package tests, `go vet`, and `govulncheck` at 60.3% aggregate coverage. `tests/run_tests.sh` passes the Docker binary and HTTP checks; staticcheck, golangci-lint, ShellCheck, strict MkDocs, and `git diff --check` also pass. The separate E2E and Incus scripts remain unrun because the test gate rejected their standalone script commands before execution; the previous Incus run reported lifecycle/overlay failures.
- [ ] Synchronize README, MkDocs pages, config/API/CLI docs, and `IDEA.md` with the implemented behavior.
- [ ] Resolve remaining verified AI-internal precedence issues through the coordinator without modifying `AI.md`.
