# CLI reference

`zipcodes-cli` performs ZIP lookups and can invoke any API route on the configured server. The server address is selected in this order: `--server`, `ZIPCODES_SERVER_PRIMARY`, then `server.primary` in `cli.yml`.

```sh
zipcodes-cli --server https://example.invalid 90210
zipcodes-cli --server https://example.invalid api GET /api/v2/zipcodes/90210
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
| `--debug[=BOOL]` | Enable or disable diagnostic output; values use the shared truthy/falsy parser |
| `--update check\|yes` | Check for a CLI release or install the verified update |
| `--version` | Print version, commit, build time, Go version, and OS/architecture |
| `--help` | Show usage |
| `--shell completions [SHELL]` | Print completions for bash, zsh, fish, sh, dash, ksh, PowerShell, or pwsh |
| `--shell init [SHELL]` | Print the shell initialization command |

Display mode is detected automatically. With no command, a native GUI build opens its home screen on a supported local display; otherwise an interactive terminal opens the Bubble Tea menu. The menu and GUI expose the same ZIP command tree, in the same order as help: `lookup`, `search`, `nearby`, `city`, `state`, `autocomplete`, `stats`, `random`, and `dataset`, followed by the generic API request screen. The optional `display.mode` setting in `cli.yml` can be `auto`, `tui`, or `gui`. A GUI build is a separate optional artifact produced with `make build BUILD_GUI=1`; its filenames use `zipcodes-cli-gui-{os}-{arch}`. The default static client does not include GUI support; an explicit `display.mode: gui` reports `gui_unavailable` and never falls back to the TUI. Remote SSH/Mosh sessions use the TUI, and a command argument always selects CLI mode. `TERM=dumb` and headless runs do not start an interface.

In CLI mode, a ZIP argument performs an exact lookup and must contain five digits; `lookup ZIPCODE` is the explicit form. The project commands map to the server's ZIP API: `search QUERY` and `autocomplete QUERY` accept multiword text, `nearby ZIPCODE [RADIUS_MILES]` accepts an optional radius, `city CITY` accepts a city name, `state STATE_ABBREVIATION` requires a two-letter state, and `stats`, `random`, and `dataset` take no arguments. User text is encoded into the URL. The `api` command accepts an HTTP method, a relative `/api/` path, and an optional JSON body. It sends the request to the configured server, uses the resolved Bearer token, and prints the response body (pretty-printed when it is JSON). Absolute URLs and paths outside `/api/` are rejected.

```sh
zipcodes-cli lookup 90210
zipcodes-cli search Los Angeles
zipcodes-cli nearby 90210 25
zipcodes-cli city Los Angeles
zipcodes-cli state CA
zipcodes-cli autocomplete Los
zipcodes-cli stats
zipcodes-cli random
zipcodes-cli dataset
zipcodes-cli api GET '/api/v2/zipcodes/search?q=Los%20Angeles&state=CA'
zipcodes-cli api POST /api/v2/server/preferences '{"theme":"dark"}'
```

For example, an exact lookup can still be run directly:

```sh
zipcodes-cli --server http://localhost:8080 90210
```

The client stores `cli.yml` under the per-user `apimgr/zipcodes/` configuration directory with mode `0600` on Unix. Use `--config NAME` for a profile in that directory (`dev` resolves to an existing `dev.yml` or `dev.yaml`), or supply a path. It reads `server.primary`, `auth.token`, `auth.token_file`, `display.mode` (`auto` by default), `theme.mode` (`auto` by default), `output.format` (`table` by default), `defaults.lang`, and the debug/color defaults. Token priority is `--token`, `--token-file`, `ZIPCODES_TOKEN`, `auth.token_file`, then `auth.token`; environment and flag values are not saved. The resolved language is sent in `Accept-Language` on autodiscovery, auth checks, lookups, API calls, and binary updates. A valid `--server` value is saved only when the stored value is empty or invalid. A `--token` is validated against `/api/{api_version}/server/auth/check`, sent as a Bearer token, and saved only when the stored token is empty or invalid. On first GUI launch, the setup screen tests the server URL and optional token before connecting. Saving the connection stores it in `cli.yml`; selecting automatic login startup writes a per-user startup entry for Linux, macOS, or Windows. TUI setup remains available in terminal sessions.

Lookup output precedence is `--output`, `ZIPCODES_OUTPUT_FORMAT`, `output.format`, then `table`. The available formats are `table`, `json`, and `plain`. For example, `zipcodes-cli --server https://example.invalid --output json 90210` prints the server's ZIP record as JSON.

At startup, the client reads autodiscovery and checks the current platform against the server's published version and SHA-256. Interactive runs prompt before an available update; use `--update check` to report status or `--update yes` to download, verify, replace, and restart the CLI. A client below `server.update.cli_min_version` refuses ordinary requests and prints the required update command. The server must have the release binary, `version.txt`, and matching `sha256.txt` provisioned in its CLI directory; see the [API reference](api.md#autodiscovery).

## Server options

`--daemon` and `server.daemonize: true` detach manual starts. Service starts and detected container runs stay in the foreground so their supervisor retains control.

