# Development

Use the project Docker images and Compose files for development; AI.md requires the Go toolchain to run in a container. The repository Makefile defines build, test, Docker, and cleanup targets. The primary maintainer runs project validation after reviewing changes.

The ZIP dataset is reproduced with `python3 scripts/update_zipcodes_data.py`, which fetches the pinned source revision and writes `src/data/zipcodes.json`. The dataset is CC BY 4.0 and its attribution and exclusions are documented in `data/licenses/zipcodes-data.md`.

Do not add guessed ZIP records or infer area codes from coordinates, states, or counties. Update source attribution whenever the pinned dataset changes.
