# Contributing to Zipcodes

## Local setup and checks

Use Docker with `casjaysdev/go:latest` for Go commands. The repository
Makefile provides `make dev`, `make build`, and `make test`. Run documentation
checks from the repository root with `mkdocs build --strict`; Python
dependencies are listed in `docs/requirements.txt`.

Before opening a pull request, run the test and lint gates, add or update tests
for behavior changes, update the relevant pages under `docs/`, and confirm the
API, browser, CLI, and documentation descriptions agree.

## Branches and pull requests

Create a topic branch from `main`, keep changes focused, and open a pull
request. The protected default branch requires review and passing CI checks;
do not force-push or delete it. The pull request template lists the evidence
and impact details required for review.

## Security reports

Do not open a public issue for a vulnerability. Use GitHub's private reporting
flow at https://github.com/apimgr/zipcodes/security/advisories/new. See
[SECURITY.md](SECURITY.md) for the response and disclosure process.

## Project guides

- [Development](../docs/development.md)
- [API](../docs/api.md)
- [Security](../docs/security.md)
- [Translation](../docs/translating.md)
