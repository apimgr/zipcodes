# CLI reference

`zipcodes-cli` performs an exact ZIP lookup against a server. The server address is selected in this order: `--server`, `ZIPCODES_SERVER_PRIMARY`, then `server.primary` in `cli.yml`.

```sh
zipcodes-cli --server https://example.invalid 90210
```

The client stores `cli.yml` under the operating system's per-user configuration directory in `apimgr/zipcodes/`. When `--server` is used and no primary server is already configured, that address is saved there. The URL must be absolute HTTP or HTTPS and must not contain credentials. This project does not define a hosted API endpoint.

## Client options

| Option | Purpose |
|---|---|
| `--config NAME` | Select a config profile under the per-user configuration directory, or use a path; `.yml` and `.yaml` are detected when the extension is omitted |
| `--server URL` | Select an API server |
| `--token TOKEN` | Validate and use an owner or operator Bearer token; save only when `auth.token` is empty or invalid |
| `--token-file FILE` | Read a Bearer token from a file for this run |
| `--output table\|json\|plain` | Choose the lookup output format |
| `--theme dark\|light\|auto` | Select the terminal theme for this run; `NO_COLOR` disables color |
| `--color auto\|yes\|no` | Control colored terminal output; automatic mode respects `NO_COLOR` |
| `--lang CODE` | Select the output language for this run |
| `--debug[=true\|false]` | Enable or disable diagnostic output |
| `--update check\|yes` | Check for a CLI release or install the verified update |
| `--version` | Print version, commit, build time, Go version, and OS/architecture |
| `--help` | Show usage |
| `--shell completions [SHELL]` | Print completions for bash, zsh, fish, sh, dash, ksh, PowerShell, or pwsh |
| `--shell init [SHELL]` | Print the shell initialization command |

Display mode is detected automatically. With no command, an interactive terminal opens the Bubble Tea lookup screen. Enter a ZIP code and press Enter; use `q` or Escape to leave. A command argument always selects plain CLI mode. `TERM=dumb` and headless runs do not start the interface. The optional `display.mode` setting in `cli.yml` can be `auto`, `tui`, or `gui`; because the default static client does not include native GUI support, `gui` uses the TUI when a terminal is available and otherwise reports `gui_unavailable`.

In CLI mode, the command accepts exactly one ZIP argument containing five digits. For example:

```sh
zipcodes-cli --server http://localhost:8080 90210
```

The client stores `cli.yml` under the per-user `apimgr/zipcodes/` configuration directory with mode `0600` on Unix. Use `--config NAME` for a profile in that directory (`dev` resolves to an existing `dev.yml` or `dev.yaml`), or supply a path. It reads `server.primary`, `auth.token`, `auth.token_file`, `display.mode` (`auto` by default), `theme.mode` (`auto` by default), `output.format` (`table` by default), `defaults.lang`, and the debug/color defaults. Token priority is `--token`, `--token-file`, `ZIPCODES_TOKEN`, `auth.token_file`, then `auth.token`; environment and flag values are not saved. A valid `--server` value is saved only when the stored value is empty or invalid. A `--token` is validated against `/api/{api_version}/server/auth/check`, sent as a Bearer token, and saved only when the stored token is empty or invalid. If no server is configured in an interactive terminal, the Bubble Tea setup screen tests the URL, optionally validates a token, and saves the connection before continuing.

Lookup output precedence is `--output`, `ZIPCODES_OUTPUT_FORMAT`, `output.format`, then `table`. The available formats are `table`, `json`, and `plain`. For example, `zipcodes-cli --server https://example.invalid --output json 90210` prints the server's ZIP record as JSON.

At startup, the client reads autodiscovery and checks the current platform against the server's published version and SHA-256. Interactive runs prompt before an available update; use `--update check` to report status or `--update yes` to download, verify, replace, and restart the CLI. A client below `server.update.cli_min_version` refuses ordinary requests and prints the required update command. The server must have the release binary, `version.txt`, and matching `sha256.txt` provisioned in its CLI directory; see the [API reference](api.md#autodiscovery).

## Server options

