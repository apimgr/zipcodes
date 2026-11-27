# Project Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never add runtime config, logs, test output, or build artifacts to source directories.
- Never change stable on-disk identifiers from IDEA.md without the specified migration.
- Never create extra top-level source roots prohibited by the project layout.

## CRITICAL - ALWAYS DO
- Keep project structure and license/attribution compliant with the spec.
- Resolve paths by OS and system/user mode using `{internal_org}/{internal_name}`.
- Create private directories/files with restrictive permissions at creation time.
- Keep runtime-generated server.yml out of version control.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Which identity is stable on disk? | `internal_name` and `internal_org` from IDEA.md. | PART 2, 4 |
| Where does runtime data go? | OS-specific config/data/log/cache paths, not the source tree. | PART 3, 4 |

For complete details, see AI.md PART 2, 3, 4.
