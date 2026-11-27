# AI Assistant Rules (PART 0, 1)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never edit AI.md; it is the implementation source of truth and read-only.
- Never guess, skip a relevant PART, add unspecified behavior, or leave stubs/TODOs.
- Never treat code, docs, TODOs, or reports as more authoritative than AI.md.
- Never claim a check passed without running it and seeing its result.

## CRITICAL - ALWAYS DO
- Read the complete relevant PART before changing its code or documentation.
- Follow AI.md exactly; resolve product-description conflicts in IDEA.md.
- Search before adding routes, config keys, dependencies, or tracking items.
- Verify implementation against the relevant PART and FINAL CHECKPOINT.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| What is authoritative? | AI.md defines implementation; IDEA.md defines product intent. | PART 0 |
| Where are template defects recorded? | Parent apimgr TODO.API.md; never TODO.AI.md. | PART 0, FINAL CHECKPOINT |

For complete details, see AI.md PART 0, 1.
