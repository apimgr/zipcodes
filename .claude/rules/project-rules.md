# Project Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never commit runtime configuration, credentials, databases, logs, or test output.
- Never omit attribution or the license for embedded third-party data.
- Never assume a platform path; use the shared path package and OS-specific rules.

## CRITICAL - ALWAYS DO
- Keep the source tree in the required layout and runtime/test state in project-scoped OS temp paths.
- Preserve required project metadata, licenses, and dataset attribution.
- Use user-only permissions for secrets and private runtime files.

For complete details, see AI.md PART 2, 3, 4.
