# Service Rules (PART 23, 24)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never run service installation, privilege changes, or host system control on the host during development.
- Never run the application as root when the service model supports a dedicated account.
- Never remove data as an incidental effect of stopping or uninstalling a service.

## CRITICAL - ALWAYS DO
- Use the specified privilege-escalation flow and validate the invoking user and executable.
- Keep service install, start, restart, reload, disable, and uninstall semantics distinct.
- Cover supported init systems and OS-specific service templates in isolated system containers.

For complete details, see AI.md PART 23, 24.
