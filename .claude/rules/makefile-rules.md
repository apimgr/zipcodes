# Makefile Rules (PART 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never run Go tools directly on the local machine.
- Never make CI depend on a Makefile target.
- Never place temporary test/build output in the source tree.

## CRITICAL - ALWAYS DO
- Provide `dev`, `local`, `build`, `test`, `release`, and `docker` targets.
- Build with the approved Docker Go image and CGO_ENABLED=0.
- Inject version/build metadata and cross-compile all required platforms.
- Keep temporary output under the project-scoped OS temp directory.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Does CI invoke make? | No; workflows run their required commands directly. | PART 25, 27 |
| Where do local release binaries go? | `binaries/` and `releases/`, as Makefile targets. | PART 25 |

For complete details, see AI.md PART 25.
