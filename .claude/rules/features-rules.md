# Feature Rules (PART 17, 18, 19, 20, 21, 22)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never depend on external cron or schedulers.
- Never expose operator notifications in the public WebUI.
- Never retain visitor location as a side effect of GeoIP lookup.
- Never prune valid backups before newly created backups pass verification.
- Never auto-install updates unless explicitly enabled by configuration.

## CRITICAL - ALWAYS DO
- Keep notifications graceful without SMTP; use logs for operator events.
- Run the persistent built-in scheduler with bounded retry, catch-up, safe recovery, and task isolation.
- Make metrics endpoints bearer-token protected and keep GeoIP failure non-fatal.
- Verify backup checksum, decryption, extraction, manifest, and database integrity.
- Keep update channels, staging, rollback, and operator confirmations per spec.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Where is backup retention configured? | `server.backup.retention`, not scheduler task settings. | PART 21 |
| Is GeoIP required for ZIP lookup? | No; lookup works when GeoIP is unavailable. | PART 19 |

For complete details, see AI.md PART 17, 18, 19, 20, 21, 22.
