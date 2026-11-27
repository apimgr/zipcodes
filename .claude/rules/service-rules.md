# Service Rules (PART 23, 24)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never perform a protected operation without checking required privileges.
- Never silently claim service lifecycle success when the platform operation failed.
- Never use system service identities or shared paths for user-level services.

## CRITICAL - ALWAYS DO
- Explain privilege requirements and use only supported escalation methods.
- Implement install, start, stop, restart, reload, and uninstall per supported OS.
- Use stable internal names and OS-native service managers.
- Verify generated service definitions and safely clean only resources owned by this service.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| What identifies a service across renames? | Stable `internal_name`/`plist_name`. | PART 2, 24 |
| How are system and user services separated? | Distinct identities, paths, and native manager units. | PART 23, 24 |

For complete details, see AI.md PART 23, 24.
