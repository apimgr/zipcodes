# Configuration Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never silently ignore unknown or invalid configuration keys.
- Never put credentials in server.yml, logs, URLs, or public responses.
- Never enable debug endpoints merely because mode is development/debug.
- Never add CLI overrides for settings that the spec reserves for config files.

## CRITICAL - ALWAYS DO
- Apply the specified CLI > environment > file > default precedence.
- Validate and normalize config at load; report safe, actionable errors.
- Keep secrets in secrets.yml and runtime config reload behavior explicit.
- Preserve production defaults and keep debug opt-in under the specified precedence.

## KEY DECISIONS
| Question | Answer | Spec Reference |
|---|---|---|
| What does development mode enable? | Verbose logging; sanitization stays enforced and debug endpoints stay off. | PART 6 |
| Where are server settings changed? | Runtime-generated server.yml/drop-ins, with secrets separated. | PART 5, 12 |

For complete details, see AI.md PART 5, 6, 12.
