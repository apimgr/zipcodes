# Security policy

## Supported versions

Report issues against the latest published release. The project does not maintain separate security support windows for older releases.

## Reporting a vulnerability

Use GitHub's private reporting flow: <https://github.com/apimgr/zipcodes/security/advisories/new>. Do not publish an exploitable vulnerability before a fix is available. Include affected versions, impact, and steps to reproduce without disclosing real user data or credentials.

The service is intentionally public and read-only. It does not provide accounts, API keys, or write endpoints. Rate limits are abuse controls and do not create access tiers.

## Data and privacy

ZIP records are public geographic reference data. The `/api/v1/me` endpoint derives location from the caller IP when GeoIP data is available; lookup results are request-scoped. Deployment operators control transport security, network exposure, and access-log retention.
