# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never add unversioned API routes outside the named exceptions.
- Never expose sensitive config, credentials, private paths, or internal addresses in health output.
- Never invent aliases or keep removed project routes for backward compatibility.
- Never disable TLS validation or weaken certificate handling.

## CRITICAL - ALWAYS DO
- Keep IDEA.md route inventory, router, OpenAPI, GraphQL, README, and docs synchronized.
- Use the specified response formats, status codes, pagination, and content negotiation.
- Provide health/version output with only public-safe fields.
- Support manual certificates and configured ACME challenge/renewal behavior.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| Where are canonical API routes? | `/api/{api_version}/...`, except explicitly listed exceptions. | PART 14 |
| Which routes are documented as the contract? | AI.md template routes plus project routes in IDEA.md. | PART 14 |

For complete details, see AI.md PART 13, 14, 15.
