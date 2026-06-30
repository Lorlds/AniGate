#!/usr/bin/env sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"

TMP_DIR=$(mktemp -d)
cleanup() {
  if [ -n "${SMOKE_PID:-}" ]; then
    kill "$SMOKE_PID" >/dev/null 2>&1 || true
    wait "$SMOKE_PID" >/dev/null 2>&1 || true
  fi
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

BIN="$TMP_DIR/anigate"
CONFIG=${ANIGATE_VERIFY_CONFIG:-configs/anigate.example.json}
ADDR=${ANIGATE_SMOKE_ADDR:-127.0.0.1:18787}

echo "==> go test ./..."
go test ./...

echo "==> go vet ./..."
go vet ./...

if [ "${ANIGATE_SKIP_RACE:-0}" != "1" ]; then
  echo "==> go test -race ./..."
  go test -race ./...
fi

echo "==> go build"
go build -trimpath -o "$BIN" ./cmd/anigate

echo "==> version"
"$BIN" version

echo "==> tools"
TOOL_COUNT=$("$BIN" tools --config "$CONFIG" | awk 'NF { count++ } END { print count + 0 }')
if [ "$TOOL_COUNT" -lt 1 ]; then
  echo "expected at least one exposed tool" >&2
  exit 1
fi
echo "$TOOL_COUNT tools"

if command -v python3 >/dev/null 2>&1; then
  echo "==> HTTP MCP smoke test on $ADDR"
  "$BIN" http --addr "$ADDR" --config "$CONFIG" >/tmp/anigate-verify-http.log 2>&1 &
  SMOKE_PID=$!
  sleep 1
  ANIGATE_SMOKE_ADDR="$ADDR" python3 - <<'PY'
import json
import os
import urllib.request

addr = os.environ["ANIGATE_SMOKE_ADDR"]
url = "http://%s/mcp" % addr

def post(method, params=None):
    payload = json.dumps({
        "jsonrpc": "2.0",
        "id": 1,
        "method": method,
        "params": params or {},
    }).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=payload,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(req, timeout=5) as res:
        return res.status, json.loads(res.read().decode("utf-8"))

status, body = post("initialize", {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {"name": "anigate-verify", "version": "1"},
})
if status != 200 or body.get("result", {}).get("serverInfo", {}).get("name") != "anigate":
    raise SystemExit("initialize smoke test failed: %r" % body)

status, body = post("tools/list")
tools = body.get("result", {}).get("tools", [])
if status != 200 or not tools:
    raise SystemExit("tools/list smoke test failed: %r" % body)

print("HTTP MCP smoke ok: %d tools" % len(tools))
PY
else
  echo "==> skipping HTTP MCP smoke test; python3 not found"
fi

echo "==> verify ok"
