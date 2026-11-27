# Backend Rules (PART 9, 10, 11, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never leak credentials, sensitive request data, or private security reports into logs or notifications.
- Never build SQL by concatenating untrusted values or bypass migrations and validation.
- Never make Tor or I2P a runtime dependency for ordinary service operation.

## CRITICAL - ALWAYS DO
- Use the canonical error, cache, database, encryption, audit, and rate-limit patterns.
- Keep SQLite and remote libSQL paths, backup/restore, and secret rotation aligned with their schemas.
- Treat Tor/I2P as isolated optional services with persistent private state, health, recovery, and shutdown.
- Keep public ZIP lookup read-only; protect ownership and operator operations according to PART 11.

For complete details, see AI.md PART 9, 10, 11, 31.
