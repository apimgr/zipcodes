# CI/CD Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never add unpinned third-party workflow actions or bypass security scans.
- Never use Dependabot or inline tool installation in workflow jobs.
- Never treat local success as proof that pushed provider CI passed.

## CRITICAL - ALWAYS DO
- Keep supported GitHub, GitLab, Gitea, Forgejo, and Jenkins pipelines aligned.
- Run equivalent test, lint, security, and cross-platform build gates per provider.
- Pin actions/images as required and use Renovate for dependency updates.
- Check the real CI run after push and fix failures before completion.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Which provider is primary here? | GitHub, inferred from the origin URL. | PART 27 |
| Are pending CI runs complete? | No; completion requires a successful run. | PART 27 |

For complete details, see AI.md PART 27.
