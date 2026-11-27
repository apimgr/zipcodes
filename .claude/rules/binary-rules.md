# Binary Rules (PART 7, 8, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never run a project binary directly on the host.
- Never enable CGO or build only the host platform for release artifacts.
- Never hardcode a hosted API URL or expose server configuration through the API.

## CRITICAL - ALWAYS DO
- Build static server and client binaries with CGO_ENABLED=0 for all eight OS/architecture pairs.
- Derive displayed binary names from the invoked executable and keep internal project names stable.
- Ship the companion CLI with complete terminal functionality and optional separately tagged GUI builds.
- Keep CLI credentials in user-only config and verify updates before atomic replacement.

For complete details, see AI.md PART 7, 8, 32.
