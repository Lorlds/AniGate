# Open Project Readiness

AniGate is usable as a public source-available project, but it is still early.
This checklist tracks what has been done and what would make it easier for
strangers to evaluate, install, trust, and contribute.

## Done

- Public README in English and Chinese.
- License file and explicit noncommercial licensing language.
- Quickstart and installer script.
- Local verification script and `Makefile`.
- GitHub Actions CI.
- Tag-based release workflow with prebuilt archives.
- MCP client examples for stdio and HTTP.
- systemd user service example.
- Contribution guide.
- Security policy.
- Code of conduct.
- Issue templates and PR template.

## Still Worth Adding

- A project logo or GitHub social preview image.
- A short demo GIF or terminal recording.
- A `docs/configuration.md` reference for every config field.
- A `docs/tools.md` reference generated from the tool registry.
- Checksums for release archives.
- Windows release artifacts if Windows support becomes a goal.
- More focused tests around config validation, HTTP auth, publish flows, and
  process cancellation.
- Dependency and vulnerability scanning once non-stdlib dependencies are added.
- A public roadmap with scoped milestones.
- A clearer commercial licensing contact path if commercial users show up.

## Current Positioning

Recommended About text:

```text
Controlled MCP gateway from ChatGPT Web to remote Linux with allowlisted tools,
bounded output, audit logs, jobs, handoff, and agent sessions.
```

Recommended topics:

```text
mcp, model-context-protocol, chatgpt, remote-linux, go, json-rpc, agent,
audit-log, secure-remote-access, developer-tools
```
