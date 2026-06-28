# Changelog

## 0.0.1.280626

- Added artifact-backed large output handling with `artifact.list`,
  `artifact.read_range`, `artifact.search`, and `artifact.stats`.
- Added bounded output envelopes for file reads, git output, job logs, audit
  tails, and agent message tails.
- Added `context.health` and `handoff.create/resume/search/digest` for
  multi-chat continuation from ChatGPT Web without loading full history.
- Added `workspace.snapshot` and `gate.stats`.
- Added allowlisted remote Git project tools: `project.list`,
  `project.ensure`, `project.open`, `project.preflight`, `project.snapshot`,
  and `project.lock_status`.
- Added file-backed task lifecycle tools: `task.start`, `task.status`,
  `task.recover`, `task.digest`, `task.finish_preview`, `task.timeline`, and
  `task.search`.
- Added guarded publish flow: `publish.preview`, `publish.branch`, and
  `publish.pr_create`.
- Added `env_allowlist`, `isolated_home`, `max_artifact_bytes`, and project
  config.
- Added task-bound agent sessions via `task_id`.
- Documented the Go/Rust/TypeScript language strategy.

## 0.0.1.260626

- Created independent AniGate repo.
- Added Mini MCP gateway with stdio and HTTP `/mcp` transports.
- Added bounded tools: `sys.info`, `fs.list`, `fs.read`, `file.search`,
  `app.run_preset`, `job.status`, and `job.logs_tail`.
- Added read-only `git.status` and `git.diff`.
- Added read-only `audit.events_tail`.
- Added file-backed job logs and event log.
- Added `policy.info`, `fs.stat`, `fs.tree`, `fs.write_preview`.
- Added `git.log`, `git.show`, and guarded `patch.apply`.
- Added HTTP bearer token auth for `/mcp`.
- Added workspace profiles: `reader`, `operator`, `agent`.
- Added typed preset arguments.
- Added `job.list` and in-process `job.cancel`.
- Added `audit.summary`.
- Added file-backed agent sessions for long ChatGPT Web to Linux conversations.