The `zipcodes` server binary accepts `--help`, `--version`, and `--status`; `--address`, `--port`, `--config`, `--data`, `--cache`, `--log` / `--logs`, `--backup`, and `--pid`; `--mode`, `--debug`, and `--color`; `--i2p` / `--i2p=false`; `--update [yes|check|branch stable|beta|devel]`; `--service {start|restart|stop|reload|--install|--uninstall|--disable|--help}`; and `--maintenance` commands. The server PID file defaults to `/var/run/apimgr/zipcodes/zipcodes.pid` in system mode and `{data_dir}/zipcodes.pid` in user mode. Container and Windows runs skip PID-file creation because their runtime supervises the process. The explicit I2P flag takes precedence over `I2P_ENABLED`, which takes precedence over `server.i2p.enabled` and its computed default. `--update` defaults to `yes`; `--maintenance update` is an alias. Run `zipcodes --update check` to check without installing. `--update branch {stable|beta|devel}` persists the channel in the server configuration. Branch aliases `daily`, `rolling`, `dev`, and `development` are normalized to `devel`.

`zipcodes --status` queries the local server health endpoint and prints the I2P eepsite state and address when available.

`zipcodes --maintenance tor status` reports the managed Tor process and vanity search state. `zipcodes --maintenance tor validate` checks the Tor settings and generated `torrc`. Both commands query the running server when available and fall back to read-only local state.

Service install, disable, and uninstall actions can alter system services and require the appropriate privileges. `--service --disable` stops and disables the service while retaining its service file and application data. `--service --help` prints service command help and current status.

`zipcodes --maintenance email test` checks the configured SMTP relay and sends one test message to `server.contact.admin.email`. SMTP credentials belong in `secrets.yml`. The command exits with an error when SMTP is unavailable or the admin recipient is unset; it does not queue or attempt delivery without a successful SMTP check.

`zipcodes --maintenance security list [--status STATUS]` lists security-report metadata without report contents. `show TRACKING_ID` displays metadata, the CVE request flag, the researcher comment, and disclosure date, but never the report body. Maintainers initiate CVE assignment after triage when requested. `export TRACKING_ID [--output FILE]` writes the already encrypted report body to a mode-`0600` file named `TRACKING_ID.asc` by default; `--output -` writes the encrypted bytes to standard output. `update TRACKING_ID` accepts `--status`, `--comment`, `--disclosure-date YYYY-MM-DD`, and `--severity`; status can only move forward or to `wont_fix`. These actions require root or interactive operator-token confirmation. A status change sends a status-only notice when SMTP passes its relay check and an encrypted contact record exists; the notice is encrypted with the researcher's submitted OpenPGP key when available.

`zipcodes --maintenance pgp generate` creates a two-year Ed25519/Curve25519 project keypair and stores its private key encrypted under the database `installation_secret`. Generation publishes to configured keyservers. `rotate` stages a replacement key, re-encrypts all retained PGP security reports in one database transaction, then atomically activates the new key; startup recovery resolves an interrupted rotation before accepting submissions. The previous encrypted private key is retained for 30 days for reports already in transit. `export public [FILE]` prints or writes the public key. `export private FILE` decrypts the private key only after a second operator-token prompt and a required reason; exports are limited to one per hour and written mode `0600`. `import FILE` imports an armored private key when no project key exists. `publish` submits the current public key to configured keyservers; due failures retry on the 15-minute `keyserver_retry` task, up to eight attempts. `delete` requires typing `DELETE KEY`, removes both key files, records a revocation state, and disables public-key publication.

The scheduled `update_check` task runs daily at 06:00. It checks stable releases by default; beta and devel channels include all more stable releases. `server.update.defer_days` delays scheduled stable/beta notices and installs, while manual checks ignore the delay. `server.update.auto_install` is disabled by default. Update notices are operator-only and never appear in public API responses.

Automatic staged restarts rely on systemd, launchd, runit, s6, or rc.d supervising exit code 75. OpenRC and SysV leave a verified update staged until the next service or system start. Windows services do not auto-install; run the update command from an elevated prompt when appropriate.
