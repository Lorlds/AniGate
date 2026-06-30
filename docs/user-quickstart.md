# AniGate User Quickstart

This guide is for people who want to run AniGate without learning the codebase.
Start with AniGate Mini. Use AniGate Max only when you need controlled
execution, edits, agents, tasks, or publishing.

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
~/.local/bin/anigate-mini
~/.local/bin/anigate-max
~/.local/bin/anigate
```

It also creates:

```text
~/.config/anigate/anigate-mini.json
~/.config/anigate/anigate-max.json
~/.config/anigate/anigate.json
```

## 3. Start AniGate Mini

For local MCP clients, stdio mode is the simplest:

```bash
~/.local/bin/anigate-mini stdio --config ~/.config/anigate/anigate-mini.json
```

For HTTP MCP clients:

```bash
~/.local/bin/anigate-mini http --addr 127.0.0.1:8787 --config ~/.config/anigate/anigate-mini.json
```

The Mini HTTP endpoint is:

```text
POST http://127.0.0.1:8787/mcp
```

## 4. Connect an MCP Client

Use one of these Mini examples:

- `examples/mcp-client.mini.stdio.json`
- `examples/mcp-client.mini.http.json`

Replace `/home/YOUR_USER` with your real home directory.

## 5. Use Max Deliberately

Max exposes execution, mutation, job, agent, project, task, publish, audit, and
diagnostic tools. Run it on a separate port and protect it with a token:

```bash
~/.local/bin/anigate-max http --addr 127.0.0.1:8788 --config ~/.config/anigate/anigate-max.json
```

Max client examples:

- `examples/mcp-client.max.stdio.json`
- `examples/mcp-client.max.http.json`

## 6. Keep It Safe

Do not expose AniGate directly to the public internet.

If you listen on anything other than `127.0.0.1`, set a strong `auth_token`.
AniGate refuses to start on non-loopback HTTP addresses without a token.

Keep workspace roots narrow. Start with Mini and a read-only workspace, then add
Max access only for directories you are comfortable letting an MCP client modify
or execute configured tools inside.

## 7. Run as a User Service

Copy the Mini service file:

```bash
mkdir -p ~/.config/systemd/user
cp docs/systemd/anigate-mini.service ~/.config/systemd/user/anigate-mini.service
systemctl --user daemon-reload
systemctl --user enable --now anigate-mini.service
```

Check logs:

```bash
journalctl --user -u anigate-mini.service -f
```

Max has a separate service example at `docs/systemd/anigate-max.service`.

## 8. Verify a Checkout

From the repository:

```bash
make verify
```

If race tests are too slow on your machine:

```bash
ANIGATE_SKIP_RACE=1 make verify
```
