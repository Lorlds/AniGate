# AniGate

Version: `0.1.3` (`semver`)

[![CI](https://github.com/Lorlds/AniGate/actions/workflows/ci.yml/badge.svg)](https://github.com/Lorlds/AniGate/actions/workflows/ci.yml)
[![Latest Release](https://img.shields.io/github/v/release/Lorlds/AniGate?label=release)](https://github.com/Lorlds/AniGate/releases/latest)
[![License](https://img.shields.io/badge/license-PolyForm%20Noncommercial-blue)](LICENSE)

[中文 README](README.zh-CN.md)

AniGate is a controlled MCP gateway from ChatGPT Web to remote Linux. It is not
a raw shell and it is not just an agent wrapper: every capability is an
allowlisted, bounded, auditable tool.

License: AniGate is source-available under the PolyForm Noncommercial License
1.0.0. Noncommercial use is permitted; commercial use requires separate
permission. Because commercial use is restricted, this is not an OSI-approved
open-source license.

## Language Strategy

- Go is the core runtime: single binary, stdlib-first, file-backed state.
- Rust is reserved for future high-performance or stronger-isolation helper
  components such as an indexer or runner.
- TypeScript is reserved for MCP client examples, schemas, demos, and web-side
  tooling. Permission decisions stay in Go.

## Current Tools

Core and policy:

- `policy.info`, `sys.info`, `gate.stats`, `gate.doctor`, `context.health`

Filesystem and search:

- `fs.list`, `fs.read`, `fs.stat`, `fs.tree`, `file.search`,
  `fs.write_preview`, `file.edit_apply`

Artifact and bounded-output handling:

- `artifact.list`, `artifact.read_range`, `artifact.search`, `artifact.stats`

Git and patch:

- `git.status`, `git.diff`, `git.log`, `git.show`, `patch.apply`

Audit, jobs, and presets:

- `audit.events_tail`, `audit.summary`, `app.run_preset`, `job.list`,
  `job.status`, `job.cancel`, `job.logs_tail`

Long-running agent sessions:

- `agent.session_start`, `agent.message_send`, `agent.session_status`,
  `agent.messages_tail`, `agent.session_list`

Project/task/publish:

- `workspace.snapshot`
- `project.list`, `project.ensure`, `project.open`, `project.preflight`,
  `project.snapshot`, `project.lock_status`
- `task.start`, `task.status`, `task.recover`, `task.digest`,
  `task.finish_preview`, `task.commit_preview`, `task.commit`,
  `task.timeline`, `task.search`
- `publish.preview`, `publish.branch`, `publish.pr_create`

Conversation handoff:

- `handoff.create`, `handoff.resume`, `handoff.search`, `handoff.digest`

## Mini and Max Levels

AniGate does not run two separate binaries for Mini and Max. It exposes one MCP
tool registry, then authorizes each call through workspace policy.

| Level | Intended use | Suggested workspace policy | Examples |
| --- | --- | --- | --- |
| Mini | Read, search, inspect, and preview without changing workspace files. | `profile: "reader"`, `read_only: true` | `fs.read`, `file.search`, `fs.write_preview`, `git.diff`, `artifact.search`, `handoff.*` |
| Max Operator | Controlled execution and workspace mutation. | `profile: "operator"`, `read_only: false` | `app.run_preset`, `patch.apply`, `file.edit_apply`, `task.commit` |
| Max Agent | Long-running configured agent work. | `profile: "agent"`, `read_only: false` | `agent.*`, task-bound agent sessions, publish flow |

Important enforcement details:

- `fs.write_preview` is Mini-safe: it returns a diff and does not write disk.
- `file.edit_apply` and `patch.apply` require a non-read-only `operator` or
  `agent` workspace.
- `app.run_preset` requires `operator` or `agent` profile. Preset commands are
  still configured argv arrays; AniGate does not expose arbitrary shell.
- `agent.*` requires `agent` profile.
- `tools/list` lists all possible tools. Use `policy.info` to inspect current
  workspaces, profiles, and the `capability_levels` map; unauthorized calls fail
  at execution time.
- For a real two-tier deployment, run two configs or two HTTP listeners: a
  Mini config with read-only `reader` workspaces, and a protected Max config
  with `operator` or `agent` workspaces.

## Design Boundaries

- No arbitrary shell tool.
- File access is limited to configured workspaces.
- Remote Git projects must be declared in the `projects` allowlist.
- Applications run only through named presets in config.
- Agent calls run only through configured argv wrappers.
- Preset and agent env vars are explicit and can be constrained by
  `env_allowlist`.
- Jobs run with isolated `HOME` when `isolated_home` is true.
- Large outputs are written to local artifacts; ChatGPT receives bounded
  previews and follow-up tool names.
- Push and PR creation require `publish.preview` and a short-lived confirmation
  token.
- Direct Web GPT edits are explicit Max actions through `file.edit_apply`;
  Mini remains read/search/preview only.
- Jobs, logs, events, artifacts, tasks, handoffs, and sessions are file-backed;
  no database is required.

AniMonitor remains separate. AniMonitor observes, notifies, digests, and
archives. AniGate reads, invokes, executes, hands off, and can later push task
events to an AniMonitor webhook.

## Quick Start

Install to `~/.local/bin/anigate` and generate a user config:

```bash
git clone https://github.com/Lorlds/AniGate.git
cd AniGate
./scripts/install.sh
```

The installer writes a local config to:

```text
~/.config/anigate/anigate.json
```

Run a local HTTP MCP server:

```bash
~/.local/bin/anigate http --addr 127.0.0.1:8787 --config ~/.config/anigate/anigate.json
```

Or run stdio mode for a local MCP client:

```bash
~/.local/bin/anigate stdio --config ~/.config/anigate/anigate.json
```

Developer checkout:

```bash
make verify
make build
./bin/anigate version
./bin/anigate stdio --config configs/anigate.example.json
```

HTTP mode:

```bash
./bin/anigate http --addr 127.0.0.1:8787 --config configs/anigate.example.json
```

Then POST JSON-RPC requests to `http://127.0.0.1:8787/mcp`. For ChatGPT Web,
expose the MCP endpoint with HTTPS or OpenAI Secure MCP Tunnel.

Useful files for new users:

- `README.zh-CN.md`: Chinese quick start and safety notes.
- `docs/user-quickstart.md`: step-by-step setup.
- `CONTRIBUTING.md`: contribution workflow and local verification.
- `SECURITY.md`: supported versions and vulnerability reporting.
- `examples/mcp-client.stdio.json`: local stdio MCP client example.
- `examples/mcp-client.http.json`: HTTP MCP client example.
- `docs/systemd/anigate.service`: optional user-level systemd service.
- `scripts/verify.sh`: local test/build/smoke verification.
- `scripts/install.sh`: local install and config generation.

Release binaries are built by GitHub Actions for `v*` tags.

## Configuration

Copy `configs/anigate.example.json` and adjust:

- `workspaces`: allowed roots for `fs.*`, `file.search`, `git.*`, project
  worktrees, and agents.
- `profile`: `reader`, `operator`, or `agent`.
- `auth_token`: optional HTTP bearer token; `ANIGATE_AUTH_TOKEN` also works.
- `env_allowlist`: optional allowlist for preset/agent env keys.
- `isolated_home`: give jobs a private `HOME` under `state_dir`.
- `presets`: named commands and typed args that `app.run_preset` may execute.
- `agents`: named long-running conversation wrappers.
- `projects`: allowlisted remote Git repositories.
- `state_dir`: file-backed jobs, logs, events, artifacts, tasks, handoffs, and
  sessions.

Preset commands are argv arrays and are executed directly with `exec.Command`;
AniGate does not invoke `/bin/sh -c`.

## Long Agent Conversations

For ChatGPT Web to keep a long interaction with the Linux side, AniGate uses
file-backed sessions:

```text
agent.session_start -> session_id
agent.message_send  -> job_id
job.status          -> poll running/done/failed/cancelled
agent.messages_tail -> read durable conversation output
```

Use `task_id` with `agent.session_start` to bind an agent session to a project
task worktree. The session survives process restarts because messages are stored
under `state_dir/agents/messages`.

Example wrappers:

```json
["codex", "exec", "{prompt}"]
["claude", "-p", "{prompt}"]
```

Use `async: true` for long jobs, then poll `job.status` and
`agent.messages_tail`.

## Context Handoff

AniGate cannot see ChatGPT Web's full token count. It estimates context pressure
from AniGate state: tool outputs, artifacts, jobs, tasks, agent transcript
length, and audit events.

When `context.health` returns `yellow` or `red`, call `handoff.create`. It
creates a compact package and a next-chat prompt. In the new chat, call
`handoff.resume` first, then use `handoff.search`, `artifact.search`, or
`task.recover` only when needed. This keeps older work searchable without
dumping all history into the new conversation.

## Project Flow

```text
project.list
project.ensure
task.start
agent.session_start task_id=<task>
task.finish_preview
task.commit_preview
task.commit
publish.preview
publish.branch or publish.pr_create
```

`project.ensure` only clones/fetches configured allowlist remotes. `task.start`
creates a branch/worktree named `anigate/<task_id>-<slug>`. `publish.branch`
and `publish.pr_create` require the token returned by `publish.preview`.
`publish.preview` rejects dirty task worktrees; call `task.commit_preview` and
`task.commit` first.
