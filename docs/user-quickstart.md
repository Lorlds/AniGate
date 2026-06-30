# AniGate User Quickstart

This guide is for people who want to run AniGate without learning the codebase.

## 1. Install Go

AniGate requires Go 1.22 or newer.

Check:

```bash
go version
```

## 2. Install AniGate

```bash
git clone https://github.com/Lorlds/AniGate.git
cd AniGate
./scripts/install.sh
```

The installer builds:

```text
~/.local/bin/anigate
```

It also creates:

```text
~/.config/anigate/anigate.json
```

## 3. Start AniGate

For local MCP clients, stdio mode is the simplest:

```bash
~/.local/bin/anigate stdio --config ~/.config/anigate/anigate.json
```

For HTTP MCP clients:

```bash
~/.local/bin/anigate http --addr 127.0.0.1:8787 --config ~/.config/anigate/anigate.json
```

The HTTP endpoint is:

```text
POST http://127.0.0.1:8787/mcp
```

## 4. Connect an MCP Client

Use one of these examples:

- `examples/mcp-client.stdio.json`
- `examples/mcp-client.http.json`

Replace `/home/YOUR_USER` with your real home directory.

## 5. Keep It Safe

Do not expose AniGate directly to the public internet.

If you listen on anything other than `127.0.0.1`, set a strong `auth_token`.
AniGate refuses to start on non-loopback HTTP addresses without a token.

Keep workspace roots narrow. Start with a read-only workspace, then add
operator or agent access only for directories you are comfortable letting an MCP
client inspect or modify.

## 6. Run as a User Service

Copy the service file:

```bash
mkdir -p ~/.config/systemd/user
cp docs/systemd/anigate.service ~/.config/systemd/user/anigate.service
systemctl --user daemon-reload
systemctl --user enable --now anigate.service
```

Check logs:

```bash
journalctl --user -u anigate.service -f
```

## 7. Verify a Checkout

From the repository:

```bash
make verify
```

If race tests are too slow on your machine:

```bash
ANIGATE_SKIP_RACE=1 make verify
```
