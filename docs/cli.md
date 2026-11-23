# CLI reference

`zipcodes-cli` performs an exact ZIP lookup against a server:

```sh
zipcodes-cli --server https://example.invalid 90210
```

Use `ZIPCODES_SERVER_PRIMARY` or `server.primary` in `cli.yml` to configure a server. No hosted endpoint is defined by this project.

The client supports `--server`, `--version`, and `--help`. The ZIP argument is required. Use `--server` to select a server for one invocation, or configure the server in `cli.yml` or `ZIPCODES_SERVER_PRIMARY`.
