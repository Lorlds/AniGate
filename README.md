# AniGate

Version: `0.0.1.280626` (`0.0.1.DDMMYY`)

AniGate is a controlled MCP gateway from ChatGPT Web to remote Linux. It is not
a raw shell and it is not just an agent wrapper: every capability is an
allowlisted, bounded, auditable tool.

## Language Strategy

- Go is the core runtime: single binary, stdlib-first, file-backed state.
- Rust is reserved for future high-performance or stronger-isolation helper
  components such as an indexer or runner.
- TypeScript is reserved for MCP client examples, schemas, demos, and web-side
  tooling. Permission decisions stay in Go.

## Current Tools

Core and policy:

- `policy.info`, `sys.info`, `gate.stats`, `context.health`

Filesystem and search:

- `fs.list`, `fs.read`, `fs.stat`, `fs.tree`, `file.search`,
  `fs.write_preview`

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
  `task.finish_preview`, `task.timeline`, `task.search`
- `publish.preview`, `publish.branch`, `publish.pr_create`

Conversation handoff:

- `handoff.create`, `handoff.resume`, `handoff.search`, `handoff.digest`

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
- Jobs, logs, events, artifacts, tasks, handoffs, and sessions are file-backed;
  no database is required.

AniMonitor remains separate. AniMonitor observes, notifies, digests, and
archives. AniGate reads, invokes, executes, hands off, and can later push task
events to an AniMonitor webhook.

## Quick Start

```bash
cd /path/to/AniGate
go test ./...
go run ./cmd/anigate version
go run ./cmd/anigate stdio --config configs/anigate.example.json
```

HTTP mode:

```bash
go run ./cmd/anigate http --addr 127.0.0.1:8787 --config configs/anigate.example.json
```

Then POST JSON-RPC requests to `http://127.0.0.1:8787/mcp`. For ChatGPT Web,
expose the MCP endpoint with HTTPS or OpenAI Secure MCP Tunnel.

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
publish.preview
publish.branch or publish.pr_create
```

`project.ensure` only clones/fetches configured allowlist remotes. `task.start`
creates a branch/worktree named `anigate/<task_id>-<slug>`. `publish.branch`
and `publish.pr_create` require the token returned by `publish.preview`.
