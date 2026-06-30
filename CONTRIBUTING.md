# Contributing to AniGate

Thanks for taking a look at AniGate.

AniGate is source-available under the PolyForm Noncommercial License 1.0.0.
Noncommercial contributions are welcome. If your use case is commercial, open
an issue first so licensing expectations are clear before significant work.

## Before You Start

AniGate is a controlled MCP gateway. Contributions should preserve these
boundaries:

- Do not add a raw shell tool.
- Keep file access inside configured workspaces.
- Prefer narrow, named MCP tools over generic command execution.
- Keep outputs bounded; large outputs should become local artifacts.
- Keep state file-backed; do not add a database dependency.
- Keep the Go runtime stdlib-first unless a dependency is clearly justified.
- Treat auth tokens, remote URLs, environment variables, and logs as sensitive.

## Local Setup

```bash
git clone https://github.com/Lorlds/AniGate.git
cd AniGate
make verify
```

If race tests are too slow on your machine:

```bash
ANIGATE_SKIP_RACE=1 make verify
```

## Common Commands

```bash
make build
make test
make vet
make race
make tools
make run-http
```

## Pull Request Checklist

Before opening a PR:

- Run `make verify`.
- Update `README.md` or `README.zh-CN.md` when user-facing behavior changes.
- Update `CHANGELOG.md` for release-worthy changes.
- Add or update tests for security boundaries, path policy, auth, jobs,
  artifacts, project/task flows, and MCP protocol behavior.
- Keep diffs focused; avoid unrelated formatting churn.

## Good First Issues

Good first contributions usually fit one of these shapes:

- Documentation improvements.
- More table-driven tests for path policy and config validation.
- More examples for MCP client configuration.
- Safer defaults or clearer error messages.
- Small tool-specific tests that improve coverage without broad refactors.

## Security Changes

For security-sensitive changes, prefer a small PR with a clear threat model:

- What could go wrong before the change?
- What input triggers the case?
- What test proves the boundary holds?

If you believe you found a vulnerability, follow `SECURITY.md` instead of
opening a public issue.
