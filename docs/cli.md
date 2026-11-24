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
| `--version` | Print version, commit, build time, Go version, and OS/architecture |
| `--help` | Show usage |
| `--shell completions [SHELL]` | Print completions for bash, zsh, fish, sh, dash, ksh, PowerShell, or pwsh |
| `--shell init [SHELL]` | Print the shell initialization command |

The command accepts exactly one ZIP argument containing five digits. For example:

```sh
zipcodes-cli --server https://your-server.example 90210
```

## Server options

The `zipcodes` server binary accepts `--help`, `--version`, and `--status`; `--address`, `--port`, `--config`, `--data`, and `--logs`; `--mode`, `--debug`, and `--color`; `--update [yes|check|branch stable|beta|devel]`; `--service install|uninstall|start|stop|restart|status`; and `--maintenance on|off|update|backup|restore`. `--update` defaults to `yes`; `--maintenance update` is an alias. Run `zipcodes --update check` to check without installing. Installing requires an elevated process. `--update branch {stable|beta|devel}` persists the channel in the server configuration. Branch aliases `daily`, `rolling`, `dev`, and `development` are normalized to `devel`.

The scheduled `update_check` task runs daily at 06:00. It checks stable releases by default; beta and devel channels include all more stable releases. `server.update.defer_days` delays scheduled stable/beta notices and installs, while manual checks ignore the delay. `server.update.auto_install` is disabled by default. Update notices are operator-only and never appear in public API responses.

Automatic staged restarts rely on systemd, launchd, runit, s6, or rc.d supervising exit code 75. OpenRC and SysV leave a verified update staged until the next service or system start. Windows services do not auto-install; run the update command from an elevated prompt when appropriate.
