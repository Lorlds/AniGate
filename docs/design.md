# AniGate Design

AniGate is a controlled MCP gateway for remote Linux use from ChatGPT Web.

```text
ChatGPT Web -> MCP HTTPS/Tunnel or stdio -> anigate-mini/anigate-max -> policy -> workspace/project/task/job
```

It exposes small, auditable tools. It does not expose `run_shell`.

AniGate is source-available under PolyForm Noncommercial 1.0.0. Noncommercial
use is permitted, but commercial use is restricted, so this is not an
OSI-approved open-source license.

## Runtime Shape

- Core implementation: Go single binary, file-backed state, stdlib-first.
- Optional future helpers: Rust for indexer/runner isolation, TypeScript for MCP
  clients and demos.
- Transports: stdio JSON-RPC and HTTP `/mcp`.
- State: no database in v1.

State layout:

```text
state/
  events.ndjson
  jobs/<job_id>.json
  logs/<job_id>.log
  artifacts/<artifact_id>.{txt,json}
  agents/sessions/<session_id>.json
  agents/messages/<session_id>.ndjson
  tasks/<task_id>.json
  handoffs/<handoff_id>.json
  publish_tokens/<token>.json
  home/
```

## Product Lines

Mini and Max are product lines, not workspace profile aliases. The selected
binary chooses the product line before any workspace authorization runs:

- `anigate-mini`: Safe MCP Preview Gateway.
- `anigate-max`: Controlled Linux MCP Workbench.
- `anigate`: legacy alias for Max.

`tools/list` only returns tools available in the selected product line.
`tools/call` applies the same product gate before dispatch, so a client cannot
call hidden Max tools by name through a Mini server.

Mini is the read/search/diff/preview gateway. It exposes only:

- `sys.info`, `policy.info`
- `fs.list`, `fs.read`, `fs.stat`, `fs.tree`, `file.search`
- `fs.write_preview` for diff-only previews without disk writes
- `git.status`, `git.diff`, `git.log`, `git.show`
- `artifact.*` for large output follow-up reads/search
- `context.health` and `handoff.*` for multi-chat continuation

Mini does not expose execution, mutation, job management, agent, project, task,
publish, audit, workspace snapshot, or gate doctor/stats tools.

Max exposes the complete tool registry. Workspace profile and `read_only` remain
a second authorization layer:

- `patch.apply`, `file.edit_apply`
- `app.run_preset`, `job.*`
- `agent.*` with file-backed sessions
- `audit.*`, `workspace.snapshot`, `gate.stats`, `gate.doctor`
- `project.*`, `task.*`, and `publish.*` actions

For Max, `file.edit_apply` and `patch.apply` require a non-read-only
`operator` or `agent` workspace. `agent.*` requires an `agent` workspace.

## Large Output Policy

Raw large output stays on the Linux side. Tools return:

```json
{
  "text": "bounded preview",
  "truncated": true,
  "artifact": {
    "id": "...",
    "kind": "git.diff",
    "path": "...",
    "bytes": 12345
  },
  "next": ["artifact.read_range", "artifact.search"]
}
```

This applies to file reads, git output, job logs, audit tails, and agent
messages. ChatGPT can search or page artifacts later without loading everything.

## Long Conversations and Handoff

MCP calls from ChatGPT Web should not depend on one permanently open process.
AniGate uses durable ids instead:

```text
task_id -> project worktree/branch
session_id -> agent conversation
job_id -> process/log state
artifact_id -> large output
handoff_id -> next-chat continuation package
```

`context.health` estimates pressure from AniGate-visible state. It cannot read
ChatGPT Web's full token count, so recommendations are heuristic. When pressure
is high, `handoff.create` writes a compact handoff plus a prompt for the next
conversation. The next conversation calls `handoff.resume` and then searches
only what it needs.

## Remote Project Flow

Projects are allowlisted in config. MCP callers cannot pass arbitrary remote
URLs.

```text
project.ensure -> clone/fetch allowlisted remote
task.start     -> branch/worktree
agent.*        -> optional task-bound long agent
file.edit_apply or patch.apply -> controlled edits
task.commit_preview -> verify changes
task.commit    -> commit task worktree
task.digest    -> compact continuation
publish.preview -> confirmation token
publish.branch or publish.pr_create
```

`publish.branch` and `publish.pr_create` require a non-expired token created by
`publish.preview`. `publish.preview` rejects uncommitted task worktrees. Force
push and protected branch pushes are not exposed.

## Security Model

- No arbitrary shell.
- No arbitrary remote Git URL.
- All paths resolve through workspace policy.
- Writable actions require non-read-only `operator` or `agent` workspaces.
- Direct Web GPT writes use `file.edit_apply`; Mini remains read/search/preview
  only.
- Agent execution requires `agent` workspace profile.
- Env vars are explicit and can be limited by `env_allowlist`.
- Jobs use isolated `HOME` when enabled.
- Remote URLs are redacted in tool output.
- Audit events are file-backed and bounded for reads.

## AniMonitor Relationship

AniMonitor should not absorb AniGate.

- AniMonitor: observe, notify, digest, rotate, archive.
- AniGate: read, invoke, execute, hand off tasks.

Integration should be event-based. AniGate can later POST task and handoff
events to an AniMonitor webhook so AniMonitor digest/archive can include remote
execution history.