The `zipcodes` server binary accepts `--help`, `--version`, and `--status`; `--address`, `--port`, `--config`, `--data`, `--cache`, `--log` / `--logs`, `--backup`, and `--pid`; `--mode`, `--debug`, and `--color`; `--i2p` / `--i2p=false`; `--update [yes|check|branch stable|beta|devel]`; `--service {start|restart|stop|reload|--install|--uninstall|--disable|--help}`; and `--maintenance` commands. Boolean options accept case-insensitive truthy values (`1`, `y`, `t`, `yes`, `true`, `on`, `ok`, `enable`, `enabled`, `yep`, `yup`, `yeah`, `aye`, `si`, `oui`, `da`, `hai`, `affirmative`, `accept`, `allow`, `grant`, `sure`, `totally`) and falsy values (`0`, `n`, `f`, `no`, `false`, `off`, `disable`, `disabled`, `nope`, `nah`, `nay`, `nein`, `non`, `niet`, `iie`, `lie`, `negative`, `reject`, `block`, `revoke`, `deny`, `never`, `noway`, `none`). Both `--flag=value` and `--flag value` are accepted. Mode precedence is `--mode`, `MODE`, `server.mode`, then `production`; debug precedence is an explicitly supplied `--debug`, `DEBUG`, the `debug` mode default, then disabled. The PID file defaults to `/var/run/apimgr/zipcodes/zipcodes.pid` in system mode and `{data_dir}/zipcodes.pid` in user mode; `--pid` overrides `PID_FILE`, which overrides that default. Container and Windows runs skip PID-file creation because their runtime supervises the process. The explicit I2P flag takes precedence over `I2P_ENABLED`, which takes precedence over `server.i2p.enabled` and its computed default. `--update` defaults to `yes`; `--maintenance update` is an alias. Run `zipcodes --update check` to check without installing. `--update branch {stable|beta|devel}` persists the channel in the server configuration. Branch aliases `daily`, `rolling`, `dev`, and `development` are normalized to `devel`.

`zipcodes --status` queries the local server health endpoint and prints the I2P eepsite state and address when available.

`zipcodes --maintenance tor status` reports the managed Tor process and vanity search state. `zipcodes --maintenance tor validate` checks the Tor settings and generated `torrc`. Both commands query the running server when available and fall back to read-only local state.

As root, service installation configures the detected system service. A non-root install uses a systemd user unit or macOS LaunchAgent when available, runs as the current user, and cannot bind ports below 1024. `--service --disable` stops and disables the selected service while retaining its service file and application data. `--service --help` prints service command help and current status.

`zipcodes --maintenance email test` checks the configured SMTP relay and sends one test message to `server.contact.admin.email`. SMTP credentials belong in `secrets.yml`. The command exits with an error when SMTP is unavailable or the admin recipient is unset; it does not queue or attempt delivery without a successful SMTP check.

`zipcodes --maintenance security list [--status STATUS]` lists security-report metadata without report contents. `show TRACKING_ID` displays metadata, the CVE request flag, the researcher comment, and disclosure date, but never the report body. Maintainers initiate CVE assignment after triage when requested. `export TRACKING_ID [--output FILE]` writes the already encrypted report body to a mode-`0600` file named `TRACKING_ID.asc` by default; `--output -` writes the encrypted bytes to standard output. `update TRACKING_ID` accepts `--status`, `--comment`, `--disclosure-date YYYY-MM-DD`, and `--severity`; status can only move forward or to `wont_fix`. These actions require root or interactive operator-token confirmation. A status change sends a status-only notice when SMTP passes its relay check and an encrypted contact record exists; the notice is encrypted with the researcher's submitted OpenPGP key when available.

`zipcodes --maintenance pgp generate` creates a two-year Ed25519/Curve25519 project keypair and stores its private key encrypted under the database `installation_secret`. Generation publishes to configured keyservers. `rotate` stages a replacement key, re-encrypts all retained PGP security reports in one database transaction, then atomically activates the new key; startup recovery resolves an interrupted rotation before accepting submissions. The previous encrypted private key is retained for 30 days for reports already in transit. `export public [FILE]` prints or writes the public key. `export private FILE` requires typing `EXPORT PRIVATE KEY` and entering a reason; operators authorized by `server.token` re-enter the token, while root does not. Exports are limited to one per hour per OS user and written mode `0600`; the audit event records the OS username, `local` operator IP, and reason. `import FILE` imports an armored private key when no project key exists and requires typing `IMPORT KEY`; an identity mismatch requires a separate `IMPORT MISMATCH` confirmation. `publish` submits the current public key to configured keyservers; due failures retry on the 15-minute `keyserver_retry` task, up to eight attempts. `delete` requires typing `DELETE KEY`, removes both keys, records a revocation state, and disables public-key publication.

The scheduled `update_check` task runs daily at 06:00. It checks stable releases by default; beta and devel channels include all more stable releases. `server.update.defer_days` delays scheduled stable/beta notices and installs, while manual checks ignore the delay. `server.update.auto_install` is disabled by default. Update notices are operator-only and never appear in public API responses.

Automatic staged restarts rely on systemd, launchd, runit, s6, or rc.d supervising exit code 75. OpenRC and SysV leave a verified update staged until the next service or system start. Windows services do not auto-install; run the update command from an elevated prompt when appropriate.
