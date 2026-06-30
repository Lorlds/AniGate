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

MINI_BIN="$TMP_DIR/anigate-mini"
MAX_BIN="$TMP_DIR/anigate-max"
LEGACY_BIN="$TMP_DIR/anigate"
MINI_CONFIG=${ANIGATE_VERIFY_MINI_CONFIG:-configs/anigate.mini.example.json}
MAX_CONFIG=${ANIGATE_VERIFY_MAX_CONFIG:-configs/anigate.max.example.json}
LEGACY_CONFIG=${ANIGATE_VERIFY_CONFIG:-configs/anigate.example.json}
ADDR=${ANIGATE_SMOKE_ADDR:-127.0.0.1:18787}

echo "==> go mod verify"
go mod verify

echo "==> go test ./..."
go test ./...

echo "==> go vet ./..."
go vet ./...

if [ "${ANIGATE_SKIP_RACE:-0}" != "1" ]; then
  echo "==> go test -race ./..."
  go test -race ./...
fi

echo "==> go build anigate-mini"
go build -trimpath -o "$MINI_BIN" ./cmd/anigate-mini

echo "==> go build anigate-max"
go build -trimpath -o "$MAX_BIN" ./cmd/anigate-max

echo "==> go build anigate"
go build -trimpath -o "$LEGACY_BIN" ./cmd/anigate

echo "==> version"
"$MINI_BIN" version
"$MAX_BIN" version
"$LEGACY_BIN" version

echo "==> mini tools"
MINI_TOOLS=$("$MINI_BIN" tools --config "$MINI_CONFIG")
MINI_COUNT=$(printf '%s\n' "$MINI_TOOLS" | awk 'NF { count++ } END { print count + 0 }')
if [ "$MINI_COUNT" -ne 21 ]; then
  echo "expected 21 Mini tools, got $MINI_COUNT" >&2
  exit 1
fi
if printf '%s\n' "$MINI_TOOLS" | grep -Eq '^(agent\.|publish\.|file\.edit_apply|patch\.apply|app\.run_preset|job\.|project\.|task\.|audit\.|workspace\.snapshot|gate\.)'; then
  echo "Mini tools exposed a Max-only tool" >&2
  printf '%s\n' "$MINI_TOOLS" >&2
  exit 1
fi
echo "$MINI_COUNT Mini tools"

echo "==> max tools"
MAX_TOOLS=$("$MAX_BIN" tools --config "$MAX_CONFIG")
MAX_COUNT=$(printf '%s\n' "$MAX_TOOLS" | awk 'NF { count++ } END { print count + 0 }')
if [ "$MAX_COUNT" -ne 56 ]; then
  echo "expected 56 Max tools, got $MAX_COUNT" >&2
  exit 1
fi
for tool in file.edit_apply patch.apply agent.message_send publish.preview; do
  if ! printf '%s\n' "$MAX_TOOLS" | grep -q "^$tool[	 ]"; then
    echo "Max tools missing $tool" >&2
    exit 1
  fi
done
echo "$MAX_COUNT Max tools"

echo "==> legacy tools"
LEGACY_COUNT=$("$LEGACY_BIN" tools --config "$LEGACY_CONFIG" | awk 'NF { count++ } END { print count + 0 }')
if [ "$LEGACY_COUNT" -ne 56 ]; then
  echo "expected 56 legacy Max tools, got $LEGACY_COUNT" >&2
  exit 1
fi
echo "$LEGACY_COUNT legacy tools"

if command -v python3 >/dev/null 2>&1; then
  echo "==> HTTP MCP smoke test on $ADDR"
  "$MINI_BIN" http --addr "$ADDR" --config "$MINI_CONFIG" >/tmp/anigate-verify-http.log 2>&1 &
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
