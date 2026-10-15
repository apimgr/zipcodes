# CLI reference

`zipcodes-cli` performs exact ZIP lookups against a server:

```sh
zipcodes-cli --server https://example.invalid 90210
```

Use `ZIPCODES_SERVER_PRIMARY` or `server.primary` in `cli.yml` to configure a server. The server URL is required because no official hosted endpoint is specified.

The current client supports `--server`, `--version`, and `--help`. Shell completions and the complete shared client configuration/interactive modes are not yet implemented.
