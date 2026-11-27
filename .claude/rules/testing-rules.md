# Testing Rules (PART 28, 29, 30)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never test with runtime data in the source tree or run project binaries on the host.
- Never claim coverage or behavior that wrapped output did not verify.
- Never ship hardcoded human-readable strings in a user-facing surface.

## CRITICAL - ALWAYS DO
- Keep unit coverage at or above the specified threshold and run Go tests in the toolchain container.
- Use bounded, isolated Docker/Incus scripts for binary, browser, and service lifecycle checks.
- Keep ReadTheDocs/MkDocs pages, navigation, CLI/API/config docs, and IDEA.md current.
- Use one locale registry and complete translations; verify keyboard and screen-size accessibility.

For complete details, see AI.md PART 28, 29, 30.
