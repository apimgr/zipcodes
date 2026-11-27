# Makefile Rules (PART 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never add non-spec Make targets or use Make as the CI workflow definition.
- Never run Go tooling locally or let clean remove outputs outside BINDIR/RELDIR.
- Never publish releases, tags, or images as part of validation.

## CRITICAL - ALWAYS DO
- Keep the seven core targets and only the explicitly permitted housekeeping/helper targets.
- Run Go tooling in bounded, resource-limited casjaysdev/go:latest containers with CGO disabled.
- Build all required server/client platforms and honor configurable output/cache paths.
- Keep release and beta behavior aligned with the artifact and publication contract.

For complete details, see AI.md PART 25.
