# Docker Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never use production or development Compose files for AI-driven runtime testing.
- Never place runtime volumes in the source tree or copy host-built binaries into images.
- Never assume optional overlay services are healthy without checking their status.

## CRITICAL - ALWAYS DO
- Keep Dockerfiles in docker/ and use the prescribed multi-stage, minimal runtime image.
- Run the server as the required unprivileged account under the specified init process.
- Test image contents, health behavior, ports, mounts, and overlay contracts in isolated containers.
- Keep test Compose state under a unique project-scoped temporary directory.

For complete details, see AI.md PART 26.
