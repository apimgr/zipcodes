# CLI reference

`zipcodes-cli` performs an exact ZIP lookup against a server. The server address is selected in this order: `--server`, `ZIPCODES_SERVER_PRIMARY`, then `server.primary` in `cli.yml`.

```sh
zipcodes-cli --server https://example.invalid 90210
```

The client stores `cli.yml` under the operating system's per-user configuration directory in `apimgr/zipcodes/`. When `--server` is used and no primary server is already configured, that address is saved there. The URL must be absolute HTTP or HTTPS and must not contain credentials. This project does not define a hosted API endpoint.

## Client options

| Option | Purpose |
|---|---|
| `--server URL` | Select an API server |
| `--tui` | Start an interactive repeated ZIP lookup interface; enter `q` to exit |
| `--theme dark\|light\|auto` | Select the TUI prompt theme; dark is the default and `NO_COLOR` disables color |
| `--version` | Print version, commit, build time, Go version, and OS/architecture |
| `--help` | Show usage |
| `--shell completions [SHELL]` | Print completions for bash, zsh, fish, sh, dash, ksh, PowerShell, or pwsh |
| `--shell init [SHELL]` | Print the shell initialization command |

The command accepts exactly one ZIP argument containing five digits. For example:

```sh
zipcodes-cli --server https://your-server.example 90210
```

To look up multiple ZIP codes in one session, run `zipcodes-cli --tui`. The selected theme is stored in `cli.yml` under the per-user `apimgr/zipcodes/` configuration directory. A `--server` value is saved as the primary only when no primary server is already configured.

Before a lookup, the client reads the server's autodiscovery document. If `server.update.cli_min_version` is set and the client version is older, the client refuses the lookup and prints the required version and update command.

## Server options

The `zipcodes` server binary accepts `--help`, `--version`, and `--status`; `--address`, `--port`, `--config`, `--data`, `--logs`, and `--backup`; `--mode`, `--debug`, and `--color`; `--update [yes|check|branch stable|beta|devel]`; `--service {start|restart|stop|reload|--install|--uninstall|--disable|--help}`; and `--maintenance` commands. `--update` defaults to `yes`; `--maintenance update` is an alias. Run `zipcodes --update check` to check without installing. `--update branch {stable|beta|devel}` persists the channel in the server configuration. Branch aliases `daily`, `rolling`, `dev`, and `development` are normalized to `devel`.

Service install, disable, and uninstall actions can alter system services and require the appropriate privileges. `--service --disable` stops and disables the service while retaining its service file and application data. `--service --help` prints service command help and current status.

`zipcodes --maintenance email test` checks the configured SMTP relay and sends one test message to `server.contact.admin.email`. SMTP credentials belong in `secrets.yml`. The command exits with an error when SMTP is unavailable or the admin recipient is unset; it does not queue or attempt delivery without a successful SMTP check.

The scheduled `update_check` task runs daily at 06:00. It checks stable releases by default; beta and devel channels include all more stable releases. `server.update.defer_days` delays scheduled stable/beta notices and installs, while manual checks ignore the delay. `server.update.auto_install` is disabled by default. Update notices are operator-only and never appear in public API responses.

Automatic staged restarts rely on systemd, launchd, runit, s6, or rc.d supervising exit code 75. OpenRC and SysV leave a verified update staged until the next service or system start. Windows services do not auto-install; run the update command from an elevated prompt when appropriate.
