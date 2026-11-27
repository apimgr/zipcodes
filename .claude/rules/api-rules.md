# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never retain removed project endpoints as compatibility shims.
- Never add GraphQL mutations; GraphQL is query-only and writes belong to REST.
- Never publish health, version, OpenAPI, or GraphQL output containing secrets.

## CRITICAL - ALWAYS DO
- Keep routes, methods, envelopes, status codes, content negotiation, OpenAPI, and docs synchronized.
- Use the configured API version and canonical path behavior.
- Validate TLS and certificate behavior using the shared SSL manager and secure defaults.

For complete details, see AI.md PART 13, 14, 15.
