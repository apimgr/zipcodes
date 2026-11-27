# Docker Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never run the application as a non-root container user when runtime privilege setup is required.
- Never put generated server config or persistent data inside the image.
- Never omit tini, health checks, OCI metadata, or required volume mounts.

## CRITICAL - ALWAYS DO
- Use the specified multi-stage Alpine image and keep runtime dependencies embedded or explicitly allowed.
- Run the entrypoint through tini and let the application establish runtime ownership/permissions.
- Keep production, development, and test Compose files aligned with the specification.
- Bind persistent config/data to documented paths and avoid committing runtime volumes.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Which processes may be external? | Only documented optional overlay daemons. | PART 7, 26, 31 |
| Who manages runtime permissions? | The application at startup. | PART 26 |

For complete details, see AI.md PART 26.
