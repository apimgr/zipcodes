# Backend Rules (PART 9, 10, 11, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never interpolate untrusted input into SQL, paths, commands, HTML, or logs.
- Never expose secrets, private report content, internal paths, or private blobs.
- Never let one task/request panic terminate the server or scheduler.
- Never require Tor or I2P binaries for the server to run.

## CRITICAL - ALWAYS DO
- Use embedded SQLite by default and the specified remote libSQL/Turso support.
- Run schema creation/migrations and integrity checks through managed startup.
- Enforce security headers, rate limits, structured audit/security logs, and rotation.
- Keep private data private; use authenticated private/no-store delivery when applicable.
- Keep overlay processes dedicated, optional as specified, monitored, and recoverable.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Are GeoIP/location lookups retained? | No; request location is informational and not persisted as a lookup side effect. | PART 19, IDEA.md |
| Are user uploads executable? | No; validate, constrain, and render untrusted content safely. | PART 11 |

For complete details, see AI.md PART 9, 10, 11, 31.
