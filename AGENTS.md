# AniGate Agent Notes

AniGate is an independent repo. Do not merge it into AniMonitor.

Core product boundary:

- AniGate is a controlled MCP gateway for ChatGPT Web to use a remote Linux
  workspace.
- It must not expose a raw shell.
- It must not collapse into only a Codex/Claude agent wrapper.
- AniMonitor is a reference for event flow, digest/rotate/archive thinking,
  thin MCP entrypoint plus local Go backend, wrappers, hooks, and audit ideas.

Mini implementation rules:

- Keep it Go stdlib-first and single-binary.
- Keep state file-backed; no DB.
- Add capabilities as narrow MCP tools rather than one generic command tool.
- Run only configured presets.
- Keep paths inside configured workspaces.
- Do not inherit the host environment by default.

Future integration:

- AniGate may push task events to an AniMonitor webhook.
- AniMonitor should remain responsible for observation, notification, digest,
  and archive. AniGate should remain responsible for read/invoke/execute.
