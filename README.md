# condition-github-actions

GitHub Actions CI condition plugin for [go-semrel](https://github.com/SemRels/semrel).

Verifies that a release is being published from a GitHub Actions workflow that
is properly configured to do so, preventing accidental releases from pull
requests, tag pushes, or environments without a valid GitHub token.

## Checks performed

| Check | Environment variable | Expected value |
|-------|----------------------|----------------|
| Running in GitHub Actions | `GITHUB_ACTIONS` | `"true"` |
| Triggered by a branch push | `GITHUB_REF_TYPE` | `"branch"` |
| Not a pull-request run | `GITHUB_EVENT_NAME` | anything except `"pull_request"` |
| Token available | `GITHUB_TOKEN` | non-empty |
| Branch matches release context *(optional)* | `GITHUB_REF_NAME` / `GITHUB_HEAD_REF` | must match `ctx.branch` when provided |

If any mandatory check fails the plugin returns a non-empty `error_message`
together with a human-readable `details` field containing a numbered failure
list and remediation hints.

## Repository layout

~~~text
api/gen/v1/              Generated gRPC / protobuf bindings (from proto/v1)
cmd/plugin/              Plugin binary entry point
internal/grpc/           gRPC adapter (CIConditionServer)
internal/plugin/         Core condition-check logic (ConditionChecker, EnvProvider)
proto/v1/                Vendored protobuf contract (semantic_release.proto)
.github/workflows/       CI, release, and security automation
~~~

## Protocol

This plugin implements the `CIConditionPlugin` gRPC service defined in
`proto/v1/semantic_release.proto` and is launched by the go-semrel host via
[hashicorp/go-plugin](https://github.com/hashicorp/go-plugin).

- **stdout** is reserved exclusively for the go-plugin handshake.
- All log output goes to **stderr**.

## Configuration example (`.semrel.yaml`)

~~~yaml
plugins:
  - name: condition-github-actions
    type: condition
    config: {}       # no additional config keys required
~~~

## Minimal workflow example

~~~yaml
name: Release
on:
  push:
    branches: [main]

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run go-semrel
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: semrel release
~~~

## Error scenarios

| Symptom | Cause | Fix |
|---------|-------|-----|
| `GITHUB_ACTIONS is not 'true'` | Plugin executed outside GitHub Actions | Run inside a GitHub Actions workflow |
| `GITHUB_REF_TYPE is "tag"` | Workflow triggered by a tag push | Use `on: push: branches:` not `tags:` |
| `GITHUB_EVENT_NAME is 'pull_request'` | Workflow triggered by a PR | Add a branch condition to exclude PRs |
| `GITHUB_TOKEN is not set` | Token env var missing | Add `env: GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}` |
| `branch mismatch` | `ctx.branch` does not match `GITHUB_REF_NAME` | Align `.semrel.yaml` branch with the workflow trigger branch |

## Development

~~~bash
# Build
make build          # → bin/plugin

# Test (all packages, with coverage)
make test
go test -v -cover ./...

# Lint (requires golangci-lint)
make lint

# Smoke test – should print startup log to stderr then exit with go-plugin message
go run cmd/plugin/main.go
~~~

## Links

- SemRel docs: <https://github.com/SemRels/semrel>
- Plugin template: <https://github.com/SemRels/plugin-template>
- ADR-001 (gRPC plugin transport): `semrel/docs/adr/ADR-001-grpc-plugin-transport.md`
