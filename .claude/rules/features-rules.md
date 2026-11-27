# Feature Rules (PART 17, 18, 19, 20, 21, 22)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never send private report content or credentials through webhook summaries.
- Never expose internal metrics or allow unverified, unauthenticated update installs.
- Never restore untrusted backup contents without validation and atomicity.

## CRITICAL - ALWAYS DO
- Implement applicable notification, scheduling, GeoIP, metrics, backup/restore, and update flows completely.
- Redact sensitive provider errors and audit state-changing operations.
- Use bounded retries, retention, checksums/signatures, and recovery behavior from each PART.
- Keep operator controls separate from public read-only ZIP routes.

For complete details, see AI.md PART 17, 18, 19, 20, 21, 22.
