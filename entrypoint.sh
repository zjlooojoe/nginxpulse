#!/bin/sh
set -e

DATA_DIR="${DATA_DIR:-/app/var/nginxpulse_data}"

APP_UID="${PUID:-}"
APP_GID="${PGID:-}"
APP_USER="nginxpulse"
APP_GROUP="nginxpulse"

if [ -n "$APP_GID" ]; then
  EXISTING_GROUP="$(awk -F: -v gid="$APP_GID" '$3==gid{print $1; exit}' /etc/group)"
  if [ -z "$EXISTING_GROUP" ]; then
    addgroup -S -g "$APP_GID" appgroup
    APP_GROUP="appgroup"
  else
    APP_GROUP="$EXISTING_GROUP"
  fi
fi

if [ -n "$APP_UID" ]; then
  EXISTING_USER="$(awk -F: -v uid="$APP_UID" '$3==uid{print $1; exit}' /etc/passwd)"
  if [ -z "$EXISTING_USER" ]; then
    adduser -S -D -H -u "$APP_UID" -G "$APP_GROUP" appuser
    APP_USER="appuser"
  else
    APP_USER="$EXISTING_USER"
  fi
fi

mkdir -p "$DATA_DIR"

if [ "$(id -u)" = "0" ]; then
  if ! su-exec "$APP_USER:$APP_GROUP" sh -lc "touch '$DATA_DIR/.write_test' && rm -f '$DATA_DIR/.write_test'" >/dev/null 2>&1; then
    chown -R "$APP_USER:$APP_GROUP" "$DATA_DIR" 2>/dev/null || true
  fi
fi

if ! su-exec "$APP_USER:$APP_GROUP" sh -lc "touch '$DATA_DIR/.write_test' && rm -f '$DATA_DIR/.write_test'" >/dev/null 2>&1; then
  echo "nginxpulse: $DATA_DIR is not writable; file logging may fail and will fall back to stdout" >&2
fi

if command -v nginx >/dev/null 2>&1; then
  su-exec "$APP_USER:$APP_GROUP" /app/nginxpulse "$@" &
  backend_pid=$!
  nginx -g 'daemon off;' &
  nginx_pid=$!

  shutdown() {
    if [ -n "${backend_pid:-}" ]; then
      kill -TERM "$backend_pid" >/dev/null 2>&1 || true
    fi
    if [ -n "${nginx_pid:-}" ]; then
      kill -TERM "$nginx_pid" >/dev/null 2>&1 || true
    fi
  }

  trap shutdown INT TERM

  while :; do
    if [ -n "${backend_pid:-}" ] && ! kill -0 "$backend_pid" >/dev/null 2>&1; then
      shutdown
      wait "$backend_pid" >/dev/null 2>&1 || true
      exit 1
    fi
    if [ -n "${nginx_pid:-}" ] && ! kill -0 "$nginx_pid" >/dev/null 2>&1; then
      shutdown
      wait "$nginx_pid" >/dev/null 2>&1 || true
      exit 1
    fi
    sleep 1
  done
fi

exec su-exec "$APP_USER:$APP_GROUP" /app/nginxpulse "$@"
