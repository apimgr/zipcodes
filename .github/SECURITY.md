# Security policy

## Supported versions

Report vulnerabilities against the latest published release. Older releases
do not have separate security support windows; update to the latest release
before reporting an issue that may already be fixed.

## Reporting a vulnerability

GitHub private vulnerability reporting is the primary channel:
<https://github.com/apimgr/zipcodes/security/advisories/new>. Do not file
vulnerabilities as public issues or publish exploit details before a fix is
available. On a mirror without private reporting, use the GitHub origin link
above. The security email published in `/.well-known/security.txt` is a
secondary contact and notification channel.

Include the affected version or commit, impact, and safe reproduction steps.
The maintainers will acknowledge the report, investigate privately, coordinate
a fix and release, and agree on disclosure timing with the reporter. Do not
include real user data, credentials, tokens, or private keys.

The running service publishes its policy at `/server/security`, machine-readable
contact and encryption information at `/.well-known/security.txt`, and the
report form at `/server/contact?security_id=...`. That form accepts reports
against a deployed instance; use GitHub private reporting for vulnerabilities
in the source project.

## Security model

ZIP lookup endpoints are public and read-only, but server-management features
include anonymous owner tokens, linked tokens, preference controls, and
security-report workflows. Tokens are credentials and must not be posted in
issues, logs, or screenshots. ZIP and GeoIP results are request-scoped; report
content and researcher contact details are encrypted at rest. Deployment
operators control transport security, network exposure, and access-log
retention.
