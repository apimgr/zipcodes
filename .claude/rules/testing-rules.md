# Testing, Documentation, and I18N Rules (PART 28, 29, 30)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never run Go tests/builds/binaries on the local machine; use the specified Docker image.
- Never leave runtime data, coverage files, or test artifacts in the source tree.
- Never claim documentation deployment exists unless its endpoint responds.
- Never ship missing locale keys or an LTR-only layout for Arabic.

## CRITICAL - ALWAYS DO
- Run unit, integration, browser, Docker, and Incus checks required by the relevant scenario.
- Keep total Go coverage at or above the Makefile threshold.
- Keep MkDocs/ReadTheDocs pages, README, CLI, API, and implementation synchronized.
- Test accessible server-rendered pages, no-JS forms, all supported locales, and RTL behavior.
- Use the required wrapper for test/lint commands and verify actual results.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Where does coverage output go? | Project-scoped OS temp directory. | PART 28 |
| When is the docs badge shown? | Only when an active documentation deployment responds. | PART 29 |

For complete details, see AI.md PART 28, 29, 30.
