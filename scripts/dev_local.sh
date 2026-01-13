#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

backend_pid=""
frontend_pid=""

ensure_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "$cmd not found in PATH. Please install it and retry." >&2
    exit 1
  fi
}

ensure_go_deps() {
  if [[ ! -f "$ROOT_DIR/go.sum" ]]; then
    echo "go.sum missing, running go mod tidy..."
    (cd "$ROOT_DIR" && GOFLAGS="-mod=mod" go mod tidy)
  fi
}

ensure_config() {
  local config_path="$ROOT_DIR/configs/nginxpulse_config.json"
  if [[ ! -f "$config_path" ]]; then
    echo "configs/nginxpulse_config.json not found, generating default config..."
    (cd "$ROOT_DIR" && go run ./cmd/nginxpulse -gen-config)
    echo "Edit configs/nginxpulse_config.json and re-run." >&2
    exit 1
  fi
}

ensure_frontend_deps() {
  local install_needed=0
  if [[ ! -d "$ROOT_DIR/webapp/node_modules" ]]; then
    install_needed=1
  elif [[ "$ROOT_DIR/webapp/package.json" -nt "$ROOT_DIR/webapp/node_modules" ]]; then
    install_needed=1
  elif [[ -f "$ROOT_DIR/webapp/package-lock.json" && "$ROOT_DIR/webapp/package-lock.json" -nt "$ROOT_DIR/webapp/node_modules" ]]; then
    install_needed=1
  fi

  if [[ "$install_needed" -eq 1 ]]; then
    echo "Installing frontend dependencies..."
    if [[ -f "$ROOT_DIR/webapp/package-lock.json" ]]; then
      (cd "$ROOT_DIR/webapp" && npm ci) || (cd "$ROOT_DIR/webapp" && npm install)
    else
      (cd "$ROOT_DIR/webapp" && npm install)
    fi
  fi
}

start_backend() {
  echo "Starting backend on http://localhost:8089"
  (cd "$ROOT_DIR" && SERVER_PORT=":8089" go run ./cmd/nginxpulse) &
  backend_pid=$!
}

start_frontend() {
  echo "Starting frontend on http://localhost:8088"
  (cd "$ROOT_DIR/webapp" && npm run dev) &
  frontend_pid=$!
}

cleanup() {
  if [[ -n "$frontend_pid" ]]; then
    kill "$frontend_pid" >/dev/null 2>&1 || true
  fi
  if [[ -n "$backend_pid" ]]; then
    kill "$backend_pid" >/dev/null 2>&1 || true
  fi
}

trap cleanup EXIT INT TERM

ensure_cmd go
ensure_cmd node
ensure_cmd npm
ensure_go_deps
ensure_config
ensure_frontend_deps

start_backend
start_frontend

wait
