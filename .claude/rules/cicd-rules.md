# CI/CD Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never replace provider workflows with Makefile targets or run publish workflows as local verification.
- Never use floating third-party action references when full commit-SHA pins are required.
- Never commit secrets or skip required security and artifact checks.

## CRITICAL - ALWAYS DO
- Keep GitHub, GitLab, Gitea, Forgejo, and Jenkins workflows aligned with the supported provider contract.
- Run lint, tests, eight-platform builds, vulnerability scans, secret scans, and image scans where prescribed.
- Validate workflow syntax and provider-specific behavior, then inspect CI after an authorized push.

For complete details, see AI.md PART 27.
