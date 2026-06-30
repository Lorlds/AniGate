#!/usr/bin/env sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PREFIX=${PREFIX:-"$HOME/.local"}
BINDIR=${BINDIR:-"$PREFIX/bin"}
CONFIG_DIR=${ANIGATE_CONFIG_DIR:-"${XDG_CONFIG_HOME:-"$HOME/.config"}/anigate"}
STATE_DIR=${ANIGATE_STATE_DIR:-"${XDG_STATE_HOME:-"$HOME/.local/state"}/anigate"}
CONFIG_FILE=${ANIGATE_CONFIG_FILE:-"$CONFIG_DIR/anigate.json"}
WORKSPACE_DIR=${ANIGATE_WORKSPACE_DIR:-"$HOME"}
TOKEN=${ANIGATE_AUTH_TOKEN:-}

if [ -z "$TOKEN" ]; then
  if command -v openssl >/dev/null 2>&1; then
    TOKEN=$(openssl rand -hex 24)
  else
    TOKEN=$(date +%s | sha256sum | awk '{print $1}')
  fi
fi

mkdir -p "$BINDIR" "$CONFIG_DIR" "$STATE_DIR"

echo "==> building $BINDIR/anigate"
cd "$ROOT"
go build -trimpath -ldflags "-s -w" -o "$BINDIR/anigate" ./cmd/anigate

if [ ! -f "$CONFIG_FILE" ]; then
  echo "==> writing $CONFIG_FILE"
  cat >"$CONFIG_FILE" <<EOF
{
  "state_dir": "$STATE_DIR",
  "auth_token": "$TOKEN",
  "max_read_bytes": 65536,
  "max_search_file_bytes": 262144,
  "max_search_results": 50,
  "max_job_log_bytes": 1048576,
  "max_artifact_bytes": 4194304,
  "env_allowlist": ["OPENAI_API_KEY", "ANTHROPIC_API_KEY"],
  "isolated_home": true,
  "workspaces": [
    {
      "name": "home",
      "path": "$WORKSPACE_DIR",
      "read_only": true,
      "profile": "reader"
    }
  ],
  "presets": [
    {
      "name": "sys_uptime",
      "description": "Show system uptime",
      "workspace": "home",
      "cwd": ".",
      "command": ["uptime"],
      "timeout_sec": 10,
      "async": false
    }
  ],
  "agents": [
    {
      "name": "echo_agent",
      "description": "Development placeholder that echoes the session prompt",
      "provider": "echo",
      "workspace": "home",
      "cwd": ".",
      "command": ["printf", "{prompt}"],
      "timeout_sec": 30,
      "max_history_messages": 20
    }
  ]
}
EOF
  chmod 600 "$CONFIG_FILE"
else
  echo "==> keeping existing $CONFIG_FILE"
fi

cat <<EOF

AniGate installed.

Binary:
  $BINDIR/anigate

Config:
  $CONFIG_FILE

Try stdio mode:
  $BINDIR/anigate stdio --config "$CONFIG_FILE"

Try local HTTP mode:
  $BINDIR/anigate http --addr 127.0.0.1:8787 --config "$CONFIG_FILE"

HTTP token header:
  Authorization: Bearer <token from $CONFIG_FILE>

Add this to PATH if needed:
  export PATH="$BINDIR:\$PATH"
EOF
