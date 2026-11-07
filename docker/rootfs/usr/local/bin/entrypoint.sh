#!/usr/bin/env bash
set -eo pipefail

APP_NAME="zipcodes"
APP_BIN="/usr/local/bin/${APP_NAME}"

export TZ="${TZ:-America/New_York}"
export CONFIG_DIR="${CONFIG_DIR:-/config/${APP_NAME}}"
export DATA_DIR="${DATA_DIR:-/data/${APP_NAME}}"
export LOGS_DIR="${LOGS_DIR:-/logs/${APP_NAME}}"
export PORT="${PORT:-80}"
export ADDRESS="${ADDRESS:-0.0.0.0}"

FLAGS=(--address "$ADDRESS" --port "$PORT")
if [[ "${DEBUG:-false}" == "true" ]]; then
    FLAGS+=(--debug)
fi

exec "$APP_BIN" "${FLAGS[@]}" "$@"
