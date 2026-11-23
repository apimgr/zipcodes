# Project SPEC

Project: ZIPCODES
Role: Efficient loader for AI.md

⚠️ **THIS FILE IS AUTO-LOADED EVERY CONVERSATION. FOLLOW IT EXACTLY.** ⚠️

Purpose:
- `AI.md` is the source of truth for implementation requirements.
- `IDEA.md` defines product behavior; update it when product behavior changes.
- Read the relevant complete AI.md PART before changing code or documentation.
- Never modify AI.md. Resolve disagreements by following AI.md and updating IDEA.md.

## Project Files

- Keep source and project documentation in this repository.
- Keep runtime data, configs, logs, test output, and build artifacts outside the source tree, except gitignored `binaries/` and `releases/` targets defined by the Makefile.
- Preserve existing user changes and read files before editing them.
- Do not create report-only files. Use TODO.AI.md for three or more implementation tasks; TODO.md is human-owned.

## Project Rules

- Go code is built, tested, and run in Docker using `casjaysdev/go:latest`; never run Go tools or project binaries directly on the host.
- Keep `CGO_ENABLED=0`; build the server and client for all eight OS/architecture combinations.
- Keep `server.yml` runtime-generated. Server configuration has no web admin panel.
- Public ZIP lookup routes are read-only. Keep REST, GraphQL, browser pages, OpenAPI, README, and docs synchronized.
- Browser pages are server-rendered, mobile-first, accessible, and work without JavaScript for core flows.
- Run required tests and lint through the session wrapper, use `.git/COMMIT_MESS`, and commit through `gitcommit --dir {project_dir} all`.

## References

- Product scope and routes: `IDEA.md`.
- Implementation requirements: `AI.md PARTS 0-32`.
- Full source-to-spec check: AI.md `FINAL CHECKPOINT`.
