# Binary Rules (PART 7, 8, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never build, test, or run Go binaries with host Go tools.
- Never enable CGO or add undocumented runtime dependencies.
- Never let CLI runtime state use system-owned server directories.
- Never expose server-only administrative operations through the client.

## CRITICAL - ALWAYS DO
- Build static server/client binaries with CGO_ENABLED=0 for all eight OS/architecture pairs.
- Keep server initialization in the binary and server config runtime-generated.
- Implement the server CLI flags and required client CLI/TUI; keep optional GUI separate.
- Keep client config/data/cache/log paths user-scoped with private permissions.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| How are Go binaries built? | Docker `casjaysdev/go:latest`, CGO_ENABLED=0. | PART 7, 25 |
| Is GUI required in every build? | No; native GUI builds are optional and separately tagged. | PART 32 |

For complete details, see AI.md PART 7, 8, 32.
