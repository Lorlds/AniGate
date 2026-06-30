# Changelog

## 0.1.2

- Added beginner-friendly install and verification scripts.
- Added `Makefile` targets for build, install, test, race, and verification.
- Added Chinese README and a step-by-step user quickstart.
- Added stdio and HTTP MCP client configuration examples.
- Added a user-level systemd service example.
- Added GitHub Actions CI and tag-based release binary publishing.
- Updated the example `sandboxes` workspace path for the new repository
  location so it still points at `/home/lorald/sandboxes`.

## 0.1.1

- Added PolyForm Noncommercial License 1.0.0 and documented AniGate as
  source-available, noncommercial software rather than OSI-approved open source.
- Made publish confirmation tokens one-time use for `publish.branch` and
  `publish.pr_create`.
- Redacted credential-bearing remote URLs from external command error text.
- Removed a full credential-URL-shaped test fixture from source text while
  keeping redaction coverage.

## 0.1.0

- Added Max-only `file.edit_apply` for explicit direct Web GPT single-file
  edits when no configured agent should be used.
- Added `task.commit_preview` and `task.commit` so task worktrees can be
  committed before publish.
- Made `publish.preview` reject dirty task worktrees and point callers to the
  commit flow.
- Added `gate.doctor` and extended `project.preflight` with structured checks.
- Preserved business `next` actions when artifact follow-up actions are added.
- Hardened HTTP startup by requiring auth tokens for non-loopback listeners and
  added stricter server timeouts.
- Tightened new AniGate state, log, artifact, task, handoff, and session file
  permissions.

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
