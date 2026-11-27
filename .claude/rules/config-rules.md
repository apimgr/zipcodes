# Configuration Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never check generated server configuration or secrets into source control.
- Never silently accept invalid configuration, expose secrets, or conflate development with debug mode.
- Never replace an operator's valid setting with an invalid default or flag value.

## CRITICAL - ALWAYS DO
- Generate server.yml and secrets.yml at runtime with restrictive permissions.
- Follow the documented precedence among flags, environment, files, and embedded defaults.
- Load ordered yaml.d area files and validate values before activation.
- Keep secret fields separate from public configuration and redact them in output.

For complete details, see AI.md PART 5, 6, 12.
