# AI Implementation Tasks

These tasks track verified work against `AI.md`. They remain open until the implementation, documentation, and required verification are complete.

## Security and contact workflows

- [ ] Implement project OpenPGP key generation, encrypted private-key storage, export/import/delete, publication, rotation, and crash recovery from PART 11.
- [ ] Complete security-report maintenance list/show/export/update workflows, audit events, status notifications, and retention behavior from PART 11.
- [ ] Complete contact workflow requirements from PARTS 11, 12, and 16, including provider CAPTCHA/no-JavaScript paths, edit/retry behavior, PGP fallback, and exact field/API errors.
- [ ] Complete project privacy, preference-language, and preference-sync behavior after resolving the conflicting `AI.md` directives with the coordinator; cookie-consent persistence now uses cookie-safe encoded JSON and its consent gate is covered by browser E2E.

## Runtime and integration behavior

- [ ] Complete PART 4/8 startup path resolution, directory permissions and ownership, `--daemon`/`server.daemonize`, and PID lifecycle sequencing; platform/container path helpers and basic PID-file handling are implemented, but startup ordering, privilege transitions, and daemonization remain incomplete.
- [ ] Complete owner and linked token lifecycle, preferences pages, and sync contract from PARTS 11 and 16.
- [ ] Complete MCP route/tool behavior and annotations from PART 11, including the sync-import method conflict recorded for coordinator review.
- [ ] Complete remote database support, backup rotation/incremental behavior, and secrets rotation from PARTS 10, 11, and 21.
- [ ] Complete service installation, update, and cross-platform binary behavior from PARTS 23–25.
- [ ] Complete TLS, Tor, I2P, and overlay lifecycle behavior from PARTS 30–31.

## Verification and documentation

- [ ] Compare every applicable `AI.md` PART against source, SSR pages, tests, CI, and documentation; implement each verified gap.
- [ ] Refactor `src/main.go` to the thin flag/config entry and move complete startup, signal handling, and graceful shutdown into a cancellable, unit-testable `run(ctx, cfg)` path as required by PART 28. One-listener HTTP serving and context-driven shutdown now live in `src/app.Run` with a request/readiness test, but the main startup, signal, reload, and overlay orchestration still remain in `main()`.
- [ ] Complete the Incus lifecycle scenarios required by PART 28; browser E2E scenarios currently pass.
- [ ] Bring aggregate test coverage to the `Makefile` threshold and pass all prescribed tests and lint gates. Latest wrapped `make test`: package tests, `go vet`, and `govulncheck` pass; aggregate coverage is 60.0%. `staticcheck ./...`, `golangci-lint run ./...`, ShellCheck, and strict MkDocs pass in the current checkout.
- [ ] Synchronize README, MkDocs pages, config/API/CLI docs, and `IDEA.md` with the implemented behavior.
- [ ] Resolve remaining verified AI-internal precedence issues through the coordinator without modifying `AI.md`.
